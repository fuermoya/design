package system

import (
	"runtime"
	"time"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/common/response"
	portalModel "github.com/fuermoya/design/server/model/portal"
	"github.com/fuermoya/design/server/model/system"
	"github.com/fuermoya/design/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DashboardApi struct{}

var sysArticleService = service.ServiceGroupApp.PortalServiceGroup.SysArticleService

// GetDashboardStats 获取仪表盘统计数据
// @Tags     Dashboard
// @Summary  获取仪表盘统计数据
// @Produce  application/json
// @Success  200  {object}  response.Response{data=map[string]interface{},msg=string}  "返回统计数据"
// @Router   /system/dashboard/stats [get]
func (d *DashboardApi) GetDashboardStats(c *gin.Context) {
	var stats = make(map[string]interface{})

	// 获取用户总数
	var userCount int64
	global.DB.Model(&system.SysUser{}).Count(&userCount)
	stats["totalUsers"] = userCount

	// 获取文章总数
	var articleCount int64
	global.DB.Model(&portalModel.SysArticle{}).Count(&articleCount)
	stats["totalArticles"] = articleCount

	// 获取总浏览量（从浏览记录表统计）
	totalViews, err := sysArticleService.GetTotalViewCount()
	if err != nil {
		global.LOG.Error("获取总浏览量失败", zap.Error(err))
		totalViews = 0
	}
	stats["totalViews"] = totalViews

	// 获取总点赞数（从点赞记录表统计）
	totalLikes, err := sysArticleService.GetTotalLikeCount()
	if err != nil {
		global.LOG.Error("获取总点赞数失败", zap.Error(err))
		totalLikes = 0
	}
	stats["totalLikes"] = totalLikes

	// 获取留言总数
	var messageCount int64
	global.DB.Model(&portalModel.SysMessage{}).Count(&messageCount)
	stats["totalMessages"] = messageCount

	// 获取今天和昨天的开始和结束时间
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := today.Add(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	yesterdayEnd := today

	// 获取今天的数据量
	var todayStats = make(map[string]int64)

	// 今天新增文章数
	var todayArticleCount int64
	global.DB.Model(&portalModel.SysArticle{}).Where("created_at >= ? AND created_at < ?", today, todayEnd).Count(&todayArticleCount)
	todayStats["articles"] = todayArticleCount

	// 今天新增浏览量
	todayViews, err := sysArticleService.GetTodayViewCount()
	if err != nil {
		global.LOG.Error("获取今天浏览量失败", zap.Error(err))
		todayViews = 0
	}
	todayStats["views"] = todayViews

	// 今天新增点赞数
	todayLikes, err := sysArticleService.GetTodayLikeCount()
	if err != nil {
		global.LOG.Error("获取今天点赞数失败", zap.Error(err))
		todayLikes = 0
	}
	todayStats["likes"] = todayLikes

	// 今天新增留言数
	var todayMessageCount int64
	global.DB.Model(&portalModel.SysMessage{}).Where("created_at >= ? AND created_at < ?", today, todayEnd).Count(&todayMessageCount)
	todayStats["messages"] = todayMessageCount

	// 获取昨天的数据量
	var yesterdayStats = make(map[string]int64)

	// 昨天新增文章数
	var yesterdayArticleCount int64
	global.DB.Model(&portalModel.SysArticle{}).Where("created_at >= ? AND created_at < ?", yesterday, yesterdayEnd).Count(&yesterdayArticleCount)
	yesterdayStats["articles"] = yesterdayArticleCount

	// 昨天新增浏览量
	yesterdayViews, err := sysArticleService.GetYesterdayViewCount()
	if err != nil {
		global.LOG.Error("获取昨天浏览量失败", zap.Error(err))
		yesterdayViews = 0
	}
	yesterdayStats["views"] = yesterdayViews

	// 昨天新增点赞数
	yesterdayLikes, err := sysArticleService.GetYesterdayLikeCount()
	if err != nil {
		global.LOG.Error("获取昨天点赞数失败", zap.Error(err))
		yesterdayLikes = 0
	}
	yesterdayStats["likes"] = yesterdayLikes

	// 昨天新增留言数
	var yesterdayMessageCount int64
	global.DB.Model(&portalModel.SysMessage{}).Where("created_at >= ? AND created_at < ?", yesterday, yesterdayEnd).Count(&yesterdayMessageCount)
	yesterdayStats["messages"] = yesterdayMessageCount

	// 计算百分比变化（今天 vs 昨天）
	calculatePercentageChange := func(today, yesterday int64) float64 {
		if yesterday == 0 {
			if today == 0 {
				return 0
			}
			return 100 // 如果昨天为0，今天有数据，则增长100%
		}
		return float64(today-yesterday) / float64(yesterday) * 100
	}

	// 添加百分比变化数据
	stats["articleChange"] = calculatePercentageChange(todayStats["articles"], yesterdayStats["articles"])
	stats["viewChange"] = calculatePercentageChange(todayStats["views"], yesterdayStats["views"])
	stats["likeChange"] = calculatePercentageChange(todayStats["likes"], yesterdayStats["likes"])
	stats["messageChange"] = calculatePercentageChange(todayStats["messages"], yesterdayStats["messages"])

	response.OkWithData(stats, c)
}

// GetSystemStatus 获取系统状态
// @Tags     Dashboard
// @Summary  获取系统状态
// @Produce  application/json
// @Success  200  {object}  response.Response{data=map[string]interface{},msg=string}  "返回系统状态"
// @Router   /system/dashboard/status [get]
func (d *DashboardApi) GetSystemStatus(c *gin.Context) {
	var status = make(map[string]interface{})

	// 获取系统运行时间
	uptime := time.Since(global.START_TIME).Seconds()
	status["uptime"] = int64(uptime)

	// 获取内存使用情况
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 计算内存使用率 (已分配内存 / 总内存)
	memoryUsage := float64(m.Alloc) / float64(m.Sys) * 100
	status["memory"] = int(memoryUsage)

	// 模拟CPU使用率 (实际项目中应该通过系统调用获取)
	// 这里使用一个简单的模拟值，实际项目中可以使用gopsutil等库
	status["cpu"] = 25 // 模拟25%的CPU使用率

	// 模拟磁盘使用率
	status["disk"] = 45 // 模拟45%的磁盘使用率

	response.OkWithData(status, c)
}

// GetDashboardCharts 获取图表数据
// @Tags     Dashboard
// @Summary  获取图表数据
// @Produce  application/json
// @Success  200  {object}  response.Response{data=map[string]interface{},msg=string}  "返回图表数据"
// @Router   /system/dashboard/charts [get]
func (d *DashboardApi) GetDashboardCharts(c *gin.Context) {
	var charts = make(map[string]interface{})

	// 获取最近7天的文章发布统计
	var articleStats []map[string]interface{}
	rows, err := global.DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM sys_articles 
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)
		GROUP BY DATE(created_at)
		ORDER BY date
	`).Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var date string
			var count int
			rows.Scan(&date, &count)
			articleStats = append(articleStats, map[string]interface{}{
				"date":  date,
				"count": count,
			})
		}
	}
	charts["articles"] = articleStats

	// 获取最近7天的用户注册统计
	var userStats []map[string]interface{}
	rows, err = global.DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM sys_users 
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)
		GROUP BY DATE(created_at)
		ORDER BY date
	`).Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var date string
			var count int
			rows.Scan(&date, &count)
			userStats = append(userStats, map[string]interface{}{
				"date":  date,
				"count": count,
			})
		}
	}
	charts["users"] = userStats

	// 获取最近7天的浏览量统计
	var viewStats []map[string]interface{}
	rows, err = global.DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			SUM(view_count) as total_views
		FROM sys_articles 
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)
		GROUP BY DATE(created_at)
		ORDER BY date
	`).Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var date string
			var totalViews int
			rows.Scan(&date, &totalViews)
			viewStats = append(viewStats, map[string]interface{}{
				"date":       date,
				"totalViews": totalViews,
			})
		}
	}
	charts["views"] = viewStats

	response.OkWithData(charts, c)
}

// GetRecentActivities 获取最近活动
// @Tags     Dashboard
// @Summary  获取最近活动
// @Produce  application/json
// @Success  200  {object}  response.Response{data=[]map[string]interface{},msg=string}  "返回最近活动"
// @Router   /system/dashboard/activities [get]
func (d *DashboardApi) GetRecentActivities(c *gin.Context) {
	var activities []map[string]interface{}

	// 获取最近的操作记录
	var records []system.SysOperationRecord
	global.DB.Order("created_at desc").Limit(10).Find(&records)

	for _, record := range records {
		activities = append(activities, map[string]interface{}{
			"id":        record.ID,
			"method":    record.Method,
			"path":      record.Path,
			"status":    record.Status,
			"user":      record.User,
			"ip":        record.Ip,
			"createdAt": record.CreatedAt,
		})
	}

	response.OkWithData(activities, c)
}

// GetSystemInfo 获取系统信息
// @Tags     Dashboard
// @Summary  获取系统信息
// @Produce  application/json
// @Success  200  {object}  response.Response{data=map[string]interface{},msg=string}  "返回系统信息"
// @Router   /system/dashboard/info [get]
func (d *DashboardApi) GetSystemInfo(c *gin.Context) {
	var info = make(map[string]interface{})

	// 获取数据库连接信息
	var dbInfo map[string]interface{}
	global.DB.Raw("SELECT VERSION() as version").Scan(&dbInfo)

	// 获取系统信息
	info["goVersion"] = runtime.Version()
	info["os"] = runtime.GOOS
	info["arch"] = runtime.GOARCH
	info["numCPU"] = runtime.NumCPU()
	info["dbVersion"] = dbInfo["version"]
	info["startTime"] = global.START_TIME.Format("2006-01-02 15:04:05")

	response.OkWithData(info, c)
}
