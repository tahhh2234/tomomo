package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uint       `gorm:"primaryKey"`
	UserID    uint       `gorm:"index;not null"`
	TokenHash string     `gorm:"size:64;uniqueIndex;not null"` // sha256 hex
	ExpiresAt time.Time  `gorm:"not null"`
	RevokedAt *time.Time `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
