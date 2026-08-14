package models

import (
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string `binding:"required" json:"name"`
	Description string `binding:"required" json:"description"`
	Location    string `binding:"required" json:"location"`
	UserID      int    `json:"user_id"`
}

var events []Event = []Event{}

func (e Event) Save() {
	events = append(events, e)
}

func GetAllEvents() []Event {
	return events
}
