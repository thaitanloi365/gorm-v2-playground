package query

import "gorm.io/gorm"

type builder struct {
	db        *gorm.DB
	rawSQL    string
	values    []interface{}
	statement *gorm.Statement
}

// Builder interface
type Builder interface {
	Raw(sql string) Builder
	Where(query interface{}, value ...interface{}) Builder
	WhereFunc(WhereFunc) Builder
	Order(order interface{}) Builder
	Scan(dest interface{}) error
	Prepare() (string, []interface{})
	Group(name string) Builder
}
