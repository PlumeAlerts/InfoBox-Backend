package main

import (
	"github.com/PlumeAlerts/InfoBox-Backend/db"
	"github.com/PlumeAlerts/InfoBox-Backend/jwt"
	"github.com/PlumeAlerts/InfoBox-Backend/requests"
	"github.com/PlumeAlerts/InfoBox-Backend/utilities"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"regexp"
)

func Initialize() {
	utilities.Validate = validator.New()
	utilities.Validate.RegisterValidation("text", ValidateText)

	db.Connect()

	r := mux.NewRouter()
	r.Use(jwt.VerifyJWT)
	r.Handle("/config", http.HandlerFunc(requests.GetConfig)).Methods("GET")
	r.Handle("/config", http.HandlerFunc(requests.PutConfig)).Methods("PUT")

	r.Handle("/ib/config", http.HandlerFunc(requests.GetIBConfig)).Methods("GET")
	r.Handle("/ib/config", http.HandlerFunc(requests.PutIBConfig)).Methods("PUT")
	r.Handle("/ib/config", http.HandlerFunc(requests.PostIBConfig)).Methods("POST")
	r.Handle("/ib/config", http.HandlerFunc(requests.DeleteIBConfig)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func ValidateText(fl validator.FieldLevel) bool {
	return regexp.MustCompile("\\p{M}*").MatchString(fl.Field().String())
}
