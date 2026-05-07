package authority

import (
	"td27/rpc/basis/internal/model/common"
)

// DictModel Dictionary entity
type DictModel struct {
	common.Td27Model
	CNName      string            `json:"cn_name" gorm:"column:cn_name;unique" binding:"required"`
	ENName      string            `json:"en_name" gorm:"column:en_name;unique" binding:"required"`
	DictDetails []DictDetailModel `json:"dictDetails"`
}

func (dm *DictModel) TableName() string {
	return "sys_management_dict"
}
