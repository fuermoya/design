package request

import (
	"github.com/fuermoya/design/server/model/common/request"
	"github.com/fuermoya/design/server/model/portal"
)

type SysThemeSearch struct {
	portal.SysTheme
	request.PageInfo
}

type SysThemeCreate struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Config      string `json:"config"`
	Preview     string `json:"preview"`
	Sort        int    `json:"sort"`
}

type SysThemeUpdate struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Config      string `json:"config"`
	Preview     string `json:"preview"`
	Sort        int    `json:"sort"`
}

type SysThemeActivate struct {
	ID uint `json:"id" binding:"required"`
}
