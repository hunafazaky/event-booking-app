package config

import (
	"fmt"
	"log"
	"os"

	"github.com/hunafazaky/event-booking-app/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func getDSN(dbURI, dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode string) string {
	if dbURI != "" {
		return dbURI
	}
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbSSLMode == "" {
		log.Fatal("Required database environment variables are missing.")
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode,
	)
}

func ConnectDB() {
	dbURI := os.Getenv("DB_URI")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB_NAME")
	dbSSLMode := os.Getenv("POSTGRES_SSLMODE")

	dsn := getDSN(dbURI, dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

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
