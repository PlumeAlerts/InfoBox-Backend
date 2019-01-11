package main

import (
	"encoding/json"
	"fmt"
	"github.com/PlumeAlerts/InfoBoxes-Backend/db"
	"github.com/PlumeAlerts/InfoBoxes-Backend/extension"
	"github.com/PlumeAlerts/InfoBoxes-Backend/requests"
	"github.com/PlumeAlerts/InfoBoxes-Backend/utilities"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"os"
	"time"
)

func Initialize() {
	client := extension.Client{ClientID: os.Getenv("EXT_CLIENT_ID"), ClientSecret: os.Getenv("EXT_CLIENT_SECRET"), OwnerID: os.Getenv("EXT_OWNER_ID")}

	utilities.Validate = validator.New()
	db.Connect()
	go InitializeSender(client)

	r := mux.NewRouter()
	r.Use(client.VerifyJWT)

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
	log.Fatal(http.ListenAndServe(":8000", handler))
	db.DB.Close()
}

func InitializeSender(client extension.Client) {

	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				triggerInfoBoxes(client)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func triggerInfoBoxes(client extension.Client) {
	var users []db.User
	db.DB.Where("last_triggered < DATE_SUB(NOW(), INTERVAL intervals MINUTE)").Order("last_triggered").Find(&users)

	for i := 0; i < len(users); i++ {
		user := users[i]

		var infoBox db.InfoBox
		db.DB.Where("id = (select min(id) from info_boxes where user_id = ? and id > ?)", user.ID, user.LastInfoBoxesId).Find(&infoBox)
		if infoBox.ID == 0 {
			db.DB.Where("id = (select min(id) from info_boxes where user_id = ? and id > ?)", user.ID, 0).Find(&infoBox)
		}
		b, _ := json.Marshal(infoBox)

		err := client.PostPubSubMessage(user.ID, string(b))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(user.ID)
		}
		db.DB.Save(&db.User{ID: user.ID, LastTriggered: time.Now(), LastInfoBoxesId: infoBox.ID, Intervals: user.Intervals})
	}
}
