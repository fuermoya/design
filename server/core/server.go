package core

import (
	"fmt"

	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/initialize"
	"github.com/fuermoya/design/server/service/system"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	// 从db加载jwt数据
	if global.DB != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)
	s := initServer(address, Router)
	global.LOG.Info("server run success on ", zap.String("address", address))
	global.LOG.Error(s.ListenAndServe().Error())
}
