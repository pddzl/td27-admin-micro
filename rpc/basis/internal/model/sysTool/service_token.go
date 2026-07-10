package sysTool

import (
	"td27/rpc/basis/internal/model/common"
)

// ServiceToken Service-to-service authentication token entity
type ServiceToken struct {
	common.Td27Model
	Name      string `json:"name" db:"name"`
	TokenHash string `json:"-" db:"token_hash"`
	Status    bool   `json:"status" db:"status"`
	ExpiresAt *int64 `json:"expiresAt" db:"expires_at"`
}

func (ServiceToken) TableName() string {
	return "sys_tool_service_token"
}

// TokenPermission Service token permission join table
type TokenPermission struct {
	TokenID      uint `json:"tokenId" db:"token_id"`
	PermissionID uint `json:"permissionId" db:"permission_id"`
}

func (TokenPermission) TableName() string {
	return "sys_tool_token_permission"
}
