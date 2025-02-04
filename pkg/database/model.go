package database

import "time"

type Message struct {
	ID        uint      `gorm:"primaryKey"`
	To        string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	Sent      bool      `gorm:"default:false"`
	SentAt    time.Time `gorm:"default:null"`
	MessageID string    `gorm:"default:null"`
}
