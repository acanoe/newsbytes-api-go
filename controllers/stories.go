package controllers

import (
	"net/http"
	"strings"

	"github.com/acanoe/newsbytes-api-go/models"
	"github.com/acanoe/newsbytes-api-go/utils"
	"github.com/gin-gonic/gin"
)

func GetStories(c *gin.Context) {
	var stories []models.Story
	var tags []string

	models.DB.Find(&stories)

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
