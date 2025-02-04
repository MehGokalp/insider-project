package repository

import (
	"github.com/mehgokalp/insider-project/pkg/database"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) List() ([]database.Message, error) {
	var messages []database.Message
	if err := r.db.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
