package request

import (
	"github.com/fuermoya/design/server/model/common/request"
	"github.com/fuermoya/design/server/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
