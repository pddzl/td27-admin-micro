package authority

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"basis/internal/model"
)

type (
	authorityRoleModel interface {
		Insert(ctx context.Context, data *AuthorityRoleEntity) error
		FindOne(ctx context.Context, id uint) (*AuthorityRoleEntity, error)
		Update(ctx context.Context, id uint, roleName string) error
		Delete(ctx context.Context, id uint) error
		UpdateRoleMenus(ctx context.Context, roleId uint, ids []uint) error
		ExistsById(ctx context.Context, id uint64) (bool, error)
	}

	defaultAuthorityRoleModel struct {
		conn  *gorm.DB
		table string
	}

	AuthorityRoleEntity struct {
		model.Td27Model
		RoleName string                 `json:"roleName" gorm:"unique" binding:"required"`
		Menus    []*AuthorityMenuEntity `json:"menus" gorm:"many2many:role_menus;"`
	}
)

func (AuthorityRoleEntity) TableName() string {
	return "authority_role"
}

func newAuthorityRoleModel(conn *gorm.DB) *defaultAuthorityRoleModel {
	return &defaultAuthorityRoleModel{
		conn:  conn,
		table: "`authority_role`",
	}
}

func (ar *defaultAuthorityRoleModel) Delete(ctx context.Context, id uint) error {
	var authorityRoleEntity AuthorityRoleEntity
	conn := ar.conn.WithContext(ctx)

	err := conn.Where("id = ?", id).Unscoped().Delete(&authorityRoleEntity).Error
	if err != nil {
		return fmt.Errorf("delete role err: %v", err)
	}

	// empty associated menus
	err = conn.Model(&authorityRoleEntity).Association("Menus").Clear()
	if err != nil {
		return fmt.Errorf("empty associated menus err: %v", err)
	}

	return nil
}

func (ar *defaultAuthorityRoleModel) FindOne(ctx context.Context, id uint) (*AuthorityRoleEntity, error) {
	return nil, nil
}

func (ar *defaultAuthorityRoleModel) Insert(ctx context.Context, data *AuthorityRoleEntity) error {
	return ar.conn.WithContext(ctx).Create(data).Error
}

func (ar *defaultAuthorityRoleModel) Update(ctx context.Context, id uint, roleName string) error {
	var authorityRoleEntity AuthorityRoleEntity
	conn := ar.conn.WithContext(ctx)

	if errors.Is(conn.Where("id = ?", id).First(&authorityRoleEntity).Error, gorm.ErrRecordNotFound) {
		return errors.New("record not found")
	}

	return conn.Model(&authorityRoleEntity).Update("role_name", roleName).Error
}

func (ar *defaultAuthorityRoleModel) UpdateRoleMenus(ctx context.Context, roleId uint, ids []uint) (err error) {
	return err
}

func (ar *defaultAuthorityRoleModel) ExistsById(ctx context.Context, id uint64) (bool, error) {
	var count int64

	err := ar.conn.WithContext(ctx).
		Model(&AuthorityRoleEntity{}).
		Where("id = ?", id).
		Limit(1).
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (ar *defaultAuthorityRoleModel) tableName() string {
	return ar.table
}
