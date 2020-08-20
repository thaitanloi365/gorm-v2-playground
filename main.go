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

	// Migrate the schema
	db.AutoMigrate(&User{})

	db.Callback().Create().Before("gorm:save_before_associations").Register("app:update_xid_when_create", func(db *gorm.DB) {
		if id, ok := db.Statement.Get("ID"); ok {
			if id, ok := id.(string); ok {
				if id == "" {
					db.Statement.SetColumn("ID", xid.New().String())
				}
			}
		}
	})

	for i := 0; i < 20; i++ {
		var user = User{
			Email: fmt.Sprintf("user_%d@test.com", i),
			Phone: "+12345678910",
		}

		var err = db.Create(&user).Error
		if err != nil {
			fmt.Println(err)
		}
	}

}
