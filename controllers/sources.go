package controllers

import (
	"net/http"

	"github.com/acanoe/newsbytes-api-go/models"
	"github.com/gin-gonic/gin"
)

func GetAvailableSources(c *gin.Context) {
	var availableSources []string

	query := models.DB.Model(&models.Story{}).Select("source").Distinct("source").Pluck("source", &availableSources)
	if query.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": query.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sources": availableSources,
	})
}

type UserSourcesInput struct {
	Sources []string `json:"sources"`
}

func SetUserSources(c *gin.Context) {
	var input UserSourcesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	query := models.DB.Model(&models.User{}).Where("id = ?", c.GetUint("user_id")).First(&user)
	if query.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": query.Error.Error(),
		})
		return
	}

	// remove old preferences
	models.DB.Model(&user).Association("UserPreferences").Unscoped().Clear()

	// replace with new preferences
	userPreferences := &models.UserPreferences{
		SourceSelection: input.Sources,
	}
	if err := models.DB.Model(&user).Association("UserPreferences").Replace(userPreferences); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "sources updated",
	})
}
