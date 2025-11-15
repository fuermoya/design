package internal

import (
	"github.com/fuermoya/design/server/global"
	"go.uber.org/zap/zapcore"
	"os"
)

var FileRotatelogs = new(fileRotatelogs)

type fileRotatelogs struct{}

// GetWriteSyncer 获取 zapcore.WriteSyncer
// Author [SliverHorn](https://github.com/SliverHorn)
func (r *fileRotatelogs) GetWriteSyncer(level string) zapcore.WriteSyncer {
	if global.CONFIG.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	} else {
		fileWriter := NewCutter(global.CONFIG.Zap.Director, level, WithCutterFormat("2006-01-02"))
		return zapcore.AddSync(fileWriter)
	}
}
