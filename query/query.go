package query

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/utils"
)

// WhereFunc where func
type WhereFunc = func(b Builder)

// New builder
func New(db *gorm.DB) Builder {
	return &builder{
		db:        db,
		values:    []interface{}{},
		rawSQL:    "",
		statement: &gorm.Statement{DB: db, Clauses: map[string]clause.Clause{}},
	}
}

func (b *builder) Raw(sql string) Builder {
	b.rawSQL = sql
	return b
}

// Exec exec
func (b *builder) Scan(dest interface{}) (err error) {
	q, vars := b.Prepare()
	err = b.db.Raw(q, vars...).Scan(dest).Error

	return
}

// Exec exec
func (b *builder) Prepare() (string, []interface{}) {
	var listClause = []string{}
	for key := range b.statement.Clauses {
		listClause = append(listClause, key)
	}
	b.statement.Build(listClause...)

	var vars = b.statement.Vars
	var rawSQL = b.statement.SQL.String()

	if b.rawSQL != "" {
		rawSQL = fmt.Sprintf("%s %s", b.rawSQL, rawSQL)
	}

	return rawSQL, vars
}

// Where where
func (b *builder) Where(query interface{}, value ...interface{}) Builder {
	if expressions := b.statement.BuildCondition(query, value); len(expressions) > 0 {
		b.statement.AddClause(clause.Where{Exprs: expressions})
	}

	return b
}

// WhereFunc where func
func (b *builder) WhereFunc(f WhereFunc) Builder {
	f(b)
	return b
}

// Order order
func (b *builder) Order(value interface{}) Builder {
	switch v := value.(type) {
	case clause.OrderByColumn:
		b.statement.AddClause(clause.OrderBy{
			Columns: []clause.OrderByColumn{v},
		})
	default:
		b.statement.AddClause(clause.OrderBy{
			Columns: []clause.OrderByColumn{{
				Column: clause.Column{Name: fmt.Sprint(value), Raw: true},
			}},
		})
	}

	return b
}

// Order order
func (b *builder) Group(name string) Builder {
	var fields = strings.FieldsFunc(name, utils.IsChar)

	b.statement.AddClause(clause.GroupBy{
		Columns: []clause.Column{{Name: name, Raw: len(fields) != 1}},
	})

	return b
}
