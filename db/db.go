package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	"os"
)

type InfoBox struct {
	ID              uint   `gorm:"primary_key" json:"id"`
	Text            string `json:"text"`
	TextSize        int    `json:"textSize" validate:"required,numeric,gte=1,lte=7"`
	URL             string `json:"url" validate:"url"`
	Icon            string `json:"icon" validate:"alphanumunicode"`
	IconColor       string `json:"iconColor" validate:"hexcolor"`
	TextColor       string `json:"textColor" validate:"hexcolor"`
	BackgroundColor string `json:"backgroundColor" validate:"hexcolor"`
	Intervals       int    `json:"intervals" validate:"required,numeric,gte=1,lte=120"`
	UserId          string
}

var DB *gorm.DB

func Connect() {
	var err error
	env := func(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}
	user := env("MYSQL_USER", "root")
	pass := env("MYSQL_PASS", "")
	addr := env("MYSQL_ADDR", "localhost:3306")
	dbname := env("MYSQL_DBNAME", "infoboxes")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true", user, pass, addr, dbname)
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic("failed to connect database")
	}
	//defer DB.Close()
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
