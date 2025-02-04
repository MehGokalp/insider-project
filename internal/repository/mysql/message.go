package mysql

import (
	"github.com/mehgokalp/insider-project/internal/domain/mysql"
	"gorm.io/gorm"
	"time"
)

type MessageRepository interface {
	List() ([]mysql.Message, error)
	GetUnsentMessages(limit int) ([]mysql.Message, error)
	UpdateSentStatus(message mysql.Message) error
}

type DBMessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *DBMessageRepository {
	return &DBMessageRepository{db: db}
}

func (r *DBMessageRepository) List() ([]mysql.Message, error) {
	var messages []mysql.Message
	if err := r.db.Where("sent = ?", true).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *DBMessageRepository) GetUnsentMessages(limit int) ([]mysql.Message, error) {
	var messages []mysql.Message
	if err := r.db.Where("sent = ?", false).Order("id ASC").Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *DBMessageRepository) UpdateSentStatus(message mysql.Message) error {
	return r.db.Model(&message).Updates(
		map[string]interface{}{
			"sent":       true,
			"sent_at":    time.Now(),
			"message_id": message.MessageId,
		}).Error
}
