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

type User struct {
	ID              string `gorm:"primary_key"`
	InfoboxInterval int
}

type InfoBox struct {
	ID uint `gorm:"primary_key"`

	Title           string `json:"title" validate:"required,text"`
	TextSize        int    `json:"textSize" validate:"required,numeric,gte=1,lte=7"`
	URL             string `json:"url" validate:"url"`
	Icon            string `json:"icon" validate:"alphanumunicode"`
	IconColor       string `json:"iconColor" validate:"hexcolor"`
	TextColor       string `json:"textColor" validate:"hexcolor"`
	BackgroundColor string `json:"backgroundColor" validate:"hexcolor"`

	UserId string
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
	user := env("MYSQL_TEST_USER", "root")
	pass := env("MYSQL_TEST_PASS", "")
	prot := env("MYSQL_TEST_PROT", "tcp")
	addr := env("MYSQL_TEST_ADDR", "localhost:3306")
	dbname := env("MYSQL_TEST_DBNAME", "infoboxes")
	netAddr := fmt.Sprintf("%s(%s)", prot, addr)
	dsn := fmt.Sprintf("%s:%s@%s/%s?multiStatements=true", user, pass, netAddr, dbname)
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic("failed to connect database")
	}
	//defer DB.Close()

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

func GetUserOrCreate(userId string) (*gorm.DB, bool) {
	user := DB.FirstOrCreate(&User{ID: userId})
	if user.RowsAffected == 0 {
		user = DB.Create(&User{ID: userId, InfoboxInterval: 15})
		return user, true
	}
	return user, false
}
