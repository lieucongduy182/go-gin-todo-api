package models

import "time"

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null; index" json:"user_id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	Completed   bool      `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `gorm:"foreignKey:UserID" json:"-"`
}

type CreateTaskInput struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description"`
}

type UpdateTaskInput struct {
	Title       string `json:"title" binding:"omitempty,min=1,max=255"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	Completed   *bool  `json:"completed" binding:"omitempty"`
}
