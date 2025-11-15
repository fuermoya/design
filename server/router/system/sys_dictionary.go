package system

import (
	v1 "github.com/fuermoya/design/server/api/v1"
	"github.com/fuermoya/design/server/middleware"
	"github.com/gin-gonic/gin"
)

type DictionaryRouter struct{}

func (s *DictionaryRouter) InitSysDictionaryRouter(Router *gin.RouterGroup) {
	sysDictionaryRouter := Router.Group("sysDictionary").Use(middleware.OperationRecord())
	sysDictionaryRouterWithoutRecord := Router.Group("sysDictionary")
	sysDictionaryApi := v1.ApiGroupApp.SystemApiGroup.DictionaryApi
	{
		sysDictionaryRouter.POST("createSysDictionary", sysDictionaryApi.CreateSysDictionary)   // 新建SysDictionary
		sysDictionaryRouter.DELETE("deleteSysDictionary", sysDictionaryApi.DeleteSysDictionary) // 删除SysDictionary
		sysDictionaryRouter.PUT("updateSysDictionary", sysDictionaryApi.UpdateSysDictionary)    // 更新SysDictionary
	}
	{
		sysDictionaryRouterWithoutRecord.GET("getSysDictionaryList", sysDictionaryApi.GetSysDictionaryList) // 获取SysDictionary列表
	}
}

// InitSysDictionaryJWTRouter 公共接口 只校验jwt
func (s *DictionaryRouter) InitSysDictionaryJWTRouter(Router *gin.RouterGroup) {
	sysDictionaryRouterWithoutRecord := Router.Group("sysPublicDictionary")
	sysDictionaryApi := v1.ApiGroupApp.SystemApiGroup.DictionaryApi
	{
		sysDictionaryRouterWithoutRecord.GET("findSysDictionary", sysDictionaryApi.FindSysDictionary) // 根据ID获取SysDictionary
	}
}
