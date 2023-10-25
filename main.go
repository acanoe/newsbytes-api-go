package main

import (
	"time"

	"github.com/acanoe/newsbytes-api-go/controllers"
	"github.com/acanoe/newsbytes-api-go/middlewares"
	"github.com/acanoe/newsbytes-api-go/models"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Define groups
	auth := r.Group("/auth")
	sources := r.Group("/sources")

	// Define routes
	r.GET("/", middlewares.JwtAuthMiddleware(), controllers.GetStories)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	auth.POST("/register", controllers.Register)
	auth.POST("/login", controllers.Login)

	sources.GET("/", controllers.GetAvailableSources)
	sources.POST("/", middlewares.JwtAuthMiddleware(), controllers.SetUserSources)
	sources.GET("/:source", controllers.GetStoriesBySource)

	return r
}

func main() {
	// Set up the router
	r := setupRouter()

	// Connect to the database
	models.ConnectDatabase()

	// News sources
	sources := []string{
		"./sources/progscrape.so",
	}

	// Update news manually
	// models.UpdateNews(sources)

	// Set schedule for updating stories
	s := gocron.NewScheduler(time.Local)
	s.EveryRandom(5, 10).Hours().Do(models.UpdateNews, sources)

	// Start the server and scheduler
	r.Run()
	s.StartAsync()
}
