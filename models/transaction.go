package models

import "time"

type Transaction struct {
	ID        uint      `json:"transaction-id" gorm:"primaryKey;autoIncrement"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
	LoanID    uint      `json:"loan_id" gorm:"not null"`
	Loan      Loan      `gorm:"foreignKey:LoanID"`
	Type      string    `json:"type" gorm:"not null"`
	Amount    float64   `json:"amount" gorm:"not null"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type Loan struct {
	ID            uint      `json:"loan-id" gorm:"primaryKey;autoIncrement"`
	UserID        uint      `json:"user_id"`
	User          User      `gorm:"foreignKey:UserID"`
	Amount        float64   `json:"amount"`
	InterestsRate float64   `json:"interests_rate"`
	DueDate       time.Time `json:"due_date"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
