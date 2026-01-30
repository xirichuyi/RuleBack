//go:build wireinject
// +build wireinject

// Package wire 依赖注入配置
package wire

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet 所有Provider的集合
var ProviderSet = wire.NewSet(
	ProvideBaseRepository,
	ProvideUserRepository,
	ProvideUserService,
	ProvideUserHandler,
	ProvideHandlers,
)

// InitializeHandlers 初始化所有Handler
func InitializeHandlers(db *gorm.DB) (*Handlers, error) {
	wire.Build(ProviderSet)
	return nil, nil
}
