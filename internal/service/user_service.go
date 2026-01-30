// Package service 业务逻辑层
package service

import (
	"errors"
	"sync"

	"gorm.io/gorm"
	apperrors "ruleback/pkg/errors"
	"ruleback/pkg/logger"

	"ruleback/internal/model"
	"ruleback/internal/repository"
)

var (
	userServiceInstance *UserService
	userServiceOnce     sync.Once
)

// UserService 用户业务逻辑层
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService 创建UserService实例（用于依赖注入）
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetUserService 获取用户Service单例（保留向后兼容）
func GetUserService() *UserService {
	userServiceOnce.Do(func() {
		userServiceInstance = &UserService{
			repo: repository.GetUserRepository(),
		}
	})
	return userServiceInstance
}

// Create 创建用户
func (s *UserService) Create(req *model.CreateUserRequest) (*model.User, error) {
	logger.Debug("开始创建用户", logger.String("username", req.Username))

	exists, err := s.repo.ExistsByUsername(req.Username)
	if err != nil {
		logger.Error("检查用户名失败", logger.Err(err), logger.String("username", req.Username))
		return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "检查用户名失败", err)
	}
	if exists {
		return nil, apperrors.New(apperrors.CodeUserExists, "用户名已存在")
	}

	exists, err = s.repo.ExistsByEmail(req.Email)
	if err != nil {
		logger.Error("检查邮箱失败", logger.Err(err), logger.String("email", req.Email))
		return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "检查邮箱失败", err)
	}
	if exists {
		return nil, apperrors.New(apperrors.CodeUserExists, "邮箱已存在")
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: s.hashPassword(req.Password),
		Nickname: req.Nickname,
		Status:   model.StatusEnabled,
	}

	if err := s.repo.Create(user); err != nil {
		logger.Error("创建用户失败", logger.Err(err), logger.String("username", req.Username))
		return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "创建用户失败", err)
	}

	logger.Info("用户创建成功", logger.Uint("user_id", user.ID), logger.String("username", user.Username))
	return user, nil
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(id uint) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		logger.Error("获取用户失败", logger.Err(err), logger.Uint("user_id", id))
		return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "获取用户失败", err)
	}
	return user, nil
}

// GetByUsername 根据用户名获取用户
func (s *UserService) GetByUsername(username string) (*model.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		logger.Error("获取用户失败", logger.Err(err), logger.String("username", username))
		return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "获取用户失败", err)
	}
	return user, nil
}

// List 获取用户列表
func (s *UserService) List(query *model.UserListQuery) ([]*model.User, int64, error) {
	query.PageQuery.SetDefaults()

	users, total, err := s.repo.List(query)
	if err != nil {
		logger.Error("获取用户列表失败", logger.Err(err))
		return nil, 0, apperrors.Wrap(apperrors.CodeDatabaseError, "获取用户列表失败", err)
	}

	return users, total, nil
}

// Update 更新用户信息
func (s *UserService) Update(id uint, req *model.UpdateUserRequest) (*model.User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		return user, nil
	}

	if err := s.repo.UpdateFields(id, updates); err != nil {
		logger.Error("更新用户失败", logger.Err(err), logger.Uint("user_id", id))
		return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "更新用户失败", err)
	}

	return s.GetByID(id)
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		logger.Error("删除用户失败", logger.Err(err), logger.Uint("user_id", id))
		return apperrors.Wrap(apperrors.CodeDatabaseError, "删除用户失败", err)
	}

	logger.Info("用户删除成功", logger.Uint("user_id", id))
	return nil
}

// hashPassword 对密码进行加密
// TODO: 使用 bcrypt 实现真正的密码加密
func (s *UserService) hashPassword(password string) string {
	return password
}

// verifyPassword 验证密码
// TODO: 使用 bcrypt 实现真正的密码验证
func (s *UserService) verifyPassword(hashedPassword, password string) bool {
	return hashedPassword == password
}
