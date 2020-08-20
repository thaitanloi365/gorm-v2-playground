package main

import (
	"fmt"

	"github.com/rs/xid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = db.Debug()

	// Migrate the schema

	if db.Migrator().HasTable(&User{}) {
		db.Migrator().DropTable(&User{})
	}

	db.AutoMigrate(&User{})

	db.Callback().Create().Before("gorm:save_before_associations").Register("app:update_xid_when_create", func(db *gorm.DB) {
		var field = db.Statement.Schema.LookUpField("ID")
		if field != nil {
			if v, isZero := field.ValueOf(db.Statement.ReflectValue); isZero {
				if id, ok := v.(string); ok {
					if id == "" {
						field.Set(db.Statement.ReflectValue, xid.New().String())
					}
				}
			}
		}

	})

	for i := 0; i < 20; i++ {
		var user = User{
			Model: Model{
				ID: fmt.Sprintf("user_%d", i),
			},
			Email: fmt.Sprintf("user_%d@test.com", i),
			Phone: "+12345678910",
		}

		var err = db.FirstOrCreate(&user).Error
		if err != nil {
			fmt.Println(err)
		}
	}
	var users []*User
	db.Find(&users)
	for _, u := range users {
		fmt.Println(u)
	}
}
