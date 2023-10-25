package models

import (
	"errors"
	"log"
	"plugin"

	"gorm.io/gorm/clause"
)

type News struct {
	Tags    []string `json:"tags"`
	Stories []Story  `json:"stories"`
}

type NewsSource interface {
	GetName() string
	GetNews() ([]Story, error)
}

func LoadSource(path string) (NewsSource, error) {
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
	source, ok := s.(NewsSource)
	if !ok {
		return nil, errors.New("source is not a NewsSource")
	}

	return source, nil
}

func UpdateNews(sources []string) {
	for _, sourcePath := range sources {
		source, err := LoadSource(sourcePath)
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

		query := DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&newsStories)
		if query.Error != nil {
			log.Printf("Error writing news stories to the database: %v", query.Error)
			continue
		}
		log.Printf("Updated %d stories from source: %s", query.RowsAffected, sourcePath)
	}
}
