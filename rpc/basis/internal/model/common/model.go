package common

import (
	"time"
)

type Td27Model struct {
	ID        uint       `json:"id" db:"id"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

type PageInfo struct {
	Page     int `json:"page" form:"page" default:"1"`
	PageSize int `json:"pageSize" form:"pageSize" default:"10"`
}

// NewPageInfo builds a PageInfo with invalid values clamped to the shared
// defaults: Page <= 0 becomes 1 and PageSize <= 0 becomes 10.
func NewPageInfo(page, pageSize int) *PageInfo {
	p := &PageInfo{Page: page, PageSize: pageSize}
	p.normalize()
	return p
}

func (p *PageInfo) normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
}
