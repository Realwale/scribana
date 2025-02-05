package models

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Title      string         `gorm:"not null" json:"title"`
	Slug       string         `gorm:"unique;not null" json:"slug"`
	Content    string         `gorm:"type:text" json:"content"`
	ImageURL   string         `json:"image_url"`
	AuthorID   uint           `json:"author_id"`
	Author     User           `json:"author"`
	CategoryID uint           `json:"category_id"`
	Category   Category       `json:"category"`
	Comments   []Comment      `json:"comments,omitempty"`
	Likes      int            `gorm:"default:0" json:"likes"`
}
