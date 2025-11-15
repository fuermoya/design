package system

import (
	"github.com/fuermoya/design/server/global"
)

type SysBaseMenu struct {
	global.BASE_MODEL
	MenuLevel     uint                                       `json:"-"`
	ParentId      string                                     `json:"parentId" gorm:"comment:父菜单ID"` // 父菜单ID
	Path          string                                     `json:"path" gorm:"comment:路由path"`    // 路由path
	Name          string                                     `json:"name" gorm:"comment:路由name"`    // 路由name
	Hidden        bool                                       `json:"hidden" gorm:"comment:是否在列表隐藏"` // 是否在列表隐藏
	Sort          int                                        `json:"sort" gorm:"comment:排序标记"`      // 排序标记
	Meta          `json:"meta" gorm:"embedded;comment:附加属性"` // 附加属性
	SysAuthoritys []SysAuthority                             `json:"authoritys" gorm:"many2many:sys_authority_menus;"`
	Children      []SysBaseMenu                              `json:"children" gorm:"-"`
}

type Meta struct {
	Title string `json:"title" gorm:"comment:菜单名"` // 菜单名
	Icon  string `json:"icon" gorm:"comment:菜单图标"` // 菜单图标
}

func (SysBaseMenu) TableName() string {
	return "sys_base_menus"
}
