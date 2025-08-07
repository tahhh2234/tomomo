package models

import "time"

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"default:1" json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `gorm:"default:'pending'" json:"status"`
	DueDate     time.Time `json:"due_date"`
	Priority    int       `gorm:"default:0" json:"priority"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
