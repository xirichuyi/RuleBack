// Package main 应用初始化和服务器管理
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ruleback/internal/config"
	"ruleback/internal/router"
	"ruleback/internal/wire"
	"ruleback/pkg/database"
	"ruleback/pkg/logger"
)

// initApp 初始化应用程序
func initApp() error {
	var err error

	cfg, err = config.Load(ConfigPath)
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	if err = initLogger(); err != nil {
		return fmt.Errorf("初始化日志失败: %w", err)
	}

	if err = initDatabase(); err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}

	if err = migrateDatabase(); err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	return nil
}

// initLogger 初始化日志系统
func initLogger() error {
	if err := logger.Init(&cfg.Log); err != nil {
		return err
	}
	logger.Info("日志系统初始化完成", logger.String("level", cfg.Log.Level))
	return nil
}

// initDatabase 初始化数据库连接
func initDatabase() error {
	if err := database.Init(&cfg.Database); err != nil {
		return err
	}
	logger.Info("数据库连接初始化完成", logger.String("host", cfg.Database.Host))
	return nil
}

// migrateDatabase 执行数据库迁移
// 使用框架时，在此处添加你的模型进行自动迁移
// 示例:
//
//	models := []interface{}{
//	    &model.User{},
//	    &model.Product{},
//	}
func migrateDatabase() error {
	models := []interface{}{
		// 在此添加你的模型
	}

	if len(models) == 0 {
		logger.Info("没有模型需要迁移")
		return nil
	}

	if err := database.AutoMigrate(models...); err != nil {
		return err
	}

	logger.Info("数据库迁移完成", logger.Int("models", len(models)))
	return nil
}

// startServer 启动HTTP服务器
func startServer() *http.Server {
	// 使用Wire初始化所有Handler
	handlers, err := wire.InitializeHandlers(database.GetDB())
	if err != nil {
		logger.Fatal("初始化Handler失败", logger.Err(err))
	}

	r := router.Setup(handlers)

	srv := &http.Server{
		Addr:         cfg.Server.GetAddress(),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		logger.Info("HTTP服务器启动",
			logger.String("env", cfg.App.Env),
			logger.String("address", cfg.Server.GetAddress()),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP服务器异常退出", logger.Err(err))
		}
	}()

	return srv
}

// gracefulShutdown 优雅关闭服务器
func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Info("收到关闭信号", logger.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("HTTP服务器关闭异常", logger.Err(err))
	}

	if err := database.Close(); err != nil {
		logger.Error("数据库关闭异常", logger.Err(err))
	}

	logger.Sync()
	logger.Info("应用已安全关闭")
}
