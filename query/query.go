package query

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/utils"
)

// New builder
func New(db *gorm.DB) Builder {
	return &builder{
		db:     db,
		logger: log.New(os.Stdout, "go-query", log.Lshortfile),
	}
}

// Exec exec
func (b *builder) Scan(dest interface{}) (err error) {
	b.db.Statement.Build("SELECT", "FROM", "WHERE", "ORDER BY", "LIMIT", "GROUP BY", "FOR")
	var sql = b.db.Statement.SQL.String()
	var vars = b.db.Statement.Vars
	err = b.db.Raw(sql, vars...).Scan(dest).Error

	return
}

func (b *builder) SkipLimitOffset(skip bool) Builder {
	b.skipLimitOffset = skip
	return b
}

// Where where
func (b *builder) Where(query interface{}, value ...interface{}) Builder {
	if expressions := b.db.Statement.BuildCondition(query, value); len(expressions) > 0 {
		b.db.Statement.AddClause(clause.Where{Exprs: expressions})
	}

	return b
}

// WhereFunc where func
func (b *builder) WhereFunc(f WhereFunc) Builder {
	f(b)
	return b
}

// WhereFunc where func
func (b *builder) Limit(limit int) Builder {
	b.limit = limit
	return b
}

// WhereFunc where func
func (b *builder) Page(page int) Builder {
	b.page = page

	return b
}

// Order order
func (b *builder) Order(value interface{}) Builder {
	switch v := value.(type) {
	case clause.OrderByColumn:
		b.db.Statement.AddClause(clause.OrderBy{
			Columns: []clause.OrderByColumn{v},
		})
	default:
		b.db.Statement.AddClause(clause.OrderBy{
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

	b.db.Statement.AddClause(clause.GroupBy{
		Columns: []clause.Column{{Name: name, Raw: len(fields) != 1}},
	})

	return b
}

func (b *builder) Paginate(dest interface{}) (*Pagination, error) {
	b.preparePaginate()

	var countStatement = &gorm.Statement{
		DB:       b.db.Statement.DB,
		Clauses:  b.db.Statement.Clauses,
		SQL:      b.db.Statement.SQL,
		Schema:   b.db.Statement.Schema,
		Table:    b.db.Statement.Table,
		Vars:     b.db.Statement.Vars,
		Selects:  b.db.Statement.Selects,
		Omits:    b.db.Statement.Omits,
		Distinct: b.db.Statement.Distinct,
	}

	// Build count sql
	countStatement.Build("SELECT", "FROM", "WHERE", "GROUP BY", "FOR", "ORDER BY")
	var countSQL = countStatement.SQL.String()
	var countSQLVars = countStatement.Vars

	// Build sql
	b.db.Statement.AddClause(clause.Limit{
		Limit:  b.limit,
		Offset: b.offset,
	})
	b.db.Statement.Build("SELECT", "FROM", "WHERE", "GROUP BY", "FOR", "ORDER BY", "LIMIT")
	var sql = b.db.Statement.SQL.String()
	var vars = b.db.Statement.Vars

	var done = make(chan bool, 1)
	var count int64

	go func() {
		var countSQL = fmt.Sprintf("SELECT COUNT(*) as count FROM (%s) t", countSQL)
		fmt.Println("countSQL", countSQL)
		fmt.Println("countSQLVars", countSQLVars)
		var err = b.db.Raw(countSQL, countSQLVars...).Row().Scan(&count)
		if err != nil {
			b.logger.Fatal(err)
		}
		done <- true
	}()

	var err = b.db.Raw(sql, vars...).Scan(&dest).Error
	if err != nil {
		b.logger.Fatal(err)
	}

	<-done

	var pagination = b.finalizePaginate(count, &dest)

	return pagination, nil
}

func (b *builder) PaginateFunc(execFunc ExecFunc) (*Pagination, error) {
	b.preparePaginate()

	// b.db.Statement.AddClause(clause.Limit{
	// 	Limit:  b.limit,
	// 	Offset: b.offset,
	// })

	var session = b.db.Session(&gorm.Session{PrepareStmt: true})

	// Build count SQL
	session.Statement.Build("SELECT", "FROM", "WHERE", "GROUP BY", "FOR")

	var countSQL = session.Statement.SQL.String()
	var countSQLVars = session.Statement.Vars

	session.Statement.Build("ORDER BY", "LIMIT")
	var sql = session.Statement.SQL.String()
	var vars = session.Statement.Vars

	var done = make(chan bool, 1)
	var count int64

	go func() {
		var countSQL = fmt.Sprintf("SELECT COUNT(*) as count FROM (%s) t", countSQL)
		var err = session.Raw(countSQL, countSQLVars...).Row().Scan(&count)
		if err != nil {
			b.logger.Fatal(err)
		}
		done <- true
	}()

	o, err := execFunc(session.Raw(sql, vars...))
	if err != nil {
		b.logger.Fatal(err)
	}

	<-done

	var pagination = b.finalizePaginate(count, &o)

	return pagination, nil
}

func (b *builder) preparePaginate() {
	if b.skipLimitOffset {
		return
	}

	if b.page < 1 {
		b.page = 1
	}
	if b.limit == 0 {
		b.limit = 10
	}

	if b.page == 1 {
		b.offset = 0
	} else {
		b.offset = (b.page - 1) * b.limit
	}
}

func (b *builder) finalizePaginate(count int64, records interface{}) *Pagination {
	var pagination = &Pagination{
		TotalRecord: int(count),
		Records:     records,
		Page:        b.page,
		PerPage:     b.limit,
		Offset:      b.offset,
	}

	if b.limit > 0 {
		pagination.PerPage = b.limit
		pagination.TotalPage = int(math.Ceil(float64(count) / float64(b.limit)))
	} else {
		pagination.TotalPage = 1
		pagination.PerPage = int(count)
	}

	if b.skipLimitOffset {
		pagination.PerPage = int(count)
		pagination.TotalRecord = int(count)
		pagination.Offset = 0

	}

	if b.page > 1 {
		pagination.PrevPage = b.page - 1
	} else {
		pagination.PrevPage = b.page
	}

	if b.page == pagination.TotalPage {
		pagination.NextPage = b.page
	} else {
		pagination.NextPage = b.page + 1
	}

	pagination.HasNext = pagination.TotalPage > pagination.Page
	pagination.HasPrev = pagination.Page > 1

	if !pagination.HasNext {
		pagination.NextPage = pagination.Page
	}

	return pagination
}
