package portal

import (
	"strconv"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/common/response"
	portalReq "github.com/fuermoya/design/server/model/portal/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysThemeApi struct{}

// CreateSysTheme 创建主题
// @Tags SysTheme
// @Summary 创建主题
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body portalReq.SysThemeCreate true "创建主题"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /sysTheme/createSysTheme [post]
func (sysThemeApi *SysThemeApi) CreateSysTheme(c *gin.Context) {
	var sysTheme portalReq.SysThemeCreate
	err := c.ShouldBindJSON(&sysTheme)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysThemeService.CreateSysTheme(sysTheme); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteSysTheme 删除主题
// @Tags SysTheme
// @Summary 删除主题
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body portalReq.SysTheme true "删除主题"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysTheme/deleteSysTheme [delete]
func (sysThemeApi *SysThemeApi) DeleteSysTheme(c *gin.Context) {
	ID := c.Query("ID")
	if err := sysThemeService.DeleteSysTheme(ID); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteSysThemeByIds 批量删除主题
// @Tags SysTheme
// @Summary 批量删除主题
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /sysTheme/deleteSysTheme [delete]
func (sysThemeApi *SysThemeApi) DeleteSysThemeByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := sysThemeService.DeleteSysThemeByIds(IDs); err != nil {
		global.LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateSysTheme 更新主题
// @Tags SysTheme
// @Summary 更新主题
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body portalReq.SysThemeUpdate true "更新主题"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /sysTheme/updateSysTheme [put]
func (sysThemeApi *SysThemeApi) UpdateSysTheme(c *gin.Context) {
	var sysTheme portalReq.SysThemeUpdate
	err := c.ShouldBindJSON(&sysTheme)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysThemeService.UpdateSysTheme(sysTheme); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindSysTheme 用id查询主题
// @Tags SysTheme
// @Summary 用id查询主题
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query portalReq.SysTheme true "用id查询主题"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /sysTheme/findSysTheme [get]
func (sysThemeApi *SysThemeApi) FindSysTheme(c *gin.Context) {
	ID := c.Query("ID")
	resysTheme, err := sysThemeService.GetSysTheme(ID)
	if err != nil {
		global.LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"resysTheme": resysTheme}, c)
	}
}

// GetSysThemeList 分页获取主题列表
// @Tags SysTheme
// @Summary 分页获取主题列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query portalReq.SysThemeSearch true "分页获取主题列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /sysTheme/getSysThemeList [get]
func (sysThemeApi *SysThemeApi) GetSysThemeList(c *gin.Context) {
	var pageInfo portalReq.SysThemeSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := sysThemeService.GetSysThemeInfoList(pageInfo); err != nil {
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

// ActivateTheme 激活主题
// @Tags SysTheme
// @Summary 激活主题
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body portalReq.SysThemeActivate true "激活主题"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"激活成功"}"
// @Router /sysTheme/activateTheme [post]
func (sysThemeApi *SysThemeApi) ActivateTheme(c *gin.Context) {
	var req portalReq.SysThemeActivate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sysThemeService.ActivateTheme(string(rune(req.ID))); err != nil {
		global.LOG.Error("激活失败!", zap.Error(err))
		response.FailWithMessage("激活失败", c)
	} else {
		response.OkWithMessage("激活成功", c)
	}
}

// GetActiveTheme 获取当前激活的主题（前台接口）
// @Tags Portal
// @Summary 获取当前激活的主题
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /portal/theme [get]
func (sysThemeApi *SysThemeApi) GetActiveTheme(c *gin.Context) {
	theme, err := sysThemeService.GetActiveTheme()
	if err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithData(gin.H{"theme": theme}, c)
	}
}

// GetPortalThemes 获取门户网站主题列表（前台接口）
// @Tags Portal
// @Summary 获取门户网站主题列表
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /portal/themes [get]
func (sysThemeApi *SysThemeApi) GetPortalThemes(c *gin.Context) {
	if list, err := sysThemeService.GetPortalThemes(); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List: list,
		}, "获取成功", c)
	}
}

// SwitchTheme 切换主题（前台接口）
// @Tags Portal
// @Summary 切换主题
// @accept application/json
// @Produce application/json
// @Param data body map[string]interface{} true "主题ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"切换成功"}"
// @Router /portal/switchTheme [post]
func (sysThemeApi *SysThemeApi) SwitchTheme(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	themeId, ok := req["themeId"].(float64)
	if !ok {
		response.FailWithMessage("主题ID无效", c)
		return
	}

	if err := sysThemeService.ActivateTheme(strconv.Itoa(int(themeId))); err != nil {
		global.LOG.Error("切换失败!", zap.Error(err))
		response.FailWithMessage("切换失败", c)
	} else {
		response.OkWithMessage("切换成功", c)
	}
}
