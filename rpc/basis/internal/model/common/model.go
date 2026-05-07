package common

import (
	"time"

	"gorm.io/gorm"
)

// Td27Model Base model for all entities
type Td27Model struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// PageInfo Pagination request parameters
type PageInfo struct {
	Page     int `json:"page" form:"page" default:"1"`
	PageSize int `json:"pageSize" form:"pageSize" default:"10"`
}

