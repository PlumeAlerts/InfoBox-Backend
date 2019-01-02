package main

import (
	"github.com/PlumeAlerts/InfoBoxes-Backend/db"
	"github.com/PlumeAlerts/InfoBoxes-Backend/jwt"
	"github.com/PlumeAlerts/InfoBoxes-Backend/requests"
	"github.com/PlumeAlerts/InfoBoxes-Backend/utilities"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
)

func Initialize() {
	utilities.Validate = validator.New()

	db.Connect()

	r := mux.NewRouter()
	r.Use(jwt.VerifyJWT)

	r.Handle("/config", http.HandlerFunc(requests.GetIBConfig)).Methods("GET")
	r.Handle("/config", http.HandlerFunc(requests.PutIBConfig)).Methods("PUT")
	r.Handle("/config", http.HandlerFunc(requests.PostIBConfig)).Methods("POST")
	r.Handle("/config", http.HandlerFunc(requests.DeleteIBConfig)).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	// Insert the middleware
	handler := c.Handler(r)
	log.Fatal(http.ListenAndServeTLS(":8000", "conf/server.crt", "conf/server.key", handler))
}
