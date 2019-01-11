package db

import (
	"fmt"
	"github.com/PlumeAlerts/InfoBoxes-Backend/extension"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"time"
)

type User struct {
	ID              string    `gorm:"primary_key" json:"id"`
	Intervals       int       `json:"intervals" validate:"required,numeric,gte=1,lte=120"`
	LastTriggered   time.Time `json:"last_triggered"`
	LastInfoBoxesId uint      `json:"last_info_boxes_id"`
}

type InfoBox struct {
	ID              uint   `gorm:"primary_key" json:"id"`
	Text            string `json:"text"`
	TextSize        int    `json:"textSize" validate:"required,numeric,gte=1,lte=7"`
	URL             string `json:"url" validate:"url"`
	Icon            string `json:"icon" validate:"required,alphanumunicode"`
	IconColor       string `json:"iconColor" validate:"hexcolor"`
	TextColor       string `json:"textColor" validate:"hexcolor"`
	BackgroundColor string `json:"backgroundColor" validate:"hexcolor"`
	UserId          string
}

var DB *gorm.DB

func env(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Connect() {
	var err error
	user := env("MYSQL_USER", "root")
	pass := env("MYSQL_PASS", "")
	addr := env("MYSQL_ADDR", "localhost:3306")
	dbname := env("MYSQL_DBNAME", "infoboxes")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true&parseTime=true&loc=Local", user, pass, addr, dbname)
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic("failed to connect database")
	}
	Migrate()
}
func Migrate() {
	driver, _ := mysql.WithInstance(DB.DB(), &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		panic(err)
	}
	m.Steps(1)
}

func GetUserOrCreate(r *http.Request) string {
	userId := r.Context().Value(extension.ChannelIDKey).(string)

	var user = User{}
	user.Intervals = 15
	DB.Where(&User{ID: userId}).FirstOrCreate(&user)
	return userId
}
