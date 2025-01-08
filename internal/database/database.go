package database

import (
	"fmt"
	"log"

	"github.com/XSidik/go-fiber-template/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg config.Config) {
	var err error
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.GetConfig().DBHost, config.GetConfig().DBPort, config.GetConfig().DBUser, config.GetConfig().DBPassword, config.GetConfig().DBName)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
}
