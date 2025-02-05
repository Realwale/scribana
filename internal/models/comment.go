package models

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Content   string         `gorm:"type:text" json:"content"`
	PostID    uint           `json:"post_id"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user"`
}
