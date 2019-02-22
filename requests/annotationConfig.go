package requests

import (
	"encoding/json"
	"fmt"
	"github.com/PlumeAlerts/StreamAnnotations-Backend/db"
	"github.com/PlumeAlerts/StreamAnnotations-Backend/utilities"
	resp "github.com/nicklaw5/go-respond"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"net/http"
)

func GetAnnotationConfig(w http.ResponseWriter, r *http.Request) {
	userId := db.GetUserIdOrCreate(r)

	// Get the annotations for the user
	var annotation []db.Annotation
	dbResponse := db.DB.Where(&db.Annotation{UserId: userId}).Find(&annotation)

	if dbResponse.Error != nil {
		resp.NewResponse(w).InternalServerError(utilities.Error{Message: "Database connection error"})
		return
	}

	resp.NewResponse(w).Ok(&annotation)
}

func PutAnnotationConfig(w http.ResponseWriter, r *http.Request) {
	userId := db.GetUserIdOrCreate(r)

	// Reads the body and converts it to an annotation
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.NewResponse(w).InternalServerError(utilities.Error{Message: "Error reading body"})
		return
	}
	var data db.Annotation
	err = json.Unmarshal(body, &data)
	if err != nil {
		resp.NewResponse(w).BadRequest(utilities.Error{Message: "Error parsing body"})
		return
	}

	// Validates the data
	err = utilities.ValidateInterface(data)
	if err != nil {
		//TODO Add error message
		validationErrors := err.(validator.ValidationErrors)
		fmt.Println(validationErrors.Error())
		resp.NewResponse(w).BadRequest(err)
		return
	}

	// Gets the annotation id for the updated data
	id, err := utilities.GetAnnotationID(r.FormValue("id"))
	if err != nil {
		resp.NewResponse(w).BadRequest(utilities.Error{Message: "Error getting annotation id"})
		return
	}

	// Gets the annotation with that id
	dbAnnotation := db.Annotation{}
	dbResponse := db.DB.Where(&db.Annotation{Id: id}).First(&dbAnnotation)
	if dbResponse.RowsAffected == 0 {
		// if the annotation isn't found return a bad request
		resp.NewResponse(w).BadRequest(utilities.Error{Message: "Failed to find annotation"})
		return
	}

	// if the annotation doesn't belong the user, return an unauthorized request
	if dbAnnotation.UserId != userId {
		resp.NewResponse(w).Unauthorized(utilities.Error{Message: "Unauthorized user trying to access annotation"})
		return
	}
	annotation := &db.Annotation{
		Id:              id,
		Text:            data.Text,
		TextSize:        data.TextSize,
		URL:             data.URL,
		Icon:            data.Icon,
		IconColor:       data.IconColor,
		TextColor:       data.TextColor,
		BackgroundColor: data.BackgroundColor,
		UserId:          userId,
	}
	dbResponse = db.DB.Save(annotation)

	if dbResponse.Error != nil {
		resp.NewResponse(w).InternalServerError(utilities.Error{Message: "Database connection error"})
		return
	}
	resp.NewResponse(w).Ok(&dbResponse.Value)
}

func PostAnnotationConfig(w http.ResponseWriter, r *http.Request) {
	userId := db.GetUserIdOrCreate(r)

	// Reads the body and converts it to an annotation
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.NewResponse(w).InternalServerError(utilities.Error{Message: "Error reading body"})
		return
	}
	var data db.Annotation
	err = json.Unmarshal(body, &data)
	if err != nil {
		resp.NewResponse(w).BadRequest(utilities.Error{Message: "Error parsing body"})
		return
	}

	// Validates the data
	err = utilities.ValidateInterface(&data)
	if _, ok := err.(*validator.ValidationErrors); ok {
		//TODO Add error message
		validationErrors := err.(validator.ValidationErrors)
		fmt.Println(validationErrors.Error())
		resp.NewResponse(w).BadRequest(nil)
		return
	}

	annotation := &db.Annotation{
		Text:            data.Text,
		TextSize:        data.TextSize,
		URL:             data.URL,
		Icon:            data.Icon,
		IconColor:       data.IconColor,
		TextColor:       data.TextColor,
		BackgroundColor: data.BackgroundColor,
		UserId:          userId,
	}
	dbResponse := db.DB.Create(annotation)

	if dbResponse.Error != nil {
		resp.NewResponse(w).InternalServerError(utilities.Error{Message: "Database connection error"})
		return
	}

	resp.NewResponse(w).Ok(&dbResponse.Value)
}

func DeleteAnnotationConfig(w http.ResponseWriter, r *http.Request) {
	userId := db.GetUserIdOrCreate(r)

	id, err := utilities.GetAnnotationID(r.FormValue("id"))
	if err != nil {
		resp.NewResponse(w).BadRequest(utilities.Error{Message: "Error getting annotation id"})
		return
	}
	annotation := db.Annotation{}
	dbAnnotation := db.DB.Where(&db.Annotation{Id: id}).First(&annotation)
	if dbAnnotation.RowsAffected == 0 {
		resp.NewResponse(w).BadRequest(utilities.Error{Message: "Failed to find annotation"})
		return
	}

	if annotation.UserId != userId {
		resp.NewResponse(w).Unauthorized(utilities.Error{Message: "Unauthorized user trying to access annotation"})
		return
	}

	dbResponse := db.DB.Delete(&db.Annotation{Id: id})
	if dbResponse.Error != nil {
		resp.NewResponse(w).InternalServerError(utilities.Error{Message: "Database connection error"})
		return
	}

	resp.NewResponse(w).Ok(nil)
}
