package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("fund-loan.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect db", err)
	}
}
