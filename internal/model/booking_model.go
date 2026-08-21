package model

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	BookingCode string `json:"booking_code"`
	Phone       string `json:"phone"`
	UserID      uint   `json:"user_id"`
	User        User   `json:"user" gorm:"foreignKey:UserID"`
	EventID     uint   `json:"event_id"`
	Event       Event  `json:"event" gorm:"foreignKey:EventID"`
}
