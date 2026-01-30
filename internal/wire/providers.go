// Package wire 依赖注入配置
package wire

import (
	"gorm.io/gorm"
	"ruleback/internal/repository"
)

// ProvideBaseRepository 提供BaseRepository实例
func ProvideBaseRepository(db *gorm.DB) *repository.BaseRepository {
	return repository.NewBaseRepository(db)
}

// Handlers 包含所有Handler实例
// 使用框架时，请在此结构体中添加你的Handler
// 示例:
//
//	type Handlers struct {
//	    UserHandler    *handler.UserHandler
//	    ProductHandler *handler.ProductHandler
//	}
type Handlers struct {
	// 在此添加你的Handler字段
}

// ProvideHandlers 提供所有Handler实例
// 使用框架时，请修改此函数以注入你的Handler
// 示例:
//
//	func ProvideHandlers(userHandler *handler.UserHandler) *Handlers {
//	    return &Handlers{UserHandler: userHandler}
//	}
func ProvideHandlers() *Handlers {
	return &Handlers{}
}
