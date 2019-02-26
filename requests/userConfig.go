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

func GetConfig(w http.ResponseWriter, r *http.Request) {
	userId := db.GetUserIdOrCreate(r)

	var ib db.User
	dbResponse := db.DB.Where(&db.User{Id: userId}).Find(&ib)
	if dbResponse.Error != nil {
		resp.NewResponse(w).InternalServerError(utilities.Error{Message: "Database connection error"})
		return
	}

	resp.NewResponse(w).Ok(&ib)
}

func PutConfig(w http.ResponseWriter, r *http.Request) {
	userId := db.GetUserIdOrCreate(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.NewResponse(w).InternalServerError(utilities.Error{Message: "Error reading body"})
		return
	}
	var data db.User
	err = json.Unmarshal(body, &data)
	if err != nil {
		resp.NewResponse(w).BadRequest(utilities.Error{Message: "Error parsing body"})
		return
	}

	err = utilities.ValidateInterface(data)

	if err != nil {
		//TODO Add error message
		validationErrors := err.(validator.ValidationErrors)
		fmt.Println(validationErrors.Error())
		resp.NewResponse(w).BadRequest(nil)
		return
	}

	user := db.User{
		Id:                 userId,
		AnnotationInterval: data.AnnotationInterval,
	}
	dbResponse := db.DB.Save(&user)

	if dbResponse.Error != nil {
		resp.NewResponse(w).InternalServerError(utilities.Error{Message: "Database connection error"})
		return
	}

	resp.NewResponse(w).Ok(&dbResponse.Value)
}
