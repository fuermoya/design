//go:build windows
// +build windows

package core

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func initServer(address string, router *gin.Engine) server {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    50 * time.Second,
		WriteTimeout:   15 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
}
