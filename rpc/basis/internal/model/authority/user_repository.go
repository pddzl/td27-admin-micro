package authority

import (
	"context"
	"errors"
	"fmt"

	"basis/internal/model"
	"basis/internal/pkg"

	"gorm.io/gorm"
)

type (
	authorityUserModel interface {
		List(ctx context.Context, page int, pageSize int) ([]AuthorityUserPlusRoleNameDTO, int64, error)
		Insert(ctx context.Context, data *AuthorityUserEntity) error
		FindOne(ctx context.Context, id uint) (*AuthorityUserPlusRoleNameDTO, error)
		Update(ctx context.Context, data *AuthorityUserEntity) (*AuthorityUserEntity, error)
		Delete(ctx context.Context, id uint) error
		ModifyPassword(ctx context.Context, id uint, oldPassword string, newPassword string) error
		SwitchUserActive(ctx context.Context, id uint, active bool) error
		FindByRoleId(ctx context.Context, roleId uint) (bool, error)
	}

	defaultAuthorityUserModel struct {
		conn  *gorm.DB
		table string
	}

	AuthorityUserEntity struct {
		model.Td27Model
		Username    string `gorm:"unique;comment:用户名"` // 用户名
		Password    string `gorm:"not null;comment:密码"`
		Phone       string `gorm:"comment:手机号"` // 手机号
		Email       string `gorm:"comment:邮箱"`  // 邮箱
		Active      bool   // 是否活跃
		RoleModelID uint   `gorm:"not null"`
	}

	AuthorityUserPlusRoleNameDTO struct {
		AuthorityUserEntity
		RoleName string
	}
)

func (AuthorityUserEntity) TableName() string {
	return "authority_user"
}

func newAuthorityUserModel(conn *gorm.DB) *defaultAuthorityUserModel {
	return &defaultAuthorityUserModel{
		conn:  conn,
		table: "`authority_user`",
	}
}

func (m *defaultAuthorityUserModel) List(ctx context.Context, page int, pageSize int) ([]AuthorityUserPlusRoleNameDTO, int64, error) {
	var authorityUserPlusRoleNameList []AuthorityUserPlusRoleNameDTO
	var total int64

	db := m.conn.WithContext(ctx).Model(&AuthorityUserEntity{})

	// 分页
	err := db.Count(&total).Error
	if err != nil {
		return authorityUserPlusRoleNameList, total, fmt.Errorf("count err: %v", err)
	} else {
		offset := pageSize * (page - 1)
		db = db.Limit(pageSize).Offset(offset)
		// 左连接 查询出role_name
		db.Select("authority_user.id,authority_user.username,authority_user.phone,authority_user.email,authority_user.active,authority_user.role_model_id,authority_role.role_name").Joins("left join authority_role on authority_user.role_model_id = authority_role.id").Scan(&authorityUserPlusRoleNameList)
	}

	return authorityUserPlusRoleNameList, total, err
}

func (m *defaultAuthorityUserModel) Delete(ctx context.Context, id uint) error {
	return m.conn.WithContext(ctx).Where("id = ?", id).Unscoped().Delete(&AuthorityUserEntity{}).Error
}

func (m *defaultAuthorityUserModel) FindOne(ctx context.Context, id uint) (*AuthorityUserPlusRoleNameDTO, error) {
	var authorityUserPlusRoleName AuthorityUserPlusRoleNameDTO
	err := m.conn.WithContext(ctx).
		Table("authority_user").
		Select("authority_user.created_at,authority_user.id,authority_user.username,authority_user.phone,authority_user.email,authority_user.active,authority_user.role_model_id,authority_role.role_name").
		Joins("inner join authority_role on authority_user.role_model_id = authority_role.id").Where("authority_user.id = ?", id).
		Scan(&authorityUserPlusRoleName).Error
	return &authorityUserPlusRoleName, err
}

func (m *defaultAuthorityUserModel) Insert(ctx context.Context, data *AuthorityUserEntity) error {
	data.Password = pkg.MD5V([]byte(data.Password))
	return m.conn.WithContext(ctx).Create(data).Error
}

func (m *defaultAuthorityUserModel) Update(ctx context.Context, newData *AuthorityUserEntity) (*AuthorityUserEntity, error) {
	// Check existence
	var existing AuthorityUserEntity
	if err := m.conn.WithContext(ctx).First(&existing, newData.ID).Error; err != nil {
		return nil, err
	}

	// todo
	// check role existence

	// Update (use map to avoid zero-value problems)
	err := m.conn.Model(&existing).Updates(map[string]interface{}{
		"username":      newData.Username,
		"phone":         newData.Phone,
		"email":         newData.Email,
		"active":        newData.Active,
		"role_model_id": newData.RoleModelID,
	}).Error
	if err != nil {
		return nil, err
	}

	// Reload updated record (DB is source of truth)
	if err = m.conn.First(&existing, newData.ID).Error; err != nil {
		return nil, err
	}

	return &existing, nil
}

func (m *defaultAuthorityUserModel) ModifyPassword(ctx context.Context, id uint, oldPassword string, newPassword string) error {
	conn := m.conn.WithContext(ctx)
	var authorityUser AuthorityUserEntity
	if errors.Is(conn.Where("id = ? and password = ?", id, pkg.MD5V([]byte(oldPassword))).First(&authorityUser).Error, gorm.ErrRecordNotFound) {
		return errors.New("wrong old password")
	}

	return conn.WithContext(ctx).Model(&authorityUser).Update("password", pkg.MD5V([]byte(newPassword))).Error
}

func (m *defaultAuthorityUserModel) SwitchUserActive(ctx context.Context, id uint, active bool) error {
	conn := m.conn.WithContext(ctx)
	var authorityUser AuthorityUserEntity
	if errors.Is(conn.Where("id = ?", id).First(&authorityUser).Error, gorm.ErrRecordNotFound) {
		return errors.New("record not found")
	}

	return conn.Model(&authorityUser).Update("active", active).Error
}

func (m *defaultAuthorityUserModel) FindByRoleId(ctx context.Context, roleId uint) (bool, error) {
	var count int64

	err := m.conn.WithContext(ctx).
		Model(&AuthorityUserEntity{}).
		Where("role_model_id = ?", roleId).
		Limit(1).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *defaultAuthorityUserModel) tableName() string {
	return m.table
}
