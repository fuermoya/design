package portal

import (
	v1 "github.com/fuermoya/design/server/api/v1"
	"github.com/fuermoya/design/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysTagRouter struct{}

// InitSysTagRouter 初始化标签路由
func (s *SysTagRouter) InitSysTagRouter(Router *gin.RouterGroup) {
	sysTagRouter := Router.Group("sysTag").Use(middleware.OperationRecord())
	sysTagRouterWithoutRecord := Router.Group("sysTag")
	sysTagApi := v1.ApiGroupApp.PortalApiGroup.SysTagApi
	{
		sysTagRouter.POST("createSysTag", sysTagApi.CreateSysTag)             // 新建标签
		sysTagRouter.DELETE("deleteSysTag", sysTagApi.DeleteSysTag)           // 删除标签
		sysTagRouter.DELETE("deleteSysTagByIds", sysTagApi.DeleteSysTagByIds) // 批量删除标签
		sysTagRouter.PUT("updateSysTag", sysTagApi.UpdateSysTag)              // 更新标签
	}
	{
		sysTagRouterWithoutRecord.GET("findSysTag", sysTagApi.FindSysTag)       // 根据ID获取标签
		sysTagRouterWithoutRecord.GET("getSysTagList", sysTagApi.GetSysTagList) // 获取标签列表
	}
}

// InitPortalTagRouter 初始化门户网站标签路由（前台接口）
func (s *SysTagRouter) InitPortalTagRouter(Router *gin.RouterGroup) {
	portalRouter := Router.Group("portal")
	sysTagApi := v1.ApiGroupApp.PortalApiGroup.SysTagApi
	{
		portalRouter.GET("tags", sysTagApi.GetPortalTags) // 获取门户网站标签列表
	}
}
