package query

import "gorm.io/gorm"

type builder struct {
	db          *gorm.DB
	rawSQL      *gorm.DB
	whereValues []interface{}
	statement   *gorm.Statement
	limit       int
	page        int
}

// Builder interface
type Builder interface {
	Where(query interface{}, value ...interface{}) Builder
}
