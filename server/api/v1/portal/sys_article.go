package portal

import (
	"strconv"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/common/response"
	"github.com/fuermoya/design/server/model/portal"
	portalReq "github.com/fuermoya/design/server/model/portal/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysArticleApi struct{}

// CreateSysArticle 创建文章
// @Tags SysArticle
// @Summary 创建文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body portalReq.SysArticleCreate true "创建文章"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sysArticle/createSysArticle [post]
func (sysArticleApi *SysArticleApi) CreateSysArticle(c *gin.Context) {
	var sysArticle portalReq.SysArticleCreate
	err := c.ShouldBindJSON(&sysArticle)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysArticleService.CreateSysArticle(sysArticle); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteSysArticle 删除文章
// @Tags SysArticle
// @Summary 删除文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body portalReq.SysArticle true "删除文章"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysArticle/deleteSysArticle [delete]
func (sysArticleApi *SysArticleApi) DeleteSysArticle(c *gin.Context) {
	ID := c.Query("id")
	if err := sysArticleService.DeleteSysArticle(ID); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteSysArticleByIds 批量删除文章
// @Tags SysArticle
// @Summary 批量删除文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysArticle/deleteSysArticle [delete]
func (sysArticleApi *SysArticleApi) DeleteSysArticleByIds(c *gin.Context) {
	var sysArticle portalReq.SysArticleDelete
	_ = c.ShouldBindJSON(&sysArticle)
	if err := sysArticleService.DeleteSysArticleByIds(sysArticle.Ids); err != nil {
		global.LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateSysArticle 更新文章
// @Tags SysArticle
// @Summary 更新文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body portalReq.SysArticleUpdate true "更新文章"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysArticle/updateSysArticle [put]
func (sysArticleApi *SysArticleApi) UpdateSysArticle(c *gin.Context) {
	var sysArticle portalReq.SysArticleUpdate
	err := c.ShouldBindJSON(&sysArticle)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysArticleService.UpdateSysArticle(sysArticle); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindSysArticle 用id查询文章
// @Tags SysArticle
// @Summary 用id查询文章
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query portalReq.SysArticle true "用id查询文章"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysArticle/findSysArticle [get]
func (sysArticleApi *SysArticleApi) FindSysArticle(c *gin.Context) {
	ID := c.Query("ID")
	resysArticle, err := sysArticleService.GetSysArticle(ID)
	if err != nil {
		global.LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"resysArticle": resysArticle}, c)
	}
}

// GetSysArticleList 分页获取文章列表
// @Tags SysArticle
// @Summary 分页获取文章列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.SysArticleSearch true "分页获取文章列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysArticle/getSysArticleList [get]
func (sysArticleApi *SysArticleApi) GetSysArticleList(c *gin.Context) {
	var pageInfo portalReq.SysArticleSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := sysArticleService.GetSysArticleInfoList(pageInfo); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetPublishedArticles 获取已发布的文章列表（前台接口）
// @Tags Portal
// @Summary 获取已发布的文章列表
// @accept application/json
// @Produce application/json
// @Param data query portalReq.SysArticleSearch true "获取已发布的文章列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /portal/articles [get]
func (sysArticleApi *SysArticleApi) GetPublishedArticles(c *gin.Context) {
	var pageInfo portalReq.SysArticleSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if list, total, err := sysArticleService.GetPublishedArticles(pageInfo); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetArticleDetail 获取文章详情（前台接口）
// @Tags Portal
// @Summary 获取文章详情
// @accept application/json
// @Produce application/json
// @Param id path int true "文章ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /portal/articles/{id} [get]
func (sysArticleApi *SysArticleApi) GetArticleDetail(c *gin.Context) {
	id := c.Param("id")

	// 获取客户端IP和User-Agent
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 增加浏览次数
	go sysArticleService.IncrementViewCount(id, ip, userAgent)

	article, err := sysArticleService.GetSysArticle(id)
	if err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	// 获取文章的浏览和点赞统计
	viewCount, _ := sysArticleService.GetArticleViewCount(id)
	likeCount, _ := sysArticleService.GetArticleLikeCount(id)

	// 检查当前IP是否已点赞
	isLiked := false
	var existingLike portal.SysArticleLike
	articleID, _ := strconv.ParseUint(id, 10, 32)
	if global.DB.Where("article_id = ? AND ip = ?", articleID, ip).First(&existingLike).Error == nil {
		isLiked = true
	}

	response.OkWithData(gin.H{
		"article":   article,
		"viewCount": viewCount,
		"likeCount": likeCount,
		"isLiked":   isLiked,
	}, c)
}

// LikeArticle 点赞文章（前台接口）
// @Tags Portal
// @Summary 点赞文章
// @accept application/json
// @Produce application/json
// @Param id path int true "文章ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"点赞成功"}"
// @Router /portal/articles/{id}/like [post]
func (sysArticleApi *SysArticleApi) LikeArticle(c *gin.Context) {
	id := c.Param("id")

	// 获取客户端IP
	ip := c.ClientIP()

	if err := sysArticleService.LikeArticle(id, ip); err != nil {
		global.LOG.Error("点赞失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("点赞成功", c)
	}
}
