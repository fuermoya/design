package portal

import (
	v1 "github.com/fuermoya/design/server/api/v1"
	"github.com/fuermoya/design/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysThemeRouter struct{}

// InitSysThemeRouter 初始化主题路由
func (s *SysThemeRouter) InitSysThemeRouter(Router *gin.RouterGroup) {
	sysThemeRouter := Router.Group("sysTheme").Use(middleware.OperationRecord())
	sysThemeRouterWithoutRecord := Router.Group("sysTheme")
	sysThemeApi := v1.ApiGroupApp.PortalApiGroup.SysThemeApi
	{
		sysThemeRouter.POST("createSysTheme", sysThemeApi.CreateSysTheme)             // 新建主题
		sysThemeRouter.DELETE("deleteSysTheme", sysThemeApi.DeleteSysTheme)           // 删除主题
		sysThemeRouter.DELETE("deleteSysThemeByIds", sysThemeApi.DeleteSysThemeByIds) // 批量删除主题
		sysThemeRouter.PUT("updateSysTheme", sysThemeApi.UpdateSysTheme)              // 更新主题
		sysThemeRouter.POST("activateTheme", sysThemeApi.ActivateTheme)               // 激活主题
	}
	{
		sysThemeRouterWithoutRecord.GET("findSysTheme", sysThemeApi.FindSysTheme)       // 根据ID获取主题
		sysThemeRouterWithoutRecord.GET("getSysThemeList", sysThemeApi.GetSysThemeList) // 获取主题列表
	}
}

// InitPortalThemeRouter 初始化门户网站主题路由（前台接口）
func (s *SysThemeRouter) InitPortalThemeRouter(Router *gin.RouterGroup) {
	portalRouter := Router.Group("portal")
	sysThemeApi := v1.ApiGroupApp.PortalApiGroup.SysThemeApi
	{
		portalRouter.GET("theme", sysThemeApi.GetActiveTheme)     // 获取当前激活的主题
		portalRouter.GET("themes", sysThemeApi.GetPortalThemes)   // 获取门户网站主题列表
		portalRouter.POST("switchTheme", sysThemeApi.SwitchTheme) // 切换主题
	}
}
