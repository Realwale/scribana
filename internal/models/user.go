package models

import (
	"gorm.io/gorm"
	"time"
)

type Role string

const (
	AdminRole  Role = "admin"
	AuthorRole Role = "author"
	ReaderRole Role = "reader"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Username  string         `gorm:"unique;not null" json:"username"`
	Password  string         `json:"-"`
	Role      Role           `gorm:"type:varchar(20);default:'reader'" json:"role"`
	Posts     []Post         `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`
	Comments  []Comment      `gorm:"foreignKey:UserID" json:"comments,omitempty"`
}
