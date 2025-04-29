package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"todo-api-2.0/models"
)

var DB *gorm.DB

// InitDB initializes the SQLite database connection
func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	DB.AutoMigrate(&models.Todo{})
}
