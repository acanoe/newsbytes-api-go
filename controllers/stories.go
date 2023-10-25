package controllers

import (
	"net/http"

	"github.com/acanoe/newsbytes-api-go/models"
	"github.com/acanoe/newsbytes-api-go/utils"
	"github.com/gin-gonic/gin"
)

func GetStories(c *gin.Context) {
	var selectedSources []string
	var stories []models.Story
	var tags []string

	// Get user's selected sources
	var user models.User
	if err := models.DB.Model(&models.User{}).Where("id = ?", c.GetUint("user_id")).Preload("UserPreferences").First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	selectedSources = user.UserPreferences.SourceSelection

	// If no sources are selected, return none
	if len(selectedSources) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No sources selected",
		})
		return
	}

	// Get stories from the selected sources
	models.DB.
		Where("source IN (?)", selectedSources).
		Order("date desc").
		Find(&stories)

	// Extract tags from stories
	for _, story := range stories {
		tags = append(tags, story.Tags...)
	}

	c.JSON(http.StatusOK, gin.H{
		"tags":    utils.RemoveDuplicateStr(tags),
		"stories": stories,
	})
}

func GetStoriesBySource(c *gin.Context) {
	source := c.Param("source")

	var stories []models.Story
	var tags []string

	models.DB.
		Where("source = ?", source).
		Order("date desc").
		Find(&stories)

	// if no stories are found, return 404
	if len(stories) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No stories found",
		})
		return
	}

	// Extract tags from stories
	for _, story := range stories {
		tags = append(tags, story.Tags...)
	}

	c.JSON(http.StatusOK, gin.H{
		"tags":    utils.RemoveDuplicateStr(tags),
		"stories": stories,
	})
}
