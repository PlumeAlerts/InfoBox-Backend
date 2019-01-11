package requests

import (
	"encoding/json"
	"fmt"
	"github.com/PlumeAlerts/InfoBoxes-Backend/db"
	"github.com/PlumeAlerts/InfoBoxes-Backend/utilities"
	resp "github.com/nicklaw5/go-respond"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"net/http"
)

func GetIBConfig(w http.ResponseWriter, r *http.Request) {
	userId := db.GetUserOrCreate(r)

	var ib []db.InfoBox
	db.DB.Where(&db.InfoBox{UserId: userId}).Find(&ib)

	resp.NewResponse(w).Ok(&ib)
}

func PutIBConfig(w http.ResponseWriter, r *http.Request) {
	userId := db.GetUserOrCreate(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var data db.InfoBox
	err = json.Unmarshal(body, &data)
	if err != nil {
		resp.NewResponse(w).BadRequest(nil)
		return
	}

	err = utilities.ValidateInterface(data)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		fmt.Println(validationErrors.Error())
		resp.NewResponse(w).BadRequest(nil)
		return
	}

	id, err := utilities.GetIBID(r.FormValue("id"))
	if err != nil {
		resp.NewResponse(w).BadRequest(nil)
		fmt.Println(err)
		return
	}

	infoBox := db.InfoBox{}
	dbIB := db.DB.Where(&db.InfoBox{ID: id}).First(&infoBox)
	if dbIB.RowsAffected == 0 {
		resp.NewResponse(w).BadRequest(nil)
		//TODO return invalid request
		return
	}

	if infoBox.UserId != userId {
		//TODO Return unauthorized
		resp.NewResponse(w).Unauthorized(nil)
		return
	}
	infoBoxes := &db.InfoBox{
		ID:              data.ID,
		Text:            data.Text,
		TextSize:        data.TextSize,
		URL:             data.URL,
		Icon:            data.Icon,
		IconColor:       data.IconColor,
		TextColor:       data.TextColor,
		BackgroundColor: data.BackgroundColor,
		UserId:          userId,
	}
	ib := db.DB.Save(infoBoxes)

	if ib.Error != nil {
		panic(ib.Error.Error())
	}
	resp.NewResponse(w).Ok(&ib.Value)
}

func PostIBConfig(w http.ResponseWriter, r *http.Request) {
	userId := db.GetUserOrCreate(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var data db.InfoBox
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		resp.NewResponse(w).BadRequest(nil)
		return
	}

	err = utilities.ValidateInterface(&data)
	if _, ok := err.(*validator.ValidationErrors); ok {
		validationErrors := err.(validator.ValidationErrors)
		fmt.Println(validationErrors.Error())
		resp.NewResponse(w).BadRequest(nil)
		return
	}

	infoBoxes := &db.InfoBox{
		Text:            data.Text,
		TextSize:        data.TextSize,
		URL:             data.URL,
		Icon:            data.Icon,
		IconColor:       data.IconColor,
		TextColor:       data.TextColor,
		BackgroundColor: data.BackgroundColor,
		UserId:          userId,
	}
	ib := db.DB.Create(infoBoxes)

	if ib.Error != nil {
		resp.NewResponse(w).InternalServerError(nil)
	}

	resp.NewResponse(w).Ok(&ib.Value)
}

func DeleteIBConfig(w http.ResponseWriter, r *http.Request) {
	userId := db.GetUserOrCreate(r)

	id, err := utilities.GetIBID(r.FormValue("id"))
	if err != nil {
		resp.NewResponse(w).BadRequest(nil)
		return
	}
	infoBox := db.InfoBox{}
	dbIB := db.DB.Where(&db.InfoBox{ID: id}).First(&infoBox)
	if dbIB.RowsAffected == 0 {
		resp.NewResponse(w).BadRequest(nil)
		//TODO return invalid request
		return
	}

	if infoBox.UserId != userId {
		resp.NewResponse(w).Unauthorized(nil)
		//TODO Return unauthorized
		return
	}

	db.DB.Delete(&db.InfoBox{ID: id})

	resp.NewResponse(w).Ok(nil)
}
