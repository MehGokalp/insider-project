package database

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	// add all models here

	return db.AutoMigrate(&Message{})
}
