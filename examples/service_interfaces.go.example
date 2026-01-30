// Package service 业务逻辑层接口定义
package service

import (
	"ruleback/internal/model"
)

// UserServiceInterface 用户业务逻辑层接口
type UserServiceInterface interface {
	Create(req *model.CreateUserRequest) (*model.User, error)
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	List(query *model.UserListQuery) ([]*model.User, int64, error)
	Update(id uint, req *model.UpdateUserRequest) (*model.User, error)
	Delete(id uint) error
}
