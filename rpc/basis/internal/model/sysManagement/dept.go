package sysManagement

import "td27/rpc/basis/internal/model/common"

// DeptModel Department entity for data permission
type DeptModel struct {
	common.Td27Model
	DeptName string `json:"deptName" db:"dept_name"`
	ParentID uint   `json:"parentId" db:"parent_id"`
	Path     string `json:"path" db:"path"`
	Level    uint   `json:"level" db:"level"`
	Sort     uint   `json:"sort" db:"sort"`
	Status   bool   `json:"status" db:"status"`
}

// GetFullPath returns full path including current department
func (d *DeptModel) GetFullPath() string {
	if d.Path == "" {
		return "/" + string(rune(d.ID))
	}
	return d.Path + string(rune(d.ID)) + "/"
}

// IsAncestorOf checks if current department is ancestor of target department
func (d *DeptModel) IsAncestorOf(targetPath string) bool {
	fullPath := d.GetFullPath()
	return len(targetPath) > len(fullPath) && targetPath[:len(fullPath)] == fullPath
}

// IsDescendantOf checks if current department is descendant of ancestor department
func (d *DeptModel) IsDescendantOf(ancestorPath string) bool {
	return len(d.Path) >= len(ancestorPath) && d.Path[:len(ancestorPath)] == ancestorPath
}

func (DeptModel) TableName() string {
	return "sys_management_dept"
}
