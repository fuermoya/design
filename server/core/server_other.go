//go:build !windows
// +build !windows

package core

import (
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 100 * time.Second
	s.WriteTimeout = 15 * time.Minute
	s.MaxHeaderBytes = 1 << 20
	return s
}
