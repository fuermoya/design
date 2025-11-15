package system

import "github.com/fuermoya/design/server/service"

type ApiGroup struct {
	JwtApi
	BaseApi
	CasbinApi
	AuthorityApi
	AuthorityMenuApi
	DictionaryApi
	OperationRecordApi
	DictionaryDetailApi
	SystemApiApi
	DashboardApi
}

var (
	jwtService              = service.ServiceGroupApp.SystemServiceGroup.JwtService
	menuService             = service.ServiceGroupApp.SystemServiceGroup.MenuService
	userService             = service.ServiceGroupApp.SystemServiceGroup.UserService
	casbinService           = service.ServiceGroupApp.SystemServiceGroup.CasbinService
	authorityService        = service.ServiceGroupApp.SystemServiceGroup.AuthorityService
	dictionaryService       = service.ServiceGroupApp.SystemServiceGroup.DictionaryService
	operationRecordService  = service.ServiceGroupApp.SystemServiceGroup.OperationRecordService
	dictionaryDetailService = service.ServiceGroupApp.SystemServiceGroup.DictionaryDetailService
	apiService              = service.ServiceGroupApp.SystemServiceGroup.ApiService
)
