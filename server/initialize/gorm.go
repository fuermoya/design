package initialize

import (
	"os"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/model/example"
	"github.com/fuermoya/design/server/model/portal"
	"github.com/fuermoya/design/server/model/system"
	"github.com/fuermoya/design/server/service"
	_ "github.com/fuermoya/design/server/source/example"
	_ "github.com/fuermoya/design/server/source/portal"
	_ "github.com/fuermoya/design/server/source/system"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	return GormMysql()
}

func InitData() {
	err := service.ServiceGroupApp.SystemServiceGroup.InitDBService.AutoInitData()
	if err != nil {
		global.LOG.Error("Create table Data failed", zap.Error(err))
		os.Exit(0)
	}
}

func RegisterTables() {
	db := global.DB
	err := db.AutoMigrate(

		system.SysApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysDictionaryDetail{},
		adapter.CasbinRule{},
		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		example.ExaFileUploadAndDownload{},

		portal.SysArticle{},
		portal.SysArticleView{},
		portal.SysArticleLike{},
		portal.SysCategory{},
		portal.SysTag{},
		portal.SysTheme{},
		portal.SysMessage{},
	)
	if err != nil {
		global.LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.LOG.Info("register table success")
}
