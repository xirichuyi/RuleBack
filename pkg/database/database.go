// Package database 数据库连接管理
package database

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"ruleback/internal/config"
)

var (
	globalDB *gorm.DB
	dbOnce   sync.Once
	initErr  error
)

// Init 初始化数据库连接（使用sync.Once确保只初始化一次）
func Init(cfg *config.DatabaseConfig) error {
	dbOnce.Do(func() {
		var dialector gorm.Dialector

		switch cfg.Driver {
		case "mysql":
			dialector = mysql.Open(cfg.GetDSN())
		case "postgres":
			dialector = postgres.Open(cfg.GetDSN())
		default:
			initErr = fmt.Errorf("不支持的数据库驱动: %s", cfg.Driver)
			return
		}

		db, err := gorm.Open(dialector, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			initErr = fmt.Errorf("连接数据库失败: %w", err)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			initErr = fmt.Errorf("获取数据库连接失败: %w", err)
			return
		}

		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

		if err := sqlDB.Ping(); err != nil {
			initErr = fmt.Errorf("数据库连接测试失败: %w", err)
			return
		}

		globalDB = db
	})

	return initErr
}

// GetDB 获取全局数据库实例
func GetDB() *gorm.DB {
	return globalDB
}

// Close 关闭数据库连接
func Close() error {
	if globalDB == nil {
		return nil
	}

	sqlDB, err := globalDB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate(models ...interface{}) error {
	if globalDB == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return globalDB.AutoMigrate(models...)
}
