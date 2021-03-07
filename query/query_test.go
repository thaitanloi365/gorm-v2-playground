package query

import (
	"encoding/json"
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// Model model
type Model struct {
	ID        uint  `gorm:"primaryKey" json:"id"`
	UpdatedAt int64 `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt int64 `gorm:"autoCreateTime" json:"created_at"`
}

type Profile struct {
	Model
	Avatar string
}

// User user
type User struct {
	Model
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	ProfileID string   `json:"profile_id"`
	Profile   *Profile `json:"profile"`
}

func TestQuery(t *testing.T) {
	var err error
	db, err = gorm.Open(sqlite.Open("test_pagination.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = db.Debug()

	New(db.Raw("SELECT users.id AS id, users.profile_id AS profile_id, profiles.id AS profile_id, profiles.created_at AS profile_created_at, profiles.avatar AS profile_avatar FROM users LEFT JOIN profiles ON profiles.id = users.profile_id")).PaginateFunc(func(db *gorm.DB) (records interface{}, err error) {
		type UserAlias struct {
			User    *User    `gorm:"embedded"`
			Profile *Profile `gorm:"embedded;embeddedPrefix:profile_"`
		}

		var users []*User
		rows, err := db.Rows()
		if err != nil {
			return users, err
		}

		for rows.Next() {
			var alias UserAlias
			var err = db.ScanRows(rows, &alias)
			if err != nil {
				fmt.Println(err)
				continue
			}

			data, err := json.Marshal(&alias)
			if err != nil {
				panic(err)
			}

			fmt.Println("data", string(data))
		}

		fmt.Println("done")
		return users, nil
	})

}

func mockData() {
	if db.Migrator().HasTable(&User{}) {
		db.Migrator().DropTable(&User{})
	}

	if db.Migrator().HasTable(&Profile{}) {
		db.Migrator().DropTable(&Profile{})
	}

	db.AutoMigrate(&Profile{})
	db.AutoMigrate(&User{})
	var createUsers = []*User{}
	// Seed database
	for i := 0; i < 20; i++ {
		var user = User{
			Email: fmt.Sprintf("user_%d@test.com", i),
			Phone: "+12345678910",
			Profile: &Profile{
				Avatar: fmt.Sprintf("avatar_%d", i),
			},
		}
		createUsers = append(createUsers, &user)
	}

	var err = db.Create(&createUsers).Error
	if err != nil {
		panic(err)
	}

	var users []*User
	p, err := New(db.Find(&users)).Where("email <> ?", "").Page(2).Limit(20).Paginate(&users)
	if err != nil {
		panic(err)
	}

	printJSON(p)

	err = New(db.Debug().Raw("SELECT * FROM users")).Scan(&users)
	if err != nil {
		panic(err)
	}
	printJSON(users)

	var profiles []*Profile
	err = New(db.Debug().Raw("SELECT * FROM profiles")).Scan(&profiles)
	if err != nil {
		panic(err)
	}

	printJSON(profiles)

}
func printJSON(in interface{}) {
	data, _ := json.MarshalIndent(&in, "", "    ")
	fmt.Println(string(data))
}
