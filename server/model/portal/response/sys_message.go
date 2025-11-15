package response

import "github.com/fuermoya/design/server/model/portal"

// SysMessageResponse 留言响应
type SysMessageResponse struct {
	portal.SysMessage
	StatusText string `json:"statusText"` // 状态文本
}
