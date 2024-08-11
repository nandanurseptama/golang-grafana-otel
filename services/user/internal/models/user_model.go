package models

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}

func (UserModel) TableName() string {
	return "users"
}
