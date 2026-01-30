// Package main 应用程序入口
package main

import (
	"fmt"
	"os"
	"time"

	"ruleback/internal/config"
)

const (
	ConfigPath      = "configs/config.yaml"
	ShutdownTimeout = 5 * time.Second
)

var cfg *config.Config

func main() {
	// 1. 初始化应用
	if err := initApp(); err != nil {
		fmt.Printf("应用初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 2. 启动服务器
	srv := startServer()

	// 3. 等待关闭信号
	gracefulShutdown(srv)
}
