package portal

import (
	"time"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/system"
)

// SysMessage 留言表
type SysMessage struct {
	global.BASE_MODEL
	Name        string         `json:"name" form:"name" gorm:"column:name;comment:留言者姓名;size:50;"`
	Email       string         `json:"email" form:"email" gorm:"column:email;comment:留言者邮箱;size:100;"`
	Phone       string         `json:"phone" form:"phone" gorm:"column:phone;comment:留言者电话;size:20;"`
	Content     string         `json:"content" form:"content" gorm:"column:content;comment:留言内容;type:text;"`
	IP          string         `json:"ip" form:"ip" gorm:"column:ip;comment:留言者IP地址;size:50;"`
	Status      int            `json:"status" form:"status" gorm:"column:status;comment:状态 1:已读 2:未读 3:已回复;default:2"`
	Reply       string         `json:"reply" form:"reply" gorm:"column:reply;comment:回复内容;type:text;"`
	ReplyTime   *time.Time     `json:"replyTime" form:"replyTime" gorm:"column:reply_time;comment:回复时间;"`
	ReplyUserID uint           `json:"replyUserId" form:"replyUserId" gorm:"column:reply_user_id;comment:回复用户ID;"`
	ReplyUser   system.SysUser `json:"replyUser" gorm:"foreignKey:ReplyUserID"`
}

func (SysMessage) TableName() string {
	return "sys_messages"
}
