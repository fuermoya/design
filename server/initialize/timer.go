package initialize

import (
	"fmt"
	"github.com/fuermoya/design/server/global"
	"github.com/fuermoya/design/server/task"
	"github.com/robfig/cron/v3"
)

func Timer() {
	var option []cron.Option
	option = append(option, cron.WithSeconds())
	// 清理DB定时任务
	_, err := global.Timer.AddTaskByFunc("ClearDB", "@daily", func() {
		err := task.ClearTable(global.DB) // 定时任务方法定在task文件包中
		if err != nil {
			fmt.Println("timer error:", err)
		}
	}, "定时清理数据库【日志，黑名单】内容", option...)
	if err != nil {
		fmt.Println("add timer error:", err)
	}
}
