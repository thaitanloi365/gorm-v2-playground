package main

import "gorm.io/gorm"

// Model model
type Model struct {
	gorm.Model
	ID        string         `gorm:"primaryKey" json:"-"`
	UpdatedAt int64          `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt int64          `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// User user
type User struct {
	Model
	Email string `json:"email"`
	Phone string `json:"phone"`
}
