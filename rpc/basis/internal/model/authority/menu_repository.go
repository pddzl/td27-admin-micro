package authority

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"gorm.io/gorm"

	"basis/internal/model"
	"database/sql/driver"
)

type (
	authorityMenuModel interface {
		Insert(ctx context.Context, data *MenuDTO) error
		Update(ctx context.Context, newData *NewMenuDTO) (*AuthorityMenuEntity, error)
		Delete(ctx context.Context, id uint) error
		GetElTreeMenus(roleId uint) ([]AuthorityMenuEntity, []uint, error)
		FindByIds(ctx context.Context, ids []uint) ([]AuthorityMenuEntity, error)
	}

	defaultAuthorityMenuModel struct {
		conn  *gorm.DB
		table string
	}

	AuthorityMenuEntity struct {
		model.Td27Model
		Pid       uint                   `json:"pid"`                       // 父菜单ID
		Name      string                 `json:"name"`                      // 路由名称
		Path      string                 `json:"path" gorm:"unique"`        // 路由路径
		Redirect  string                 `json:"redirect"`                  // 重定向
		Component string                 `json:"component" gorm:"not null"` // 前端组件
		Sort      uint                   `json:"sort" gorm:"not null"`      // 排序
		Meta      Meta                   `json:"meta" gorm:"type:json"`     // 元数据
		Children  []AuthorityMenuEntity  `json:"children" gorm:"-"`
		Roles     []*AuthorityRoleEntity `json:"roles" gorm:"many2many:role_menus;"`
	}

	MenuDTO struct {
		Pid       uint   `json:"pid"`                          // 默认0 根目录
		Name      string `json:"name"`                         // 名称
		Path      string `json:"path" binding:"required"`      // 路径
		Redirect  string `json:"redirect"`                     // 重定向
		Component string `json:"component" binding:"required"` // 前端组件
		Sort      uint   `json:"sort" binding:"required"`      // 排序
		Meta      Meta   `json:"meta"`
	}

	NewMenuDTO struct {
		ID uint
		MenuDTO
	}

	Meta struct {
		Hidden     bool   `json:"hidden,omitempty"`  // 菜单是否隐藏
		Title      string `json:"title,omitempty"`   // 菜单名
		SvgIcon    string `json:"svgIcon,omitempty"` // svg图标
		ElIcon     string `json:"elIcon,omitempty"`  // element-plus图标
		Affix      bool   `json:"affix,omitempty"`   // 是否固定
		KeepAlive  bool   `json:"keepAlive,omitempty"`
		AlwaysShow bool   `json:"alwaysShow,omitempty"` // 是否一直显示根路由
	}
)

func (m Meta) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	return string(b), err
}

func (m *Meta) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), m)
}

func (AuthorityMenuEntity) TableName() string {
	return "authority_menu"
}

func newAuthorityMenuModel(conn *gorm.DB) *defaultAuthorityMenuModel {
	return &defaultAuthorityMenuModel{
		conn:  conn,
		table: "`authority_menu`",
	}
}

func (m *defaultAuthorityMenuModel) Delete(ctx context.Context, id uint) error {
	var entity AuthorityMenuEntity
	conn := m.conn.WithContext(ctx)

	if errors.Is(conn.Where("id = ?", id).First(&entity).Error, gorm.ErrRecordNotFound) {
		return errors.New("authority_menu record not found")
	}

	return conn.Unscoped().Select("Roles").Delete(&entity).Error
}

func (m *defaultAuthorityMenuModel) Insert(ctx context.Context, data *MenuDTO) error {
	return m.conn.WithContext(ctx).Create(data).Error
}

func (m *defaultAuthorityMenuModel) Update(ctx context.Context, newData *NewMenuDTO) (*AuthorityMenuEntity, error) {
	conn := m.conn.WithContext(ctx)

	// Check existence
	var authorityMenuEntity AuthorityMenuEntity
	if err := conn.First(&authorityMenuEntity, newData.ID).Error; err != nil {
		return nil, err
	}

	// update
	err := conn.Model(&authorityMenuEntity).Updates(map[string]interface{}{
		"pid":       newData.Pid,
		"name":      newData.Name,
		"path":      newData.Path,
		"component": newData.Component,
		"redirect":  newData.Redirect,
		"sort":      newData.Sort,
		"meta":      newData.Meta,
	}).Error

	// Reload updated record
	if err = conn.First(&authorityMenuEntity, newData.ID).Error; err != nil {
		return nil, err
	}

	return &authorityMenuEntity, nil
}

func getTreeMap(menuListFormat []AuthorityMenuEntity, menuList []AuthorityMenuEntity) {
	for index, menuF := range menuListFormat {
		for _, menu := range menuList {
			if menuF.ID == menu.Pid {
				// menuF 只是个复制值
				//menuF.Children = append(menuF.Children, menu)
				menuListFormat[index].Children = append(menuListFormat[index].Children, menu)
			}
		}
		if len(menuListFormat[index].Children) > 0 {
			// 排序
			sort.Slice(menuListFormat[index].Children, func(i, j int) bool {
				return menuListFormat[index].Children[i].Sort < menuListFormat[index].Children[j].Sort
			})
			getTreeMap(menuListFormat[index].Children, menuList)
		}
	}
}

// GetElTreeMenus 获取所有menu
func (m *defaultAuthorityMenuModel) GetElTreeMenus(roleId uint) ([]AuthorityMenuEntity, []uint, error) {
	var authorityMenuEntityList []AuthorityMenuEntity
	conn := m.conn.WithContext(context.Background())

	err := conn.Find(&authorityMenuEntityList).Error
	if err != nil {
		return nil, nil, fmt.Errorf("GetElTreeMenus menus error: %v", err)
	}

	menuListFormat := make([]AuthorityMenuEntity, 0)
	for _, menu := range authorityMenuEntityList {
		if menu.Pid == 0 {
			menuListFormat = append(menuListFormat, menu)
		}
	}

	getTreeMap(menuListFormat, authorityMenuEntityList)

	var authorityRoleEntity AuthorityRoleEntity
	err = conn.Where("id = ?", roleId).Preload("Menus").First(&authorityRoleEntity).Error
	if err != nil {
		return nil, nil, fmt.Errorf("GetElTreeMenus preload role error: %v", err)
	}

	// 前端el-tree 选中数据
	// 去掉夫菜单，防止直接选中父级造成全选
	roleIds := make([]uint, 0)
	count := 0
	for _, menu := range authorityRoleEntity.Menus {
		for _, menu1 := range authorityRoleEntity.Menus {
			if menu.ID == menu1.Pid {
				count++
				break
			}
		}
		if count == 0 {
			roleIds = append(roleIds, menu.ID)
		} else {
			count--
		}
	}

	return menuListFormat, roleIds, nil
}

func (m *defaultAuthorityMenuModel) FindByIds(ctx context.Context, ids []uint) ([]AuthorityMenuEntity, error) {
	return nil, nil
}

func (m *defaultAuthorityMenuModel) tableName() string {
	return m.table
}
