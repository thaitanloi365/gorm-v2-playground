package query

import (
	"encoding/json"
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

	var users []*User
	p, err := New(db.Find(&users)).Where("email <> ?", "").Page(2).Limit(20).Paginate(&users)
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

	// var u User
	// err = db.Raw(`SELECT * FROM users WHERE email = ? LIMIT 1`, "user_19@test.com").Scan(&u).Error
	// if err != nil {
	// 	panic(err)
	// }

	// db.Where("email = ?", "user_19@test.com").First(&u)
	// fmt.Println(u)

	// fmt.Println("r", r)
	// fmt.Println(result)
}
