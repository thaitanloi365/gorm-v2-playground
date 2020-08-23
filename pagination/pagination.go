package pagination

import (
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
)

// New new builder
func New(db *gorm.DB) Builder {
	return &builder{
		db: db,

		skipLimitOffset: false,
		orderBy:         []interface{}{},
	}
}

func (b *builder) WithRawSQL(rawSQL *gorm.DB) Builder {
	b.rawSQL = rawSQL
	return b
}

func (b *builder) Limit(limit int) Builder {
	b.limit = limit
	return b
}

func (b *builder) Page(page int) Builder {
	b.page = page
	return b
}

func (b *builder) Order(orderBy interface{}) Builder {
	b.orderBy = append(b.orderBy, orderBy)
	return b
}

func (b *builder) SkipLimitOffset(skip bool) Builder {
	b.skipLimitOffset = skip
	return b
}

func (b *builder) Paginate(result interface{}) *Pagination {

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

	var done = make(chan bool, 1)
	var count int64

	var sqlString = b.rawSQL.Statement.SQL.String()

	go func(sqlString string) {
		var sessionConfig = &gorm.Session{PrepareStmt: true}
		if b.rawSQL == nil {
			b.db.Session(sessionConfig).Model(result).Count(&count)
		} else {
			var sqlString = b.rawSQL.Statement.SQL.String()
			var vars = b.rawSQL.Statement.Vars
			var rawSQL = fmt.Sprintf("SELECT COUNT(*) as count FROM (%s) t", sqlString)
			b.rawSQL.Session(sessionConfig).Raw(rawSQL, vars).Row().Scan(&count)
		}
		done <- true
	}(sqlString)

	if b.rawSQL == nil {
		var session = b.db.Session(&gorm.Session{PrepareStmt: true})
		if b.skipLimitOffset == false {
			if len(b.orderBy) > 0 {
				for _, o := range b.orderBy {
					session = session.Order(o)
				}
			}
			if b.limit > 0 {
				session = session.Limit(b.limit)
			}
			if b.offset >= 0 {
				session = session.Offset(b.offset)
			}
		}
		session.Find(result)
	} else {
		var session = b.rawSQL.Session(&gorm.Session{PrepareStmt: true})
		var sqlString = b.rawSQL.Statement.SQL.String()
		var vars = b.rawSQL.Statement.Vars
		var rawSQL = sqlString
		if b.skipLimitOffset == false {
			if len(b.orderBy) > 0 {
				var orderByString = []string{}
				for _, orderBy := range b.orderBy {
					if v, ok := orderBy.(string); ok {
						orderByString = append(orderByString, v)
					}
				}

				rawSQL = fmt.Sprintf("%s ORDER %s", rawSQL, strings.Join(orderByString, ","))
			}
			if b.limit > 0 {
				rawSQL = fmt.Sprintf("%s LIMIT %d", rawSQL, b.limit)
			}
			if b.offset >= 0 {
				rawSQL = fmt.Sprintf("%s OFFSET %d", rawSQL, b.offset)
			}
		}

		session.Raw(rawSQL, vars).Scan(result)
	}

	<-done

	var pagination = &Pagination{
		TotalRecord: int(count),
		Records:     result,
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
