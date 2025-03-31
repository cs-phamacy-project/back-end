package models

import (
	"time"
)

// Payment โครงสร้างข้อมูลการชำระเงิน
type Payment struct {
	ID          uint      `gorm:"primaryKey;column:payment_id" json:"id"`
	OrderID     uint      `gorm:"not null;column:order_id" json:"order_id"`
	ChargeID    string    `gorm:"not null;column:charge_id" json:"charge_id"`
	Amount      float64   `gorm:"not null;column:amount" json:"amount"`
	Status      string    `gorm:"not null;column:status" json:"status"`
	PaymentType string    `gorm:"not null;column:payment_type" json:"payment_type"`
	QRCode      string    `gorm:"column:qr_code" json:"qr_code"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}