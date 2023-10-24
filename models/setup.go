package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})

	if err != nil {
		panic("Cannot connect to database")
	}

	err = database.AutoMigrate(&Story{})

	if err != nil {
		return
	}

	DB = database
}
