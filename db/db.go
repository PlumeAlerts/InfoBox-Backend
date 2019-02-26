package db

import (
	"fmt"
	"github.com/PlumeAlerts/StreamAnnotations-Backend/extension"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
	"time"
)

type User struct {
	Id                 string    `gorm:"primary_key" json:"id"`
	LastTriggered      time.Time `json:"last_triggered"`
	LastAnnotationId   int       `json:"last_annotation_id"`
	AnnotationInterval int       `json:"annotation_interval" validate:"required,numeric,gte=1,lte=120"`
}

type Annotation struct {
	Id              int    `gorm:"primary_key" json:"id"`
	Text            string `json:"text"`
	TextSize        int    `json:"textSize" validate:"required,numeric,gte=1,lte=7"`
	URL             string `json:"url" validate:"url"`
	Icon            string `json:"icon" validate:"required,alphanumunicode"`
	IconColor       string `json:"iconColor" validate:"hexcolor"`
	TextColor       string `json:"textColor" validate:"hexcolor"`
	BackgroundColor string `json:"backgroundColor" validate:"hexcolor"`
	UserId          string `json:"user_id"`
}

var DB *gorm.DB

func Connect(address string, port string, dbname string, username string, password string) {
	var err error
	DB, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", address, port, dbname, username, password))
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	Migrate()
}

func Migrate() {
	driver, _ := postgres.WithInstance(DB.DB(), &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		panic(err)
	}
	m.Steps(1)
}

func GetUserIdOrCreate(r *http.Request) string {
	userId := r.Context().Value(extension.ChannelIDKey).(string)

	var user = User{}
	user.AnnotationInterval = 15
	DB.Where(&User{Id: userId}).FirstOrCreate(&user)
	return userId
}
