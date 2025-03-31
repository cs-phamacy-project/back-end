package models

import (
	"time"

	"gorm.io/gorm"
)

type Conversation struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null"`
	User       User           `json:"-" gorm:"foreignKey:UserID"`
	PharmacyID uint           `json:"pharmacy_id" gorm:"not null"`
	Pharmacy   Pharmacy       `json:"-" gorm:"foreignKey:PharmacyID"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	Messages   []Message      `json:"messages,omitempty" gorm:"foreignKey:ConversationID"`
}