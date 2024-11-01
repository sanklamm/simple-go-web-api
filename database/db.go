package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"

	"github.com/sanklamm/simple-go-web-api/models"
)

var DB *gorm.DB
var TestDB bool = false

func ConnectDatabase() {
	var database *gorm.DB
	var err error

	if TestDB {
		database, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	} else {
		database, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("Could not connect to the database!", err)
	}

	database.AutoMigrate(&models.Product{}, &models.User{})

	DB = database
}
