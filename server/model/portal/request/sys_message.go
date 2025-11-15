package request

import (
	"github.com/fuermoya/design/server/model/common/request"
	"github.com/fuermoya/design/server/model/portal"
)

// SysMessageSearch 留言搜索条件
type SysMessageSearch struct {
	portal.SysMessage
	request.PageInfo
}

// SysMessageCreate 创建留言请求
type SysMessageCreate struct {
	Name    string `json:"name" binding:"required"`        // 留言者姓名
	Email   string `json:"email" binding:"required,email"` // 留言者邮箱
	Phone   string `json:"phone"`                          // 留言者电话
	Content string `json:"content" binding:"required"`     // 留言内容
}

// SysMessageReply 回复留言请求
type SysMessageReply struct {
	ID    uint   `json:"id" binding:"required"`    // 留言ID
	Reply string `json:"reply" binding:"required"` // 回复内容
}

// SysMessageMarkRead 标记已读请求
type SysMessageMarkRead struct {
	ID uint `json:"id" binding:"required"` // 留言ID
}
