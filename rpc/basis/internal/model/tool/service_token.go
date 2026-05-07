package tool

import (
	"td27/rpc/basis/internal/model/common"
)

// ServiceToken Service-to-service authentication token entity
type ServiceToken struct {
	common.Td27Model
	Name      string `json:"name" gorm:"size:100;unique;not null;comment:令牌名称/描述"`
	TokenHash string `json:"-" gorm:"size:255;unique;not null;comment:令牌哈希值"`
	Status    bool   `json:"status" gorm:"default:true;comment:是否启用"`
	ExpiresAt *int64 `json:"expiresAt" gorm:"comment:过期时间戳(秒)"`
}

func (ServiceToken) TableName() string {
	return "sys_tool_service_token"
}

// TokenPermission Service token permission join table
type TokenPermission struct {
	TokenID      uint `json:"tokenId" gorm:"primaryKey"`
	PermissionID uint `json:"permissionId" gorm:"primaryKey"`
}

func (TokenPermission) TableName() string {
	return "sys_tool_token_permission"
}
