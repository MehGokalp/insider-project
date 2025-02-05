package repository

import (
	"github.com/mehgokalp/insider-project/pkg/database"
	"gorm.io/gorm"
	"time"
)

type MessageRepository interface {
	List() ([]database.Message, error)
	GetUnsentMessages(limit int) ([]database.Message, error)
	UpdateSentStatus(message database.Message) error
}

type DBMessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *DBMessageRepository {
	return &DBMessageRepository{db: db}
}

func (r *DBMessageRepository) List() ([]database.Message, error) {
	var messages []database.Message
	if err := r.db.Where("sent = ?", true).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *DBMessageRepository) GetUnsentMessages(limit int) ([]database.Message, error) {
	var messages []database.Message
	if err := r.db.Where("sent = ?", false).Order("id ASC").Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *DBMessageRepository) UpdateSentStatus(message database.Message) error {
	return r.db.Model(&message).Updates(
		map[string]interface{}{
			"sent":       true,
			"sent_at":    time.Now(),
			"message_id": message.MessageId,
		}).Error
}
