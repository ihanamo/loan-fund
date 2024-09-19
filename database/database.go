package database

import (
	"loan-fund/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("fund-loan.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect db", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Loan{}, &models.Transaction{})
	if err != nil {
		log.Fatal("Failed to migrate db", err)
	}
}
