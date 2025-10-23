package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Completed   bool           `gorm:"default:false; index" json:"completed"`
	Priority    string         `gorm:"default:'medium'" json:"priority"`
	DueDate     *time.Time     `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"omitempty,min=1,max=255"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	Completed   *bool  `json:"completed" binding:"omitempty"`
}

type TaskResponse struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int64       `json:"total_pages"`
}
