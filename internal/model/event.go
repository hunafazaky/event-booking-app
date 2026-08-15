package model

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string    `json:"name" binding:"required" `
	Description string    `json:"description" binding:"required" `
	Location    string    `json:"location" binding:"required" `
	UserID      int       `json:"user_id"`
	User        User      `json:"-" gorm:"foreignKey:user_id"`
	DateTime    time.Time `json:"date_time" binding:"required" `
}

var events []Event = []Event{}

func (e Event) Save() {
	events = append(events, e)
}

func GetAllEvents() []Event {
	return events
}
