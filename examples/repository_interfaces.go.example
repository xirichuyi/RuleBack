// Package repository 数据访问层接口定义
package repository

import (
	"ruleback/internal/model"
)

// UserRepositoryInterface 用户数据访问层接口
type UserRepositoryInterface interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	List(query *model.UserListQuery) ([]*model.User, int64, error)
	UpdateFields(id uint, fields map[string]interface{}) error
	UpdateStatus(id uint, status model.Status) error
	UpdatePassword(id uint, password string) error
}
