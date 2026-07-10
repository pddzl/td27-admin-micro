package sysTool

import (
	"td27/rpc/basis/internal/model/common"
)

// FileModel File upload entity
type FileModel struct {
	common.Td27Model
	FileName string `json:"fileName" db:"file_name"`
	FullPath string `json:"fullPath" db:"full_path"`
	Mime     string `json:"mime" db:"mime"`
}

func (FileModel) TableName() string {
	return "sys_tool_file"
}
