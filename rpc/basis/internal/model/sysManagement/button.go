package sysManagement

import "td27/rpc/basis/internal/model/common"

// ButtonModel Button permission entity
type ButtonModel struct {
	common.Td27Model
	ButtonCode  string `json:"button_code" db:"button_code"`
	ButtonName  string `json:"button_name" db:"button_name"`
	Description string `json:"description" db:"description"`
	PagePath    string `json:"page_path" db:"page_path"`
}

func (ButtonModel) TableName() string {
	return "sys_management_button"
}

// DTOs for button operations
type ButtonDto struct {
	ID            uint   `json:"id"`
	ButtonCode    string `json:"buttonCode"`
	ButtonName    string `json:"buttonName"`
	Description   string `json:"description"`
	PagePath      string `json:"pagePath"`
	HasPermission bool   `json:"hasPermission"`
}

type CreateButtonReq struct {
	ButtonCode  string `json:"buttonCode" binding:"required,max=100"`
	ButtonName  string `json:"buttonName" binding:"required,max=100"`
	Description string `json:"description"`
	PagePath    string `json:"pagePath" binding:"required,max=200"`
}

type UpdateButtonReq struct {
	ID          uint   `json:"id" binding:"required"`
	ButtonCode  string `json:"buttonCode" binding:"required,max=100"`
	ButtonName  string `json:"buttonName" binding:"required,max=100"`
	Description string `json:"description"`
	PagePath    string `json:"pagePath" binding:"required,max=200"`
}

type ListButtonReq struct {
	common.PageInfo
	PagePath string `json:"pagePath" form:"pagePath"`
}

type CheckButtonReq struct {
	ButtonCode string `json:"buttonCode" binding:"required"`
}
