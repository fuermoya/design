package request

import (
	"github.com/fuermoya/design/server/model/common/request"
	"github.com/fuermoya/design/server/model/portal"
)

type SysCategorySearch struct {
	portal.SysCategory
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
}

type SysCategoryCreate struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
}

type SysCategoryUpdate struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	Status      int    `json:"status"`
}

type SysCategoryDelete struct {
	Ids []uint `json:"ids" form:"ids"`
}
