package main

import (
	"errors"
	"log"
	"plugin"
	"time"

	"github.com/acanoe/newsbytes-api-go/controllers"
	"github.com/acanoe/newsbytes-api-go/middlewares"
	"github.com/acanoe/newsbytes-api-go/models"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"gorm.io/gorm/clause"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Define groups
	auth := r.Group("/auth")

	// Define routes
	r.GET("/", middlewares.JwtAuthMiddleware(), controllers.GetStories)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	auth.POST("/register", controllers.Register)
	auth.POST("/login", controllers.Login)

	return r
}

func loadSource(path string) (models.NewsSource, error) {
	// load plugin
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	// get source
	s, err := p.Lookup("Source")
	if err != nil {
		return nil, err
	}

	// check if source is a NewsSource
	source, ok := s.(models.NewsSource)
	if !ok {
		return nil, errors.New("source is not a NewsSource")
	}

	return source, nil
}

func updateNews() {
	sources := []string{
		"./sources/progscrape.so",
	}

	for _, sourcePath := range sources {
		source, err := loadSource(sourcePath)
		if err != nil {
			log.Printf("Error loading source %s: %v", sourcePath, err)
			continue
		}

		log.Printf("Getting news from source: %s", sourcePath)
		newsStories, err := source.GetNews()
		if err != nil {
			log.Printf("Error getting news from source %s: %v", sourcePath, err)
			continue
		}

		result := models.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&newsStories)
		if result.Error != nil {
			log.Printf("Error writing news stories to the database: %v", result.Error)
			continue
		}
		log.Printf("Updated %d stories from source: %s", result.RowsAffected, sourcePath)
	}
}

func main() {
	// Set up the router
	r := setupRouter()

	// Connect to the database
	models.ConnectDatabase()

	// Update news manually
	// updateNews()

	// Set schedule for updating stories
	s := gocron.NewScheduler(time.Local)
	s.EveryRandom(5, 10).Hours().Do(updateNews)

	// Start the server and scheduler
	r.Run()
	s.StartAsync()
}
