package example

import (
	"github.com/fuermoya/design/server/model/common/request"
)

type ExaFileUploadAndDownloadSearch struct {
	request.PageInfo
	Name string `json:"name" form:"name"` // 文件名
	Tag  string `json:"tag" form:"tag"`   // 文件标签
}
