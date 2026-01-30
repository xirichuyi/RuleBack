// Package repository 数据访问层
package repository

import (
	"sync"

	"gorm.io/gorm"
	"ruleback/pkg/database"
)

var (
	baseRepoInstance *BaseRepository
	baseRepoOnce     sync.Once
)

// BaseRepository 基础Repository，提供通用的数据库操作方法
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository 创建BaseRepository实例（用于依赖注入）
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// GetBaseRepository 获取基础Repository单例（保留向后兼容）
func GetBaseRepository() *BaseRepository {
	baseRepoOnce.Do(func() {
		baseRepoInstance = &BaseRepository{
			db: database.GetDB(),
		}
	})
	return baseRepoInstance
}

// DB 获取数据库实例
func (r *BaseRepository) DB() *gorm.DB {
	return r.db
}

// Create 创建记录
func (r *BaseRepository) Create(model interface{}) error {
	return r.db.Create(model).Error
}

// Update 更新记录
func (r *BaseRepository) Update(model interface{}) error {
	return r.db.Save(model).Error
}

// Delete 删除记录（软删除）
func (r *BaseRepository) Delete(model interface{}) error {
	return r.db.Delete(model).Error
}

// DeleteByID 根据ID删除记录
func (r *BaseRepository) DeleteByID(model interface{}, id uint) error {
	return r.db.Delete(model, id).Error
}

// GetByID 根据ID获取记录
func (r *BaseRepository) GetByID(model interface{}, id uint) error {
	return r.db.First(model, id).Error
}

// Paginate 分页查询
func (r *BaseRepository) Paginate(page, pageSize int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// OrderBy 排序查询
func (r *BaseRepository) OrderBy(field, order string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" {
			return db
		}
		if order != "asc" && order != "desc" {
			order = "desc"
		}
		return db.Order(field + " " + order)
	}
}

// Transaction 执行事务
func (r *BaseRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
