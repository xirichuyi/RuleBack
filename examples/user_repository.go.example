package repository

import (
	"sync"

	"gorm.io/gorm"
	"ruleback/internal/model"
)

var (
	userRepoInstance *UserRepository
	userRepoOnce     sync.Once
)

// UserRepository 用户数据访问层
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository 创建UserRepository实例（用于依赖注入）
func NewUserRepository(base *BaseRepository) *UserRepository {
	return &UserRepository{BaseRepository: base}
}

// GetUserRepository 获取用户Repository单例（保留向后兼容）
func GetUserRepository() *UserRepository {
	userRepoOnce.Do(func() {
		userRepoInstance = &UserRepository{
			BaseRepository: GetBaseRepository(),
		}
	})
	return userRepoInstance
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.DB().Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.DB().Save(user).Error
}

// Delete 删除用户（软删除）
func (r *UserRepository) Delete(id uint) error {
	return r.DB().Delete(&model.User{}, id).Error
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.DB().First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB().Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB().Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.DB().Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.DB().Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// List 获取用户列表
func (r *UserRepository) List(query *model.UserListQuery) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	db := r.DB().Model(&model.User{})
	db = r.applyFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Scopes(r.OrderBy(query.SortQuery.SortBy, query.SortQuery.SortOrder))
	db = db.Scopes(r.Paginate(query.Page, query.PageSize))

	if err := db.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// applyFilters 应用查询过滤条件
func (r *UserRepository) applyFilters(db *gorm.DB, query *model.UserListQuery) *gorm.DB {
	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.Email != "" {
		db = db.Where("email LIKE ?", "%"+query.Email+"%")
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	return db
}

// UpdateFields 更新指定字段
func (r *UserRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.DB().Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
}

// UpdateStatus 更新用户状态
func (r *UserRepository) UpdateStatus(id uint, status model.Status) error {
	return r.UpdateFields(id, map[string]interface{}{"status": status})
}

// UpdatePassword 更新用户密码
func (r *UserRepository) UpdatePassword(id uint, password string) error {
	return r.UpdateFields(id, map[string]interface{}{"password": password})
}
