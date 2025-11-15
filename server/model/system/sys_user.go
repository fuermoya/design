package system

import (
	"github.com/fuermoya/design/server/global"
	"github.com/google/uuid"
)

type SysUser struct {
	global.BASE_MODEL
	UUID        uuid.UUID    `json:"uuid" gorm:"index;comment:用户UUID"`                                                     // 用户UUID
	Username    string       `json:"userName" gorm:"index;comment:用户登录名"`                                                  // 用户登录名
	Password    string       `json:"-"  gorm:"comment:用户登录密码"`                                                             // 用户登录密码
	NickName    string       `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                            // 用户昵称
	HeaderImg   string       `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	AuthorityId uint         `json:"authorityId" gorm:"default:888;comment:用户角色ID"`                                        // 用户角色ID
	Authority   SysAuthority `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	Enable      int          `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
}

func (SysUser) TableName() string {
	return "sys_users"
}
