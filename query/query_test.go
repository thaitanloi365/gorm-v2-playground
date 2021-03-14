package query

import (
	"encoding/json"
	"fmt"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CustomString string

var db *DBTest

type DBTest struct {
	*gorm.DB
}

func (db *DBTest) GetGorm() *gorm.DB {
	return db.DB
}

func (db *DBTest) WithGorm(gdb *gorm.DB) DB {
	db.DB = gdb
	return db
}

func (db *DBTest) SetDebug(debug bool) {
	if debug {
		db.DB = db.DB.Debug()
	}
}

// Model model
type Model struct {
	ID        uint  `gorm:"primaryKey" json:"id"`
	UpdatedAt int64 `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedAt int64 `gorm:"autoCreateTime" json:"created_at"`
}

type Profile struct {
	Model
	Avatar string `json:"avatar"`
}

type CreditCard struct {
	Model
	UserID uint   `json:"user_id"`
	Last4  string `json:"last_4"`
}

// User user
type User struct {
	Model
	Email       string        `json:"email"`
	Phone       string        `json:"phone"`
	ProfileID   string        `json:"profile_id"`
	Profile     *Profile      `json:"profile"`
	CreditCards []*CreditCard `json:"credit_cards"`
}

func TestQueryWhereNamed(t *testing.T) {
	initDB()
	New(db, "SELECT users.id AS user_id, users.profile_id AS user_profile_id, users.profile_id AS profile_id, profiles.id AS profile_id, profiles.created_at AS profile_created_at, profiles.avatar AS profile_avatar FROM users LEFT JOIN profiles ON profiles.id = users.profile_id").
		PagingFunc(func(db, rawSQL DB) (interface{}, error) {
			// type UserAlias struct {
			// 	User    *User    `gorm:"embedded;embeddedPrefix:user_"`
			// 	Profile *Profile `gorm:"embedded;embeddedPrefix:profile_"`
			// }

			var users []*User
			// var alias []*UserAlias
			var results []map[string]interface{}

			db.GetGorm().Model(&User{}).Find(&results)

			printJSON(results)
			// rows, err := db.Rows()
			// if err != nil {
			// 	return users, err
			// }

			// for rows.Next() {
			// 	var alias UserAlias
			// 	var err = db.ScanRows(rows, &alias)
			// 	if err != nil {
			// 		fmt.Println(err)
			// 		continue
			// 	}

			// 	data, err := json.Marshal(&alias)
			// 	if err != nil {
			// 		panic(err)
			// 	}

			// 	fmt.Println("data", string(data))
			// }

			// fmt.Println("done")
			return users, nil

		})
}
func TestQuery(t *testing.T) {
	initDB()

	var customString CustomString = "33"
	var customStrings = []CustomString{
		"1", "2",
	}

	var result = New(db, `
  SELECT u.*, (
    CASE
      WHEN u.id = @user_id OR u.email = @user_name OR u.email IN @user_names THEN @then_value
      ELSE FALSE
    END
  ) ,
  row_to_json(p) AS profile
  FROM users u LEFT JOIN profiles p ON p.id = u.profile_id
  `).
		// NamedMap(map[string]interface{}{
		// 	"user_id":    1,
		// 	"user_name":  &customString,
		// 	"user_names": customStrings,
		// 	"then_value": true,
		// }).
		Where("u.id < ?", 5).
		Where(map[string]interface{}{
			"user_id":    1,
			"user_name":  &customString,
			"user_names": customStrings,
			"then_value": true,
		}).
		Having("u.id > 0").
		OrderBy("u.id DESC").
		GroupBy("u.id, p.*").
		Page(1).
		Limit(20).
		PagingFunc(func(db, rawSQL DB) (interface{}, error) {
			type UserAlias struct {
				User    *User    `gorm:"embedded;embeddedPrefix:user_"`
				Profile *Profile `gorm:"embedded;embeddedPrefix:profile_"`
			}

			// var users []*User
			var results []map[string]interface{}

			// var where = map[string]interface{}{
			// 	"user_id": 10,
			// }
			// db.GetGorm().Model(&User{}).Where("id = ?", 1).Find(&results)
			rawSQL.GetGorm().Find(&results)

			return results, nil
		})

	printJSON(result)
}

func initDB() {
	var uri = "host=localhost user=postgres password= dbname=gorm_test port=5432 sslmode=disable"
	gdb, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = &DBTest{
		DB: gdb,
	}
	db.SetDebug(true)
}

func mockData() {
	if db.Migrator().HasTable(&User{}) {
		db.Migrator().DropTable(&User{})
	}

	if db.Migrator().HasTable(&Profile{}) {
		db.Migrator().DropTable(&Profile{})
	}

	if db.Migrator().HasTable(&CreditCard{}) {
		db.Migrator().DropTable(&CreditCard{})
	}

	db.AutoMigrate(&CreditCard{})
	db.AutoMigrate(&Profile{})
	db.AutoMigrate(&User{})
	var createUsers = []*User{}
	// Seed database
	for i := 1; i <= 20; i++ {
		var user = User{
			Email: fmt.Sprintf("user_%d@test.com", i),
			Phone: "+12345678910",
			Profile: &Profile{
				Avatar: fmt.Sprintf("avatar_%d", i),
			},
			CreditCards: []*CreditCard{
				{
					Last4: "1111",
				},
				{
					Last4: "1112",
				},
				{
					Last4: "1113",
				},
				{
					Last4: "1114",
				},
				{
					Last4: "1115",
				},
			},
		}
		createUsers = append(createUsers, &user)
	}

	var err = db.Create(&createUsers).Error
	if err != nil {
		panic(err)
	}

}
func printJSON(in interface{}) {
	data, _ := json.MarshalIndent(&in, "", "    ")
	fmt.Println(string(data))
}
