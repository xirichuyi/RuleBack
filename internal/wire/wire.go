//go:build wireinject
// +build wireinject

// Package wire 依赖注入配置
package wire

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet 所有Provider的集合
// 使用框架时，请在此处添加你的Provider
// 示例:
//
//	var ProviderSet = wire.NewSet(
//	    ProvideBaseRepository,
//	    ProvideUserRepository,
//	    ProvideUserService,
//	    ProvideUserHandler,
//	    ProvideHandlers,
//	)
var ProviderSet = wire.NewSet(
	ProvideBaseRepository,
	ProvideHandlers,
)

// InitializeHandlers 初始化所有Handler
// 使用框架时，Wire会根据ProviderSet自动生成依赖注入代码
func InitializeHandlers(db *gorm.DB) (*Handlers, error) {
	wire.Build(ProviderSet)
	return nil, nil
}
