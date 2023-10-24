package models

type News struct {
	Tags    []string `json:"tags"`
	Stories []Story  `json:"stories"`
}

type NewsSource interface {
	GetNews() ([]Story, error)
}
