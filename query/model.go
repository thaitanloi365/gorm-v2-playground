package query

import (
	"log"

	"gorm.io/gorm"
)

// WhereFunc where func
type WhereFunc = func(b Builder)

// ExecFunc exec func
type ExecFunc = func(db *gorm.DB) (records interface{}, err error)

type builder struct {
	db              *gorm.DB
	statement       *gorm.Statement
	page            int
	limit           int
	offset          int
	skipLimitOffset bool
	logger          *log.Logger
}

// Pagination ...
type Pagination struct {
	HasNext     bool        `json:"has_next"`
	HasPrev     bool        `json:"has_prev"`
	PerPage     int         `json:"per_page"`
	NextPage    int         `json:"next_page"`
	Page        int         `json:"current_page"`
	PrevPage    int         `json:"prev_page"`
	Offset      int         `json:"offset"`
	Records     interface{} `json:"records"`
	TotalRecord int         `json:"total_record"`
	TotalPage   int         `json:"total_page"`
	Metadata    interface{} `json:"metadata"`
}

// Builder interface
type Builder interface {
	Where(query interface{}, value ...interface{}) Builder
	WhereFunc(WhereFunc) Builder
	Order(order interface{}) Builder
	Scan(dest interface{}) error
	Prepare() *gorm.Statement
	Group(name string) Builder
	Limit(int) Builder
	Page(int) Builder
	Paginate(dest interface{}) (*Pagination, error)
}
