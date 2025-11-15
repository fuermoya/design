package global

import (
	"sync"
	"time"

	"github.com/fuermoya/design/server/utils/timer"
	"github.com/songzhibin97/gkit/cache/local_cache"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	"github.com/fuermoya/design/server/config"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	DB                  *gorm.DB
	CONFIG              config.Server
	VP                  *viper.Viper
	LOG                 *zap.Logger
	Timer               timer.Timer = timer.NewTimerTask()
	Concurrency_Control             = &singleflight.Group{}
	BlackCache          local_cache.Cache
	lock                sync.RWMutex
	START_TIME          time.Time = time.Now() // 系统启动时间
)
