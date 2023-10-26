package models

import (
	"fmt"
	"log"

	"github.com/acanoe/newsbytes-api-go/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Cannot load .env file")
	}

	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbName := utils.GetEnv("DB_NAME", "mydatabase")
	dbUser := utils.GetEnv("DB_USER", "myuser")
	dbPass := utils.GetEnv("DB_PASS", "mypassword")
	sslMode := utils.GetEnv("SSL_MODE", "disable")
	timeZone := utils.GetEnv("TIMEZONE", "Asia/Jakarta")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", dbHost, dbUser, dbPass, dbName, dbPort, sslMode, timeZone)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	err = autoMigrateStructs(database, &Story{}, &User{}, &UserPreferences{})
	if err != nil {
		return
	}

	DB = database
}
