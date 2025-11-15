package portal

import (
	"strconv"
	"time"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/common/response"
	"github.com/fuermoya/design/server/model/portal"
	portalReq "github.com/fuermoya/design/server/model/portal/request"
	"github.com/fuermoya/design/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SysMessageApi struct{}

// CreateMessage 创建留言
func (m *SysMessageApi) CreateMessage(c *gin.Context) {
	var req portalReq.SysMessageCreate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 获取客户端IP
	clientIP := c.ClientIP()

	message := portal.SysMessage{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Content: req.Content,
		IP:      clientIP,
		Status:  2, // 未读状态
	}

	err = sysMessageService.CreateMessage(message)
	if err != nil {
		global.LOG.Error("创建留言失败!", zap.Error(err))
		response.FailWithMessage("创建留言失败", c)
		return
	}

	response.OkWithMessage("留言提交成功", c)
}

// GetMessageList 获取留言列表
func (m *SysMessageApi) GetMessageList(c *gin.Context) {
	var pageInfo portalReq.SysMessageSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	list, total, err := sysMessageService.GetMessageList(pageInfo)
	if err != nil {
		global.LOG.Error("获取留言列表失败!", zap.Error(err))
		response.FailWithMessage("获取留言列表失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetMessageById 根据ID获取留言
func (m *SysMessageApi) GetMessageById(c *gin.Context) {
	id := c.Query("id")
	messageId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	message, err := sysMessageService.GetMessageById(uint(messageId))
	if err != nil {
		global.LOG.Error("获取留言失败!", zap.Error(err))
		response.FailWithMessage("获取留言失败", c)
		return
	}

	response.OkWithData(message, c)
}

// UpdateMessage 更新留言
func (m *SysMessageApi) UpdateMessage(c *gin.Context) {
	var message portal.SysMessage
	err := c.ShouldBindJSON(&message)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err = sysMessageService.UpdateMessage(message)
	if err != nil {
		global.LOG.Error("更新留言失败!", zap.Error(err))
		response.FailWithMessage("更新留言失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// DeleteMessage 删除留言
func (m *SysMessageApi) DeleteMessage(c *gin.Context) {
	var message portal.SysMessage
	err := c.ShouldBindJSON(&message)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	err = sysMessageService.DeleteMessage(message)
	if err != nil {
		global.LOG.Error("删除留言失败!", zap.Error(err))
		response.FailWithMessage("删除留言失败", c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// ReplyMessage 回复留言
func (m *SysMessageApi) ReplyMessage(c *gin.Context) {
	var req portalReq.SysMessageReply
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	// 获取当前用户信息
	userID := utils.GetUserID(c)
	message := portal.SysMessage{

		Reply:       req.Reply,
		ReplyUserID: userID,

		Status: 3, // 已回复状态
	}

	now := time.Now()
	message.ReplyTime = &now
	message.ID = req.ID

	err = sysMessageService.UpdateMessage(message)
	if err != nil {
		global.LOG.Error("回复留言失败!", zap.Error(err))
		response.FailWithMessage("回复留言失败", c)
		return
	}

	response.OkWithMessage("回复成功", c)
}

// MarkAsRead 标记为已读
func (m *SysMessageApi) MarkAsRead(c *gin.Context) {
	var req portalReq.SysMessageMarkRead
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	message := portal.SysMessage{
		Status: 1, // 已读状态
	}
	message.ID = req.ID
	err = sysMessageService.UpdateMessage(message)
	if err != nil {
		global.LOG.Error("标记已读失败!", zap.Error(err))
		response.FailWithMessage("标记已读失败", c)
		return
	}

	response.OkWithMessage("标记成功", c)
}

// GetMessageStats 获取留言统计
func (m *SysMessageApi) GetMessageStats(c *gin.Context) {
	stats, err := sysMessageService.GetMessageStats()
	if err != nil {
		global.LOG.Error("获取留言统计失败!", zap.Error(err))
		response.FailWithMessage("获取留言统计失败", c)
		return
	}

	response.OkWithData(stats, c)
}
