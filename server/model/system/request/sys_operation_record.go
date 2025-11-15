package request

import (
	"github.com/fuermoya/design/server/model/common/request"
	"github.com/fuermoya/design/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
