package v1

import (
	"github.com/fuermoya/design/server/api/v1/example"
	"github.com/fuermoya/design/server/api/v1/portal"
	"github.com/fuermoya/design/server/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
	PortalApiGroup  portal.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
