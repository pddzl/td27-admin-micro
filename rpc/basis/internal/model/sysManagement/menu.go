package sysManagement

import (
	"td27/rpc/basis/internal/model/common"
)

// MenuModel Menu entity
type MenuModel struct {
	common.Td27Model
	MenuName   string `json:"menu_name" db:"menu_name"`
	Icon       string `json:"icon" db:"icon" gorm:"size:100;comment:图标"`
	Path       string `json:"path" db:"path"`
	Component  string `json:"component" db:"component" gorm:"size:200;comment:前端组件"`
	Redirect   string `json:"redirect" db:"redirect" gorm:"size:200;comment:重定向"`
	ParentID   uint   `json:"parentId" db:"parent_id" gorm:"index;comment:父菜单ID"`
	Sort       uint   `json:"sort" db:"sort" gorm:"default:0;comment:排序"`
	Hidden     bool   `json:"hidden" db:"hidden" gorm:"default:false;comment:是否隐藏"`
	KeepAlive  bool   `json:"keepAlive" db:"keep_alive" gorm:"default:false;comment:缓存"`
	Affix      bool   `json:"affix" db:"affix" gorm:"default:false;comment:是否固定"`
	AlwaysShow bool   `json:"alwaysShow" db:"always_show" gorm:"default:false;comment:一直显示根路由"`
	Title      string `json:"title" db:"title" gorm:"unique;comment:菜单名"`
}

func (MenuModel) TableName() string {
	return "sys_management_menu"
}
