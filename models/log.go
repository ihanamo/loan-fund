package models

import "time"

// Log model to track user actions such as loan issuance, repayments, deposits
type Log struct {
	ID            uint        `json:"log_id" gorm:"primaryKey;autoIncrement"`
	UserID        uint        `json:"user_id" gorm:"not null"`
	User          User        `gorm:"foreignKey:UserID"`
	TransactionID *uint       `json:"transaction_id" gorm:"not null"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID"`
	Type          string      `json:"type" gorm:"not null"` // "loan", "repayment", "deposit"
	Message       string      `json:"message"`
	CreatedAt     time.Time   `json:"created_at"`
}
