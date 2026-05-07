package authority

import "td27/rpc/basis/internal/model/common"

// DeptModel Department entity for data permission
type DeptModel struct {
	common.Td27Model
	DeptName string `json:"deptName" gorm:"unique;size:100;not null;comment:部门名称"`
	ParentID uint   `json:"parentId" gorm:"index;comment:父部门ID"`
	Path     string `json:"path" gorm:"size:500;index;comment:部门路径(materialized path),如:/1/2/3/"`
	Level    uint   `json:"level" gorm:"NOT NULL;comment:depth level"`
	Sort     uint   `json:"sort" gorm:"default:0"`
	Status   bool   `json:"status" gorm:"default:true"`
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
