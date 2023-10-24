package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type Tag []string

func (t *Tag) Scan(src any) error {
	bytes, ok := src.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", src)
	}
	*t = Tag(strings.Split(string(bytes), ","))
	return nil
}
func (t Tag) Value() (driver.Value, error) {
	if len(t) == 0 {
		return nil, nil
	}
	return strings.Join(t, ","), nil
}

type Story struct {
	ID       uint      `json:"-" gorm:"primaryKey"`
	Source   string    `json:"source"`
	Date     time.Time `json:"date"`
	Title    string    `json:"title"`
	Href     string    `json:"href" gorm:"unique"`
	Reddit   string    `json:"reddit"`
	Lobsters string    `json:"lobsters"`
	Hnews    string    `json:"hnews"`
	Tags     Tag       `json:"tags" gorm:"type:text"`
}
