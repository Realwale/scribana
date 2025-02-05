package models

import (
	"gorm.io/gorm"
	"time"
)

type Category struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Slug      string         `gorm:"unique;not null" json:"slug"`
	Posts     []Post         `json:"posts,omitempty"`
}
