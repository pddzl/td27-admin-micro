package sysManagement

import (
	"td27/rpc/basis/internal/model/common"
)

// ApiModel API entity
type ApiModel struct {
	common.Td27Model
	Path        string `json:"path" db:"path"`
	Method      string `json:"method" db:"method"`
	GroupEN     string `json:"group_en" db:"group_en"`
	GroupCN     string `json:"group_cn" db:"group_cn"`
	Description string `json:"description" db:"description"`
}

func (ApiModel) TableName() string {
	return "sys_management_api"
}
