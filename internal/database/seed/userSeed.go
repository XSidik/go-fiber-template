package seed

import (
	"github.com/XSidik/go-fiber-template/internal/models"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	password, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 14)
	users := []models.UserModel{
		{UserName: "admin", Password: string(password)},
	}

	for _, user := range users {
		var existing models.UserModel
		if err := db.Where("user_name = ?", user.UserName).First(&existing).Error; err != nil {
			if err := db.Create(&user).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
