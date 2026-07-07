package common

import (
	"time"
)

type Td27Model struct {
	ID        uint       `json:"id" db:"id"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

type PageInfo struct {
	Page     int `json:"page" form:"page" default:"1"`
	PageSize int `json:"pageSize" form:"pageSize" default:"10"`
}
