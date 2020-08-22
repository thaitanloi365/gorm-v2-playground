package query

import (
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestQuery(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test_pagination.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = db.Debug()

	// Model model
	type Model struct {
		ID        string `gorm:"primaryKey" json:"-"`
		UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at"`
		CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	}

	// User user
	type User struct {
		Model
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	if db.Migrator().HasTable(&User{}) {
		db.Migrator().DropTable(&User{})
	}

	db.AutoMigrate(&User{})

	var createUsers = []*User{}
	// Seed database
	for i := 0; i < 20; i++ {
		var user = User{
			Model: Model{
				ID: fmt.Sprintf("user_%d", i),
			},
			Email: fmt.Sprintf("user_%d@test.com", i),
			Phone: "+12345678910",
		}
		createUsers = append(createUsers, &user)
	}

	err = db.Create(&createUsers).Error
	if err != nil {
		panic(err)
	}
	var b = New(db).Where("email = ?", "1234")
	fmt.Println(b)
}
