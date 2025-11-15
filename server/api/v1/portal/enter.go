package portal

import "github.com/fuermoya/design/server/service"

type ApiGroup struct {
	SysArticleApi
	SysCategoryApi
	SysTagApi
	SysThemeApi
	SysMessageApi
}

var (
	// 统一管理所有service调用
	sysArticleService  = service.ServiceGroupApp.PortalServiceGroup.SysArticleService
	sysCategoryService = service.ServiceGroupApp.PortalServiceGroup.SysCategoryService
	sysTagService      = service.ServiceGroupApp.PortalServiceGroup.SysTagService
	sysThemeService    = service.ServiceGroupApp.PortalServiceGroup.SysThemeService
	sysMessageService  = service.ServiceGroupApp.PortalServiceGroup.SysMessageService
)
