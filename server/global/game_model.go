package global

import (
	"time"
)

type GAME_MODEL struct {
	ID          uint      `gorm:"primarykey"` // 主键ID
	PlatformId  string    `json:"platformId" gorm:"column:platform_id;index:idx_time_platform_account,priority:2;index:idx_platform_server_player,priority:3;comment:平台id"`
	ServerId    int32     `json:"serverId" gorm:"column:server_id;comment:服务器id;index:idx_platform_server_player,priority:2;"`
	Account     string    `json:"account" gorm:"column:account;index:idx_time_platform_account,priority:1;comment:玩家账号"`
	PlayerId    int32     `json:"playerId" gorm:"column:player_id;index:idx_platform_server_player,priority:1;comment:玩家id"`
	PlayerName  string    `json:"playerName" gorm:"column:player_name;comment:玩家名称"`
	PlayerLevel int32     `json:"playerLevel" gorm:"column:player_level;comment:玩家等级"`
	CreatedAt   time.Time `json:"CreatedAt" gorm:"column:created_at;index;comment:创建时间"` // 创建时间
}
