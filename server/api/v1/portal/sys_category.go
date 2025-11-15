package portal

import (
	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/common/response"
	portalReq "github.com/fuermoya/design/server/model/portal/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysCategoryApi struct{}

// CreateSysCategory 创建分类
// @Tags SysCategory
// @Summary 创建分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body portalReq.SysCategoryCreate true "创建分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sysCategory/createSysCategory [post]
func (sysCategoryApi *SysCategoryApi) CreateSysCategory(c *gin.Context) {
	var sysCategory portalReq.SysCategoryCreate
	err := c.ShouldBindJSON(&sysCategory)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysCategoryService.CreateSysCategory(sysCategory); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteSysCategory 删除分类
// @Tags SysCategory
// @Summary 删除分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body portalReq.SysCategory true "删除分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysCategory/deleteSysCategory [delete]
func (sysCategoryApi *SysCategoryApi) DeleteSysCategory(c *gin.Context) {
	ID := c.Query("id")
	if err := sysCategoryService.DeleteSysCategory(ID); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteSysCategoryByIds 批量删除分类
// @Tags SysCategory
// @Summary 批量删除分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysCategory/deleteSysCategory [delete]
func (sysCategoryApi *SysCategoryApi) DeleteSysCategoryByIds(c *gin.Context) {
	var sysCategory portalReq.SysCategoryDelete
	_ = c.ShouldBindJSON(&sysCategory)
	if err := sysCategoryService.DeleteSysCategoryByIds(sysCategory.Ids); err != nil {
		global.LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateSysCategory 更新分类
// @Tags SysCategory
// @Summary 更新分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body portalReq.SysCategoryUpdate true "更新分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysCategory/updateSysCategory [put]
func (sysCategoryApi *SysCategoryApi) UpdateSysCategory(c *gin.Context) {
	var sysCategory portalReq.SysCategoryUpdate
	err := c.ShouldBindJSON(&sysCategory)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysCategoryService.UpdateSysCategory(sysCategory); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindSysCategory 用id查询分类
// @Tags SysCategory
// @Summary 用id查询分类
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query portalReq.SysCategory true "用id查询分类"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysCategory/findSysCategory [get]
func (sysCategoryApi *SysCategoryApi) FindSysCategory(c *gin.Context) {
	ID := c.Query("ID")
	resysCategory, err := sysCategoryService.GetSysCategory(ID)
	if err != nil {
		global.LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"resysCategory": resysCategory}, c)
	}
}

// GetSysCategoryList 分页获取分类列表
// @Tags SysCategory
// @Summary 分页获取分类列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.SysCategorySearch true "分页获取分类列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysCategory/getSysCategoryList [get]
func (sysCategoryApi *SysCategoryApi) GetSysCategoryList(c *gin.Context) {
	var pageInfo portalReq.SysCategorySearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := sysCategoryService.GetSysCategoryInfoList(pageInfo); err != nil {
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

// GetPortalCategories 获取门户网站分类列表（前台接口）
// @Tags Portal
// @Summary 获取门户网站分类列表
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /portal/categories [get]
func (sysCategoryApi *SysCategoryApi) GetPortalCategories(c *gin.Context) {
	if list, err := sysCategoryService.GetPortalCategories(); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List: list,
		}, "获取成功", c)
	}
}
