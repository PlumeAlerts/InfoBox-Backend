package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PlumeAlerts/StreamAnnotations-Backend/db"
	"github.com/PlumeAlerts/StreamAnnotations-Backend/extension"
	"github.com/PlumeAlerts/StreamAnnotations-Backend/requests"
	"github.com/PlumeAlerts/StreamAnnotations-Backend/utilities"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"time"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(clientId string, clientSecret string, ownerId string, address string, port string, dbname string, username string, password string) {
	client := extension.Client{ClientID: clientId, ClientSecret: clientSecret, OwnerID: ownerId}

	utilities.Validate = validator.New()
	db.Connect(address, port, dbname, username, password)
	go InitializeSender(client)

	a.Router = mux.NewRouter()
	a.Router.Use(client.VerifyJWT)

	a.Router.Handle("/config", http.HandlerFunc(requests.GetConfig)).Methods("GET")
	a.Router.Handle("/config", http.HandlerFunc(requests.PutConfig)).Methods("PUT")

	a.Router.Handle("/annotation/config", http.HandlerFunc(requests.GetAnnotationConfig)).Methods("GET")
	a.Router.Handle("/annotation/config", http.HandlerFunc(requests.PutAnnotationConfig)).Methods("PUT")
	a.Router.Handle("/annotation/config", http.HandlerFunc(requests.PostAnnotationConfig)).Methods("POST")
	a.Router.Handle("/annotation/config", http.HandlerFunc(requests.DeleteAnnotationConfig)).Methods("DELETE")
}

func (a *App) Run() {
	c := cors.New(cors.Options{
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	handler := c.Handler(a.Router)
	//TODO Look into making this error out properly
	log.Fatal(http.ListenAndServe(":8000", handler))
	err := db.DB.Close()
	if err != nil {
		log.Fatal(err)
	}
}
func InitializeSender(client extension.Client) {

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				triggerAnnotations(client)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func triggerAnnotations(client extension.Client) {
	var users []db.User
	db.DB.Where("last_triggered < now() - interval '1 second' * users.annotation_interval").Order("last_triggered").Find(&users)

	for i := 0; i < len(users); i++ {
		user := users[i]

		//TODO Should fix the count issue
		count := 0
		db.DB.Where("user_id = ?", user.Id).Find(&db.Annotation{}).Count(&count)
		if count == 0 {
			continue
		}

		var annotation db.Annotation
		db.DB.Where("id = (select min(id) from annotation where user_id = ? and id > ?)", user.Id, user.LastAnnotationId).Find(&annotation)
		if annotation.Id == 0 {
			db.DB.Where("id = (select min(id) from annotation where user_id = ? and id > ?)", user.Id, 0).Find(&annotation)
		}
		b, _ := json.Marshal(annotation)

		err := client.PostPubSubMessage(user.Id, string(b))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(user.Id)
		}
		db.DB.Save(&db.User{Id: user.Id, LastTriggered: time.Now(), LastAnnotationId: annotation.Id})
	}
}
