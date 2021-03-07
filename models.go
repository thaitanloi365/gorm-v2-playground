package main

// Model model
type Model struct {
	ID        string    `gorm:"primaryKey" json:"-"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt DeletedAt `json:"-"`
}

// User user
type User struct {
	Model
	Email string `json:"email"`
	Phone string `json:"phone"`
}
