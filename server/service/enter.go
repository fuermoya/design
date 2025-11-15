package service

import (
	"github.com/fuermoya/design/server/service/example"
	"github.com/fuermoya/design/server/service/portal"
	"github.com/fuermoya/design/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	PortalServiceGroup  portal.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
