package requests

import (
	"encoding/json"
	"fmt"
	"github.com/PlumeAlerts/InfoBox-Backend/db"
	"github.com/PlumeAlerts/InfoBox-Backend/jwt"
	"github.com/PlumeAlerts/InfoBox-Backend/utilities"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"net/http"
)

func GetIBConfig(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(jwt.ChannelIDKey).(string)

	var ib []db.InfoBox
	db.DB.Where(&db.InfoBox{UserId: userId}).Find(&ib)

	utilities.RespondWithJSON(w, 200, &ib)
}

func PutIBConfig(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(jwt.ChannelIDKey).(string)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var data db.InfoBox
	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	err = utilities.ValidateInterface(data)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		fmt.Println(validationErrors.Error())
		w.WriteHeader(400)
		return
	}

	id, err := utilities.GetIBID(r.FormValue("id"))
	if err != nil {
		w.WriteHeader(400)
		fmt.Println(err)
		return
	}

	infoBox := db.InfoBox{}
	dbIB := db.DB.Where(&db.InfoBox{ID: id}).First(&infoBox)
	if dbIB.RowsAffected == 0 {
		//TODO return invalid request
		w.WriteHeader(400)
		return
	}

	if infoBox.UserId != userId {
		//TODO Return unauthorized
		w.WriteHeader(403)
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
		Intervals:       data.Intervals,
		UserId:          userId,
	}
	ib := db.DB.Save(infoBoxes)

	if ib.Error != nil {
		panic(ib.Error.Error())
	}
	utilities.RespondWithJSON(w, 200, &ib.Value)
}

func PostIBConfig(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(jwt.ChannelIDKey).(string)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var data db.InfoBox
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	err = utilities.ValidateInterface(&data)
	if _, ok := err.(*validator.ValidationErrors); ok {
		validationErrors := err.(validator.ValidationErrors)
		fmt.Println(validationErrors.Error())
		w.WriteHeader(400)
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
		Intervals:       data.Intervals,

		UserId: userId,
	}
	ib := db.DB.Create(infoBoxes)

	if ib.Error != nil {
		panic(ib.Error.Error())
	}
	utilities.RespondWithJSON(w, 200, &ib.Value)
}

func DeleteIBConfig(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(jwt.ChannelIDKey).(string)

	id, err := utilities.GetIBID(r.FormValue("id"))
	if err != nil {
		w.WriteHeader(400)
		fmt.Println(err)
		return
	}
	infoBox := db.InfoBox{}
	dbIB := db.DB.Where(&db.InfoBox{ID: id}).First(&infoBox)
	if dbIB.RowsAffected == 0 {
		//TODO return invalid request
		w.WriteHeader(400)
		return
	}

	if infoBox.UserId != userId {
		//TODO Return unauthorized
		w.WriteHeader(403)
		return
	}

	db.DB.Delete(&db.InfoBox{ID: id})
	utilities.RespondWithJSON(w, 200, nil)
}
