package repository

import "gorm.io/gorm"

func userSummary(db *gorm.DB) *gorm.DB {
	return db.Select("id", "name", "email")
}
