package system

import (
	v1 "github.com/fuermoya/design/server/api/v1"
	"github.com/fuermoya/design/server/middleware"
	"github.com/gin-gonic/gin"
)

type DashboardRouter struct{}

func (s *DashboardRouter) InitDashboardRouter(Router *gin.RouterGroup) {
	dashboardRouter := Router.Group("dashboard").Use(middleware.OperationRecord())
	dashboardApi := v1.ApiGroupApp.SystemApiGroup.DashboardApi
	{
		dashboardRouter.GET("stats", dashboardApi.GetDashboardStats)        // 获取统计数据
		dashboardRouter.GET("status", dashboardApi.GetSystemStatus)         // 获取系统状态
		dashboardRouter.GET("charts", dashboardApi.GetDashboardCharts)      // 获取图表数据
		dashboardRouter.GET("activities", dashboardApi.GetRecentActivities) // 获取最近活动
		dashboardRouter.GET("info", dashboardApi.GetSystemInfo)             // 获取系统信息
	}

}
