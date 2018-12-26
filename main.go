package main

import (
	"github.com/PlumeAlerts/InfoBox-Backend/db"
	"github.com/PlumeAlerts/InfoBox-Backend/jwt"
	"github.com/PlumeAlerts/InfoBox-Backend/requests"

	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"regexp"
)

func main() {
	requests.Validate = validator.New()
	requests.Validate.RegisterValidation("text", ValidateText)

	db.Connect()

	r := mux.NewRouter()
	r.Use(jwt.VerifyJWT)
	r.Handle("/ib/config", http.HandlerFunc(requests.GetIBConfig)).Methods("GET")
	r.Handle("/ib/config", http.HandlerFunc(requests.PutIBConfig)).Methods("PUT")
	r.Handle("/ib/config", http.HandlerFunc(requests.PostIBConfig)).Methods("POST")
	r.Handle("/ib/config", http.HandlerFunc(requests.DeleteIBConfig)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))

	//message := Message{
	//	Title:       "Twitter",
	//	IconURL:     "https://localhost.rig.twitch.tv:8080/Twitter_Logo_Blue.svg",
	//	URL:         "https://twitter.com/lclc98",
	//	Description: "Checkout my twitter",
	//	Duration:    10,
	//}
	//
	//SendPubSubMessage("30550166", message)
}

func ValidateText(fl validator.FieldLevel) bool {
	return regexp.MustCompile("\\p{M}*").MatchString(fl.Field().String())
}
