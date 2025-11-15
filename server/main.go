package main

import (
	"go.uber.org/zap"

	"github.com/fuermoya/design/server/core"
	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/initialize"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title                       design Swagger API接口文档
// @version                     v2.6.1
// @description                 使用gin+vue进行极速开发的全栈开发基础平台
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {
	global.VP = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.LOG)
	global.DB = initialize.Gorm() // gorm连接数据库
	initialize.Timer()
	if global.DB != nil {
		initialize.RegisterTables() // 初始化表
		initialize.InitData()       //初始化数据
		// 程序结束前关闭数据库链接
		db, _ := global.DB.DB()
		defer db.Close()
	}
	core.RunWindowsServer()
}
