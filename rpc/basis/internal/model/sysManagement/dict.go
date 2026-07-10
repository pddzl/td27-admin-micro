package sysManagement

import (
	"td27/rpc/basis/internal/model/common"
)

// DictModel Dictionary entity
type DictModel struct {
	common.Td27Model
	CNName      string            `json:"cn_name" db:"cn_name"`
	ENName      string            `json:"en_name" db:"en_name"`
	DictDetails []DictDetailModel `json:"dictDetails"`
}

func (dm *DictModel) TableName() string {
	return "sys_management_dict"
}
