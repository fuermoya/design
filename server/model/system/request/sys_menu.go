package request

import (
	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/system"
)

// Add menu authority info structure
type AddMenuAuthorityInfo struct {
	Menus       []system.SysBaseMenu `json:"menus"`
	AuthorityId uint                 `json:"authorityId"` // 角色ID
}

func DefaultMenu() []system.SysBaseMenu {
	return []system.SysBaseMenu{{
		BASE_MODEL: global.BASE_MODEL{ID: 5},
		MenuLevel:  0,
		Hidden:     true,
		ParentId:   "1",
		Path:       "person/person",
		Name:       "person",
		Sort:       4,
		Meta:       system.Meta{Title: "个人信息", Icon: "message"}},
	}
}
