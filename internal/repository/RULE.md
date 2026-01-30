# internal/repository 模块 AI 代码生成规则

> **模块职责**: 数据访问层，封装所有数据库操作

---

## 一、本模块的文件结构

```
internal/repository/
├── base.go              # 基础Repository（勿修改）
├── interfaces.go        # Repository接口定义（按需创建）
├── xxx_repository.go    # 业务Repository（按需创建）
└── RULE.md             # 本规则文件
```

**规则**: 每个数据模型对应一个 Repository 文件

---

## 二、创建新Repository的完整流程

### 步骤1: 创建文件

文件命名: `{实体名小写}_repository.go`，如 `order_repository.go`

### 步骤2: 定义结构体和构造函数

**推荐方式（Wire依赖注入）：**

```go
package repository

import (
    "gorm.io/gorm"
    "ruleback/internal/model"
)

// OrderRepository 订单数据访问层
type OrderRepository struct {
    *BaseRepository
}

// NewOrderRepository 创建OrderRepository实例（用于Wire依赖注入）
func NewOrderRepository(base *BaseRepository) *OrderRepository {
    return &OrderRepository{BaseRepository: base}
}
```

**兼容方式（单例模式，保留向后兼容）：**

```go
package repository

import (
    "sync"

    "gorm.io/gorm"
    "ruleback/internal/model"
)

var (
    orderRepoInstance *OrderRepository
    orderRepoOnce     sync.Once
)

// OrderRepository 订单数据访问层
type OrderRepository struct {
    *BaseRepository
}

// NewOrderRepository 创建OrderRepository实例（用于Wire依赖注入）
func NewOrderRepository(base *BaseRepository) *OrderRepository {
    return &OrderRepository{BaseRepository: base}
}

// GetOrderRepository 获取OrderRepository单例（兼容旧代码）
func GetOrderRepository() *OrderRepository {
    orderRepoOnce.Do(func() {
        orderRepoInstance = &OrderRepository{
            BaseRepository: GetBaseRepository(),
        }
    })
    return orderRepoInstance
}
```

### 步骤3: 实现基础CRUD方法

```go
// Create 创建订单
func (r *OrderRepository) Create(order *model.Order) error {
    return r.DB().Create(order).Error
}

// Update 更新订单
func (r *OrderRepository) Update(order *model.Order) error {
    return r.DB().Save(order).Error
}

// Delete 删除订单（软删除）
func (r *OrderRepository) Delete(id uint) error {
    return r.DB().Delete(&model.Order{}, id).Error
}

// GetByID 根据ID获取订单
func (r *OrderRepository) GetByID(id uint) (*model.Order, error) {
    var order model.Order
    err := r.DB().First(&order, id).Error
    if err != nil {
        return nil, err
    }
    return &order, nil
}
```

### 步骤4: 实现列表查询方法

```go
// List 获取订单列表
func (r *OrderRepository) List(query *model.OrderListQuery) ([]*model.Order, int64, error) {
    var items []*model.Order
    var total int64

    db := r.DB().Model(&model.Order{})
    db = r.applyFilters(db, query)

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    db = db.Scopes(r.OrderBy(query.SortBy, query.SortOrder))
    db = db.Scopes(r.Paginate(query.Page, query.PageSize))

    if err := db.Find(&items).Error; err != nil {
        return nil, 0, err
    }

    return items, total, nil
}

// applyFilters 应用查询过滤条件
func (r *OrderRepository) applyFilters(db *gorm.DB, query *model.OrderListQuery) *gorm.DB {
    if query.UserID != nil {
        db = db.Where("user_id = ?", *query.UserID)
    }
    if query.Status != nil {
        db = db.Where("status = ?", *query.Status)
    }
    return db
}
```

### 步骤5: 实现更新方法

```go
// UpdateFields 更新指定字段
func (r *OrderRepository) UpdateFields(id uint, fields map[string]interface{}) error {
    return r.DB().Model(&model.Order{}).Where("id = ?", id).Updates(fields).Error
}

// UpdateStatus 更新订单状态
func (r *OrderRepository) UpdateStatus(id uint, status model.OrderStatus) error {
    return r.UpdateFields(id, map[string]interface{}{"status": status})
}
```

---

## 三、方法命名规范

| 操作类型 | 命名规范 | 示例 |
|---------|---------|------|
| 创建 | `Create` | `Create(model)` |
| 更新 | `Update` | `Update(model)` |
| 删除 | `Delete` | `Delete(id)` |
| 按ID查询 | `GetByID` | `GetByID(id)` |
| 按字段查询 | `GetBy` + 字段名 | `GetByOrderNo(orderNo)` |
| 检查存在 | `ExistsBy` + 字段名 | `ExistsByOrderNo(orderNo)` |
| 列表查询 | `List` | `List(query)` |
| 更新字段 | `UpdateFields` | `UpdateFields(id, fields)` |

---

## 四、返回值规范

| 操作类型 | 返回值 |
|---------|-------|
| 单条查询 | `(*Model, error)` |
| 列表查询 | `([]*Model, int64, error)` |
| 存在检查 | `(bool, error)` |
| 写操作 | `error` |

---

## 五、禁止行为

| 禁止 | 正确做法 |
|------|---------|
| 在Repository中包含业务逻辑 | 业务逻辑放在Service层 |
| 使用硬编码SQL | 使用GORM方法 |
| 调用其他Repository | 在Service层协调 |
| 使用装饰性分隔线注释 | 使用简洁单行注释 |

**注意**: 推荐使用 `New*` 构造函数配合Wire依赖注入，`Get*` 单例方法保留用于向后兼容

---

## 六、已存在的基础方法（BaseRepository）

```go
r.DB()                          // 获取数据库实例
r.Paginate(page, pageSize)      // 分页Scope
r.OrderBy(field, order)         // 排序Scope
r.Transaction(fn)               // 事务支持
```

---

## 七、完整文件模板

```go
package repository

import (
    "sync"

    "gorm.io/gorm"
    "ruleback/internal/model"
)

var (
    xxxRepoInstance *XxxRepository
    xxxRepoOnce     sync.Once
)

// XxxRepository Xxx数据访问层
type XxxRepository struct {
    *BaseRepository
}

// NewXxxRepository 创建XxxRepository实例（用于Wire依赖注入）
func NewXxxRepository(base *BaseRepository) *XxxRepository {
    return &XxxRepository{BaseRepository: base}
}

// GetXxxRepository 获取XxxRepository单例（兼容旧代码）
func GetXxxRepository() *XxxRepository {
    xxxRepoOnce.Do(func() {
        xxxRepoInstance = &XxxRepository{
            BaseRepository: GetBaseRepository(),
        }
    })
    return xxxRepoInstance
}

// Create 创建记录
func (r *XxxRepository) Create(xxx *model.Xxx) error {
    return r.DB().Create(xxx).Error
}

// Update 更新记录
func (r *XxxRepository) Update(xxx *model.Xxx) error {
    return r.DB().Save(xxx).Error
}

// Delete 删除记录
func (r *XxxRepository) Delete(id uint) error {
    return r.DB().Delete(&model.Xxx{}, id).Error
}

// GetByID 根据ID获取记录
func (r *XxxRepository) GetByID(id uint) (*model.Xxx, error) {
    var xxx model.Xxx
    if err := r.DB().First(&xxx, id).Error; err != nil {
        return nil, err
    }
    return &xxx, nil
}

// List 获取列表
func (r *XxxRepository) List(query *model.XxxListQuery) ([]*model.Xxx, int64, error) {
    var items []*model.Xxx
    var total int64

    db := r.DB().Model(&model.Xxx{})
    db = r.applyFilters(db, query)

    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    db = db.Scopes(r.OrderBy(query.SortBy, query.SortOrder))
    db = db.Scopes(r.Paginate(query.Page, query.PageSize))

    if err := db.Find(&items).Error; err != nil {
        return nil, 0, err
    }

    return items, total, nil
}

func (r *XxxRepository) applyFilters(db *gorm.DB, query *model.XxxListQuery) *gorm.DB {
    // 添加筛选条件
    return db
}

// UpdateFields 更新指定字段
func (r *XxxRepository) UpdateFields(id uint, fields map[string]interface{}) error {
    return r.DB().Model(&model.Xxx{}).Where("id = ?", id).Updates(fields).Error
}
```
