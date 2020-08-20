package main

// Model model
type Model struct {
	ID string `gorm:"primaryKey" json:"-"`
}

// User user
type User struct {
	Model
	Email string `json:"email"`
	Phone string `json:"phone"`
}
