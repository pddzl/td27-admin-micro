package sysManagement

import (
	"td27/rpc/basis/internal/model/common"
)

// DictDetailModel Dictionary detail entity
type DictDetailModel struct {
	common.Td27Model
	Label       string             `json:"label" db:"label"`
	Value       string             `json:"value" db:"value"`
	Sort        int                `json:"sort" db:"sort"`
	DictModelID int                `json:"dictId" db:"dict_id" gorm:"column:dict_id" binding:"required"`
	ParentID    *int               `json:"parentId" db:"parent_id" gorm:"column:parent_id"`
	Children    []*DictDetailModel `json:"children" gorm:"-"`
	Description string             `json:"description" db:"description"`
}

func (ddm *DictDetailModel) TableName() string {
	return "sys_management_dict_detail"
}
