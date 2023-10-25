package controllers

import (
	"net/http"
	"strings"

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
	models.DB.Where("source IN (?)", selectedSources).Order("date desc").Find(&stories)

	// Get all the tags from the stories
	var rawTags []string
	models.DB.Model(&stories).Pluck("tags", &rawTags)

	for _, tag := range rawTags {
		tags = append(tags, strings.Split(tag, ",")...)
	}

	tags = utils.RemoveDuplicateStr(tags)

	c.JSON(http.StatusOK, gin.H{
		"tags":    tags,
		"stories": stories,
	})
}
