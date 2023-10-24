package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/acanoe/newsbytes-api-go/models"
)

type Progscrape struct{}

// Filter stories by date
func filterStoriesByDate(stories []models.Story, date time.Time) []models.Story {
	filteredStories := []models.Story{}
	for _, story := range stories {
		if story.Date.Truncate(24 * time.Hour).Equal(date) {
			filteredStories = append(filteredStories, story)
		}
	}
	return filteredStories
}

func (p Progscrape) GetNews() ([]models.Story, error) {
	url := "https://progscrape.com/feed.json"

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	var newsData models.News
	err = json.NewDecoder(res.Body).Decode(&newsData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	filteredStories := filterStoriesByDate(newsData.Stories, today)

	return filteredStories, nil
}

var Source Progscrape
