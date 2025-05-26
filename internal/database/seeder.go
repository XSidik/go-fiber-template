package database

import (
	"fmt"

	"github.com/XSidik/go-fiber-template/internal/database/seed"

	"gorm.io/gorm"
)

func RunSeeders(db *gorm.DB) error {
	if err := seed.SeedUsers(db); err != nil {
		return fmt.Errorf("user seeder error: %w", err)
	}

	fmt.Println("✅ All seeders executed.")
	return nil
}
