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
	}

	defaultAuthorityRoleModel struct {
		conn     *gorm.DB
		table    string
		userRepo AuthorityUserModel
		menuRepo AuthorityMenuModel
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
		conn:     conn,
		table:    "`authority_role`",
		userRepo: NewAuthorityUserModel(conn),
		menuRepo: NewAuthorityMenuModel(conn),
	}
}

func (m *defaultAuthorityRoleModel) Delete(ctx context.Context, id uint) error {
	var authorityRoleEntity AuthorityRoleEntity
	conn := m.conn.WithContext(ctx)

	if errors.Is(conn.Where("id = ?", id).First(&authorityRoleEntity).Error, gorm.ErrRecordNotFound) {
		return errors.New("record not found")
	}

	userExist, err := m.userRepo.FindByRoleId(ctx, id)
	if err != nil {
		return err
	}
	if userExist {
		return errors.New("user exist under role")
	}

	err = conn.Delete(&authorityRoleEntity).Error
	if err != nil {
		return fmt.Errorf("delete role err: %v", err)
	}

	// empty associated menus
	err = conn.Model(&authorityRoleEntity).Association("Menus").Clear()
	if err != nil {
		return fmt.Errorf("empty associated menus err: %v", err)
	}

	// todo
	// casbin rule

	return nil
}

func (m *defaultAuthorityRoleModel) FindOne(ctx context.Context, id uint) (*AuthorityRoleEntity, error) {
	return nil, nil
}

func (m *defaultAuthorityRoleModel) Insert(ctx context.Context, data *AuthorityRoleEntity) error {
	return m.conn.WithContext(ctx).Create(data).Error
}

func (m *defaultAuthorityRoleModel) Update(ctx context.Context, id uint, roleName string) error {
	var authorityRoleEntity AuthorityRoleEntity
	conn := m.conn.WithContext(ctx)

	if errors.Is(conn.Where("id = ?", id).First(&authorityRoleEntity).Error, gorm.ErrRecordNotFound) {
		return errors.New("record not found")
	}

	return conn.Model(&authorityRoleEntity).Update("role_name", roleName).Error
}

func (m *defaultAuthorityRoleModel) UpdateRoleMenus(ctx context.Context, roleId uint, ids []uint) (err error) {
	var authorityRoleEntity AuthorityRoleEntity
	conn := m.conn.WithContext(ctx)

	if errors.Is(conn.Where("id = ?", roleId).First(&authorityRoleEntity).Error, gorm.ErrRecordNotFound) {
		return errors.New("record not found")
	}

	authorityMenuEntity, err := m.menuRepo.FindByIds(ctx, ids)
	if err != nil {
		return fmt.Errorf("find menus err: %v", err)
	}

	err = conn.Model(&authorityRoleEntity).Association("Menus").Replace(authorityMenuEntity)
	if err != nil {
		return fmt.Errorf("replace menu err: %v", err)
	}

	return err
}

func (m *defaultAuthorityRoleModel) tableName() string {
	return m.table
}
