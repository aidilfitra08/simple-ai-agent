package database

import (
	"log"

	"github.com/aidilfitra08/simple-ai-agent/config"
	"github.com/aidilfitra08/simple-ai-agent/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect establishes a database connection
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()
	log.Printf("Connecting to DB provider: %s", cfg.DBProvider)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	log.Println("Database migration completed")
	return nil
}
