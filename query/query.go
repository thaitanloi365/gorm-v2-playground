package query

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// New builder
func New(db *gorm.DB) Builder {
	return &builder{
		db:          db,
		whereValues: []interface{}{},
		statement:   &gorm.Statement{DB: db, Clauses: map[string]clause.Clause{}},
	}
}

func (b *builder) Where(query interface{}, value ...interface{}) Builder {
	// whereValues = append(whereValues,)
	b.statement.Where(query, value)
	fmt.Println("where", b.statement.SQL.String(), b.statement.Vars)
	var clause = clause.Where{}
	return b
}
