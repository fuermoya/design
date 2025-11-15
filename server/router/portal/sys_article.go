package portal

import (
	v1 "github.com/fuermoya/design/server/api/v1"
	"github.com/fuermoya/design/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysArticleRouter struct{}

// InitSysArticleRouter 初始化文章路由
func (s *SysArticleRouter) InitSysArticleRouter(Router *gin.RouterGroup) {
	sysArticleRouter := Router.Group("sysArticle").Use(middleware.OperationRecord())
	sysArticleRouterWithoutRecord := Router.Group("sysArticle")
	sysArticleApi := v1.ApiGroupApp.PortalApiGroup.SysArticleApi
	{
		sysArticleRouter.POST("createSysArticle", sysArticleApi.CreateSysArticle)             // 新建文章
		sysArticleRouter.DELETE("deleteSysArticle", sysArticleApi.DeleteSysArticle)           // 删除文章
		sysArticleRouter.DELETE("deleteSysArticleByIds", sysArticleApi.DeleteSysArticleByIds) // 批量删除文章
		sysArticleRouter.PUT("updateSysArticle", sysArticleApi.UpdateSysArticle)              // 更新文章
	}
	{
		sysArticleRouterWithoutRecord.GET("findSysArticle", sysArticleApi.FindSysArticle)       // 根据ID获取文章
		sysArticleRouterWithoutRecord.GET("getSysArticleList", sysArticleApi.GetSysArticleList) // 获取文章列表
	}
}

// InitPortalRouter 初始化门户网站路由（前台接口）
func (s *SysArticleRouter) InitPortalRouter(Router *gin.RouterGroup) {
	portalRouter := Router.Group("portal")
	sysArticleApi := v1.ApiGroupApp.PortalApiGroup.SysArticleApi
	{
		portalRouter.GET("articles", sysArticleApi.GetPublishedArticles)  // 获取已发布的文章列表
		portalRouter.GET("articles/:id", sysArticleApi.GetArticleDetail)  // 获取文章详情
		portalRouter.POST("articles/:id/like", sysArticleApi.LikeArticle) // 点赞文章

	}
}
