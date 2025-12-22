package model

import (
	"time"

	"gorm.io/gorm"
)

type Td27Model struct {
	ID        uint           `gorm:"primarykey"`                                   // 主键ID
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;default:null"` // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;default:null"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index"`                                        // 删除时间
}
