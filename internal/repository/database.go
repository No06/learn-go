package repository

import (
	"fmt"
	"log"

	"hinoob.net/learn-go/internal/config"
	"hinoob.net/learn-go/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB initializes the database connection and migrates the schema
func InitDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.DBName,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.SSLMode,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully.")

	// Auto-migrate the schema
	migrateDatabase()
}

func migrateDatabase() {
	log.Println("Running database migrations...")
	err := DB.AutoMigrate(
		&model.User{},
		&model.Assignment{},
		&model.Submission{},
		&model.Comment{},
		&model.SubmissionFile{},
		&model.Class{},
		&model.TimeSlot{},
		&model.Course{},
		&model.Message{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed.")
}
