package request

import (
	"github.com/fuermoya/design/server/model/common/request"
	"github.com/fuermoya/design/server/model/portal"
)

type SysTagSearch struct {
	portal.SysTag
	request.PageInfo
	Keyword string `json:"keyword" form:"keyword"`
}

type SysTagCreate struct {
	Name   string `json:"name" binding:"required"`
	Color  string `json:"color"`
	Status int    `json:"status"`
}

type SysTagUpdate struct {
	ID     uint   `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Color  string `json:"color"`
	Status int    `json:"status"`
}

type SysTagDelete struct {
	Ids []uint `json:"ids" form:"ids"`
}
