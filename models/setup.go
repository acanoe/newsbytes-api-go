package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func autoMigrateStructs(db *gorm.DB, structs ...interface{}) error {
	for _, str := range structs {
		if err := db.AutoMigrate(str); err != nil {
			return err
		}
	}
	return nil
}

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})

	if err != nil {
		panic("Cannot connect to database")
	}

	err = autoMigrateStructs(database, &Story{}, &User{})
	if err != nil {
		return
	}

	DB = database
}
