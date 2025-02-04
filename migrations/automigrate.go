package migrations

import (
	"github.com/mehgokalp/insider-project/internal/domain/mysql"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	// add all models here

	return db.AutoMigrate(&mysql.Message{})
}
