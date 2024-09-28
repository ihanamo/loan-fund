package handlers

import (
	"loan-fund/database"
	"loan-fund/models"
	"time"
)

func LogAction(userID uint, actionType, message string, transactionID *uint) error {
	logEntry := models.Log{
		UserID:        userID,
		Type:          actionType,
		Message:       message,
		TransactionID: transactionID,
		CreatedAt:     time.Now(),
	}

	if result := database.DB.Create(&logEntry); result.Error != nil {
		return result.Error
	}

	return nil
}
