package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	ConversationID uint           `json:"conversation_id" gorm:"not null"`
	Conversation   Conversation   `json:"-" gorm:"foreignKey:ConversationID"`
	SenderID       uint           `json:"sender_id" gorm:"not null"`
	SenderType     string         `json:"sender_type" gorm:"not null"` // "user" หรือ "pharmacy"
	Content        string         `json:"content" gorm:"not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}