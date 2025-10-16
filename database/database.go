package database

import (
	"fmt"
	"log"

	"github.com/lieucongduy182/go-gin-todo-api/config"
	"github.com/lieucongduy182/go-gin-todo-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	cfg := config.AppConfig

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	fmt.Println("✅ Connected to Database!")
	// Migrate the schema
	if err := DB.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		log.Fatal("Failed to migrate database", err)
	}

	fmt.Println("✅ Database migrated successfully!")
}

func GetDB() *gorm.DB {
	return DB
}
