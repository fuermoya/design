package portal

import (
	v1 "github.com/fuermoya/design/server/api/v1"
	"github.com/fuermoya/design/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysMessageRouter struct{}

// InitSysMessageRouter 初始化留言路由（后台管理）
func (s *SysMessageRouter) InitSysMessageRouter(Router *gin.RouterGroup) {
	sysMessageRouter := Router.Group("sysMessage").Use(middleware.OperationRecord())
	sysMessageRouterWithoutRecord := Router.Group("sysMessage")
	messageApi := v1.ApiGroupApp.PortalApiGroup.SysMessageApi
	{
		sysMessageRouter.DELETE("deleteMessage", messageApi.DeleteMessage) // 删除留言
		sysMessageRouter.PUT("updateMessage", messageApi.UpdateMessage)    // 更新留言
		sysMessageRouter.POST("replyMessage", messageApi.ReplyMessage)     // 回复留言
		sysMessageRouter.POST("markAsRead", messageApi.MarkAsRead)         // 标记为已读
	}
	{
		sysMessageRouterWithoutRecord.GET("getMessageList", messageApi.GetMessageList)   // 获取留言列表
		sysMessageRouterWithoutRecord.GET("getMessageById", messageApi.GetMessageById)   // 根据ID获取留言
		sysMessageRouterWithoutRecord.GET("getMessageStats", messageApi.GetMessageStats) // 获取留言统计
	}
}

// InitPortalMessageRouter 初始化门户留言路由（前台接口）
func (s *SysMessageRouter) InitPortalMessageRouter(Router *gin.RouterGroup) {
	portalRouter := Router.Group("portal")
	messageApi := v1.ApiGroupApp.PortalApiGroup.SysMessageApi
	{
		portalRouter.POST("message", messageApi.CreateMessage) // 创建留言
	}
}
