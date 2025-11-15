package system

import (
	"github.com/fuermoya/design/server/global"
)

type JwtBlacklist struct {
	global.BASE_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
