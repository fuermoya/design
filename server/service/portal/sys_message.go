package portal

import (
	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/portal"
	portalReq "github.com/fuermoya/design/server/model/portal/request"
	portalRes "github.com/fuermoya/design/server/model/portal/response"
)

type SysMessageService struct{}

// CreateMessage 创建留言
func (messageService *SysMessageService) CreateMessage(message portal.SysMessage) (err error) {
	err = global.DB.Create(&message).Error
	return err
}

// GetMessageList 获取留言列表
func (messageService *SysMessageService) GetMessageList(info portalReq.SysMessageSearch) (list []portalRes.SysMessageResponse, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB.Model(&portal.SysMessage{})

	// 条件查询
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Email != "" {
		db = db.Where("email LIKE ?", "%"+info.Email+"%")
	}
	if info.Status != 0 {
		db = db.Where("status = ?", info.Status)
	}

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 获取列表 - 明确指定要查询的字段，避免包含虚拟字段
	err = db.Select("id, created_at, updated_at, deleted_at, name, email, phone, content, ip, status, reply, reply_time, reply_user_id").
		Limit(limit).Offset(offset).Order("created_at DESC").Preload("ReplyUser").Find(&list).Error

	// 设置状态文本
	for i := range list {
		list[i].StatusText = getStatusText(list[i].Status)
	}

	return list, total, err
}

// GetMessageById 根据ID获取留言
func (messageService *SysMessageService) GetMessageById(id uint) (message portalRes.SysMessageResponse, err error) {
	err = global.DB.Select("id, created_at, updated_at, deleted_at, name, email, phone, content, ip, status, reply, reply_time, reply_user_id").
		Preload("ReplyUser").First(&message, id).Error
	if err != nil {
		return
	}
	message.StatusText = getStatusText(message.Status)
	return message, err
}

// UpdateMessage 更新留言
func (messageService *SysMessageService) UpdateMessage(message portal.SysMessage) (err error) {
	err = global.DB.Model(&message).Updates(message).Error
	return err
}

// DeleteMessage 删除留言
func (messageService *SysMessageService) DeleteMessage(message portal.SysMessage) (err error) {
	err = global.DB.Delete(&message).Error
	return err
}

// GetMessageStats 获取留言统计
func (messageService *SysMessageService) GetMessageStats() (stats map[string]interface{}, err error) {
	var total, unread, read, replied int64

	// 总数
	err = global.DB.Model(&portal.SysMessage{}).Count(&total).Error
	if err != nil {
		return
	}

	// 未读
	err = global.DB.Model(&portal.SysMessage{}).Where("status = ?", 2).Count(&unread).Error
	if err != nil {
		return
	}

	// 已读
	err = global.DB.Model(&portal.SysMessage{}).Where("status = ?", 1).Count(&read).Error
	if err != nil {
		return
	}

	// 已回复
	err = global.DB.Model(&portal.SysMessage{}).Where("status = ?", 3).Count(&replied).Error
	if err != nil {
		return
	}

	stats = map[string]interface{}{
		"total":   total,
		"unread":  unread,
		"read":    read,
		"replied": replied,
	}

	return stats, nil
}

// getStatusText 获取状态文本
func getStatusText(status int) string {
	switch status {
	case 1:
		return "已读"
	case 2:
		return "未读"
	case 3:
		return "已回复"
	default:
		return "未知"
	}
}
