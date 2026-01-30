// Package wire 依赖注入配置
package wire

import (
	"gorm.io/gorm"
	"ruleback/internal/handler"
	"ruleback/internal/repository"
	"ruleback/internal/service"
)

// ProvideBaseRepository 提供BaseRepository实例
func ProvideBaseRepository(db *gorm.DB) *repository.BaseRepository {
	return repository.NewBaseRepository(db)
}

// ProvideUserRepository 提供UserRepository实例
func ProvideUserRepository(base *repository.BaseRepository) *repository.UserRepository {
	return repository.NewUserRepository(base)
}

// ProvideUserService 提供UserService实例
func ProvideUserService(repo *repository.UserRepository) *service.UserService {
	return service.NewUserService(repo)
}

// ProvideUserHandler 提供UserHandler实例
func ProvideUserHandler(svc *service.UserService) *handler.UserHandler {
	return handler.NewUserHandler(svc)
}

// Handlers 包含所有Handler实例
type Handlers struct {
	UserHandler *handler.UserHandler
}

// ProvideHandlers 提供所有Handler实例
func ProvideHandlers(userHandler *handler.UserHandler) *Handlers {
	return &Handlers{
		UserHandler: userHandler,
	}
}
