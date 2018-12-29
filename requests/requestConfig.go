package requests

import (
	"encoding/json"
	"fmt"
	"github.com/PlumeAlerts/InfoBox-Backend/db"
	"github.com/PlumeAlerts/InfoBox-Backend/utilities"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"net/http"
)

func GetConfig(w http.ResponseWriter, r *http.Request) {
	userId := db.GetUserOrCreate(r)

	var ib db.User
	db.DB.Where(&db.User{ID: userId}).Find(&ib)

	b, err := utilities.InterfaceToJson(&ib)
	if err {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(b)
}

func PutConfig(w http.ResponseWriter, r *http.Request) {
	userId := db.GetUserOrCreate(r)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var data db.User
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

	infoBox := db.User{}
	dbIB := db.DB.Where(&db.User{ID: userId}).First(&infoBox)
	if dbIB.RowsAffected == 0 {
		//TODO return invalid request
		w.WriteHeader(400)
		return
	}

	infoBoxes := &db.User{
		ID:              userId,
		InfoboxInterval: data.InfoboxInterval,
	}
	ib := db.DB.Update(infoBoxes)

	if ib.Error != nil {
		panic(ib.Error.Error())
	}
	b, errs := utilities.InterfaceToJson(&ib.Value)
	if errs {
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(200)
	w.Write(b)
}
