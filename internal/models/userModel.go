package models

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	UserName string `gorm:"unique" json:"user_name"`
	Password string `json:"-"`
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&UserModel{})
}
