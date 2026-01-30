// Package model 数据模型定义
package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型，所有数据模型都应嵌入此结构体
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BaseModelWithoutSoftDelete 不带软删除的基础模型
type BaseModelWithoutSoftDelete struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Status 通用状态类型
type Status int8

const (
	StatusDisabled Status = 0
	StatusEnabled  Status = 1
)

// String 返回状态的字符串表示
func (s Status) String() string {
	switch s {
	case StatusDisabled:
		return "disabled"
	case StatusEnabled:
		return "enabled"
	default:
		return "unknown"
	}
}

// IsEnabled 判断是否为启用状态
func (s Status) IsEnabled() bool {
	return s == StatusEnabled
}

// PageQuery 分页查询参数
type PageQuery struct {
	Page     int `form:"page" json:"page" binding:"min=1"`
	PageSize int `form:"page_size" json:"page_size" binding:"min=1,max=100"`
}

// GetOffset 计算偏移量
func (p *PageQuery) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取限制数量
func (p *PageQuery) GetLimit() int {
	return p.PageSize
}

// SetDefaults 设置默认值
func (p *PageQuery) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// SortQuery 排序查询参数
type SortQuery struct {
	SortBy    string `form:"sort_by" json:"sort_by"`
	SortOrder string `form:"sort_order" json:"sort_order"`
}

// GetOrderClause 获取排序子句
func (s *SortQuery) GetOrderClause(defaultField string) string {
	field := s.SortBy
	if field == "" {
		field = defaultField
	}

	order := s.SortOrder
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return field + " " + order
}
