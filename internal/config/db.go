package config

import (
	"log"
	"os"

	"github.com/hunafazaky/event-booking-app/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DB_URI")
	if dsn == "" {
		log.Fatal("Environment variable is not set.")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database. Err:", err)
	}

	if err := db.AutoMigrate(&model.Event{}, model.User{}); err != nil {
		log.Fatal("Failed to auto-migrate. Err:", err)
	}

	DB = db
	log.Println("Connected to the database successfully.")
}
