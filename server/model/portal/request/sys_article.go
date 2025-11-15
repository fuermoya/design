package request

import (
	"github.com/fuermoya/design/server/model/common/request"
	"github.com/fuermoya/design/server/model/portal"
)

type SysArticleSearch struct {
	portal.SysArticle
	request.PageInfo
	CategoryID *int   `json:"categoryId" form:"categoryId"`
	Status     *int   `json:"status" form:"status"`
	Keyword    string `json:"keyword" form:"keyword"`
}

type SysArticleDelete struct {
	Ids []uint `json:"ids" form:"ids"`
}

type SysArticleCreate struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary"`
	CoverImage string `json:"coverImage"`
	CategoryID uint   `json:"categoryId" binding:"required"`
	TagIDs     []uint `json:"tagIds"`
	Status     int    `json:"status"`
	Sort       int    `json:"sort"`
}

type SysArticleUpdate struct {
	ID         uint   `json:"id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary"`
	CoverImage string `json:"coverImage"`
	CategoryID uint   `json:"categoryId" binding:"required"`
	TagIDs     []uint `json:"tagIds"`
	Status     int    `json:"status"`
	Sort       int    `json:"sort"`
}
