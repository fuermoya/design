package router

import (
	"github.com/fuermoya/design/server/router/example"
	"github.com/fuermoya/design/server/router/portal"
	"github.com/fuermoya/design/server/router/system"
)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Portal  portal.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
