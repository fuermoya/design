package portal

import (
	v1 "github.com/fuermoya/design/server/api/v1"
	"github.com/fuermoya/design/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysCategoryRouter struct{}

// InitSysCategoryRouter 初始化分类路由
func (s *SysCategoryRouter) InitSysCategoryRouter(Router *gin.RouterGroup) {
	sysCategoryRouter := Router.Group("sysCategory").Use(middleware.OperationRecord())
	sysCategoryRouterWithoutRecord := Router.Group("sysCategory")
	sysCategoryApi := v1.ApiGroupApp.PortalApiGroup.SysCategoryApi
	{
		sysCategoryRouter.POST("createSysCategory", sysCategoryApi.CreateSysCategory)             // 新建分类
		sysCategoryRouter.DELETE("deleteSysCategory", sysCategoryApi.DeleteSysCategory)           // 删除分类
		sysCategoryRouter.DELETE("deleteSysCategoryByIds", sysCategoryApi.DeleteSysCategoryByIds) // 批量删除分类
		sysCategoryRouter.PUT("updateSysCategory", sysCategoryApi.UpdateSysCategory)              // 更新分类
	}
	{
		sysCategoryRouterWithoutRecord.GET("findSysCategory", sysCategoryApi.FindSysCategory)       // 根据ID获取分类
		sysCategoryRouterWithoutRecord.GET("getSysCategoryList", sysCategoryApi.GetSysCategoryList) // 获取分类列表
	}
}

// InitPortalCategoryRouter 初始化门户网站分类路由（前台接口）
func (s *SysCategoryRouter) InitPortalCategoryRouter(Router *gin.RouterGroup) {
	portalRouter := Router.Group("portal")
	sysCategoryApi := v1.ApiGroupApp.PortalApiGroup.SysCategoryApi
	{
		portalRouter.GET("categories", sysCategoryApi.GetPortalCategories) // 获取门户网站分类列表
	}
}
