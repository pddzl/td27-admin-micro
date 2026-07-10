package sysTool

import (
	"time"

	"td27/rpc/basis/internal/model/common"
)

// CacheModel System cache entity
type CacheModel struct {
	common.Td27Model
	Username  string    `json:"user" gorm:"column:username;comment:用户名"`
	Key       string    `json:"key" gorm:"uniqueIndex;size:255;comment:缓存键"`
	Value     string    `json:"value" gorm:"type:text;comment:缓存值"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"index;comment:过期时间"`
}

func (CacheModel) TableName() string {
	return "sys_tool_cache"
}
