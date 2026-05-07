package authority

import (
	"gorm.io/gorm"
)

// DataScope Data permission scope
type DataScope string

const (
	DataScopeAll    DataScope = "all"    // All data
	DataScopeDept   DataScope = "dept"   // Current department data
	DataScopeSelf   DataScope = "self"   // Only own data
	DataScopeCustom DataScope = "custom" // Custom SQL filter
)

// DataPermission Data permission configuration
type DataPermission struct {
	Scope     DataScope // Data access scope
	DeptID    uint      // Department ID (used when Scope is dept)
	UserID    uint      // User ID (used when Scope is self)
	CustomSQL string    // Custom SQL filter (used when Scope is custom)
}

// ApplyDataScope Applies data permission filter to GORM query
func ApplyDataScope(db *gorm.DB, perm *DataPermission, tableAlias, deptColumn, userColumn string) *gorm.DB {
	if perm == nil {
		return db
	}

	if tableAlias != "" {
		tableAlias = tableAlias + "."
	}

	switch perm.Scope {
	case DataScopeAll:
		// No filter for all data
		return db
	case DataScopeDept:
		// Filter by department
		if deptColumn != "" {
			return db.Where(tableAlias+deptColumn+" = ?", perm.DeptID)
		}
		return db
	case DataScopeSelf:
		// Filter by current user
		if userColumn != "" {
			return db.Where(tableAlias+userColumn+" = ?", perm.UserID)
		}
		return db
	case DataScopeCustom:
		// Apply custom SQL filter
		if perm.CustomSQL != "" {
			return db.Where(perm.CustomSQL)
		}
		return db
	default:
		// Default to self only
		if userColumn != "" {
			return db.Where(tableAlias+userColumn+" = ?", perm.UserID)
		}
		return db
	}
}
