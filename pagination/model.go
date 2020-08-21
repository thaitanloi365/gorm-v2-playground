package pagination

import (
	"gorm.io/gorm"
)

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

// Builder builder
type Builder interface {
	WithRawSQL(rawSQL *gorm.DB) Builder
	Limit(int) Builder
	Page(int) Builder
	Order(interface{}) Builder
	Paginate(result interface{}) *Pagination
	SkipLimitOffset(bool) Builder
}

type builder struct {
	db              *gorm.DB
	rawSQL          *gorm.DB
	page            int
	limit           int
	offset          int
	skipLimitOffset bool
	orderBy         []interface{}
}
