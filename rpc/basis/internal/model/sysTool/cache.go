package sysTool

import (
	"time"

	"td27/rpc/basis/internal/model/common"
)

// CacheModel System cache entity
type CacheModel struct {
	common.Td27Model
	Username  string    `json:"user" db:"username"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	ExpiresAt time.Time `json:"expiresAt" db:"expires_at"`
}

func (CacheModel) TableName() string {
	return "sys_tool_cache"
}
