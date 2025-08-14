package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;size:255" json:"email"`
	PasswordHash string `json:"-"`
	Name         string `json:"name"`
}
