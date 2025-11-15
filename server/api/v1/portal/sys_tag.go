package portal

import (
	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/common/response"
	portalReq "github.com/fuermoya/design/server/model/portal/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysTagApi struct{}

// CreateSysTag 创建标签
// @Tags SysTag
// @Summary 创建标签
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body portalReq.SysTagCreate true "创建标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sysTag/createSysTag [post]
func (sysTagApi *SysTagApi) CreateSysTag(c *gin.Context) {
	var sysTag portalReq.SysTagCreate
	err := c.ShouldBindJSON(&sysTag)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysTagService.CreateSysTag(sysTag); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteSysTag 删除标签
// @Tags SysTag
// @Summary 删除标签
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body portalReq.SysTag true "删除标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysTag/deleteSysTag [delete]
func (sysTagApi *SysTagApi) DeleteSysTag(c *gin.Context) {
	ID := c.Query("id")
	if err := sysTagService.DeleteSysTag(ID); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteSysTagByIds 批量删除标签
// @Tags SysTag
// @Summary 批量删除标签
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysTag/deleteSysTag [delete]
func (sysTagApi *SysTagApi) DeleteSysTagByIds(c *gin.Context) {
	var sysTag portalReq.SysTagDelete
	_ = c.ShouldBindJSON(&sysTag)
	if err := sysTagService.DeleteSysTagByIds(sysTag.Ids); err != nil {
		global.LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateSysTag 更新标签
// @Tags SysTag
// @Summary 更新标签
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body portalReq.SysTagUpdate true "更新标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysTag/updateSysTag [put]
func (sysTagApi *SysTagApi) UpdateSysTag(c *gin.Context) {
	var sysTag portalReq.SysTagUpdate
	err := c.ShouldBindJSON(&sysTag)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysTagService.UpdateSysTag(sysTag); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindSysTag 用id查询标签
// @Tags SysTag
// @Summary 用id查询标签
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query portalReq.SysTag true "用id查询标签"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysTag/findSysTag [get]
func (sysTagApi *SysTagApi) FindSysTag(c *gin.Context) {
	ID := c.Query("ID")
	resysTag, err := sysTagService.GetSysTag(ID)
	if err != nil {
		global.LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"resysTag": resysTag}, c)
	}
}

// GetSysTagList 分页获取标签列表
// @Tags SysTag
// @Summary 分页获取标签列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.SysTagSearch true "分页获取标签列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysTag/getSysTagList [get]
func (sysTagApi *SysTagApi) GetSysTagList(c *gin.Context) {
	var pageInfo portalReq.SysTagSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := sysTagService.GetSysTagInfoList(pageInfo); err != nil {
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

// GetPortalTags 获取门户网站标签列表（前台接口）
// @Tags Portal
// @Summary 获取门户网站标签列表
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /portal/tags [get]
func (sysTagApi *SysTagApi) GetPortalTags(c *gin.Context) {
	if list, err := sysTagService.GetPortalTags(); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List: list,
		}, "获取成功", c)
	}
}
