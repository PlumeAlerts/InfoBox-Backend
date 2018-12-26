package requests

import (
	"encoding/json"
	"fmt"
	"github.com/PlumeAlerts/InfoBox-Backend/db"
	"github.com/PlumeAlerts/InfoBox-Backend/jwt"
	"gopkg.in/go-playground/validator.v9"
	"io/ioutil"
	"net/http"
	"strconv"
)

var Validate *validator.Validate

func GetIBConfig(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(jwt.ChannelIDKey).(string)

	_, created := db.GetUserOrCreate(userId)
	if created {
		b, err := InterfaceToJson(&[]db.InfoBox{})
		if err {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(b)
	}

	var ib []db.InfoBox
	db.DB.Where(&db.InfoBox{UserId: userId}).Find(&ib)

	b, err := InterfaceToJson(&ib)
	if err {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(b)
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

	if _, ok := err.(*validator.InvalidValidationError); ok {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	err = Validate.Struct(&data)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		println(validationErrors.Error())
		w.WriteHeader(400)
		return
	}

	infoBoxes := &db.InfoBox{
		ID:              data.ID,
		Title:           data.Title,
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
		panic(ib.Error.Error())
	}
	b, error := InterfaceToJson(&ib.Value)
	if error {
		w.WriteHeader(400)
		return
		//panic(err)
	}
	w.WriteHeader(200)
	w.Write(b)
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
		w.WriteHeader(400)
		return
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	err = Validate.Struct(&data)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		println(validationErrors.Error())
		w.WriteHeader(400)
		return
	}

	db.GetUserOrCreate(userId)

	infoBoxes := &db.InfoBox{Title: data.Title,
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
		panic(ib.Error.Error())
	}
	b, error := InterfaceToJson(&ib.Value)
	if error {
		w.WriteHeader(400)
		return
		//panic(err)
	}
	w.WriteHeader(200)
	w.Write(b)
}

func DeleteIBConfig(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(jwt.ChannelIDKey).(string)

	id, err := getId(r.FormValue("id"))
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
	w.WriteHeader(200)
}

func getId(id string) (uint, error) {
	err := Validate.Var(id, "required,numeric")
	if err != nil {
		return 0, err
	}

	t, _ := strconv.ParseUint(id, 10, 32)
	return uint(t), nil
}
func InterfaceToJson(obj interface{}) ([]byte, bool) {
	b, err := json.Marshal(&obj)
	if err != nil {
		return nil, true
	}
	return b, false
}
