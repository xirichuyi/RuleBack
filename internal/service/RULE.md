# internal/service 模块 AI 代码生成规则

> **模块职责**: 业务逻辑层，处理所有业务规则和数据转换

---

## 一、本模块的文件结构

```
internal/service/
├── interfaces.go       # Service接口定义（按需创建）
├── xxx_service.go      # 业务Service（按需创建）
└── RULE.md            # 本规则文件
```

**规则**: 每个业务实体对应一个 Service 文件

---

## 二、创建新Service的完整流程

### 步骤1: 创建文件

文件命名: `{实体名小写}_service.go`，如 `order_service.go`

### 步骤2: 定义结构体和构造函数

**推荐方式（Wire依赖注入）：**

```go
package service

import (
    "errors"

    "gorm.io/gorm"
    apperrors "ruleback/pkg/errors"
    "ruleback/pkg/logger"
    "ruleback/internal/model"
    "ruleback/internal/repository"
)

// OrderService 订单业务逻辑层
type OrderService struct {
    repo *repository.OrderRepository
}

// NewOrderService 创建OrderService实例（用于Wire依赖注入）
func NewOrderService(repo *repository.OrderRepository) *OrderService {
    return &OrderService{repo: repo}
}
```

**兼容方式（单例模式，保留向后兼容）：**

```go
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
    orderServiceInstance *OrderService
    orderServiceOnce     sync.Once
)

// OrderService 订单业务逻辑层
type OrderService struct {
    repo *repository.OrderRepository
}

// NewOrderService 创建OrderService实例（用于Wire依赖注入）
func NewOrderService(repo *repository.OrderRepository) *OrderService {
    return &OrderService{repo: repo}
}

// GetOrderService 获取OrderService单例（兼容旧代码）
func GetOrderService() *OrderService {
    orderServiceOnce.Do(func() {
        orderServiceInstance = &OrderService{
            repo: repository.GetOrderRepository(),
        }
    })
    return orderServiceInstance
}
```

### 步骤3: 实现Create方法

```go
// Create 创建订单
func (s *OrderService) Create(userID uint, req *model.CreateOrderRequest) (*model.Order, error) {
    logger.Debug("开始创建订单", logger.Uint("user_id", userID))

    order := &model.Order{
        UserID:  userID,
        OrderNo: s.generateOrderNo(),
        Amount:  req.Amount,
        Status:  model.OrderStatusPending,
    }

    if err := s.repo.Create(order); err != nil {
        logger.Error("创建订单失败", logger.Err(err))
        return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "创建订单失败", err)
    }

    logger.Info("订单创建成功", logger.Uint("order_id", order.ID))
    return order, nil
}
```

### 步骤4: 实现GetByID方法

```go
// GetByID 根据ID获取订单
func (s *OrderService) GetByID(id uint) (*model.Order, error) {
    order, err := s.repo.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, apperrors.ErrNotFound
        }
        logger.Error("获取订单失败", logger.Err(err), logger.Uint("order_id", id))
        return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "获取订单失败", err)
    }
    return order, nil
}
```

### 步骤5: 实现List方法

```go
// List 获取订单列表
func (s *OrderService) List(query *model.OrderListQuery) ([]*model.Order, int64, error) {
    query.PageQuery.SetDefaults()

    orders, total, err := s.repo.List(query)
    if err != nil {
        logger.Error("获取订单列表失败", logger.Err(err))
        return nil, 0, apperrors.Wrap(apperrors.CodeDatabaseError, "获取订单列表失败", err)
    }

    return orders, total, nil
}
```

### 步骤6: 实现Update方法

```go
// Update 更新订单
func (s *OrderService) Update(id uint, req *model.UpdateOrderRequest) (*model.Order, error) {
    if _, err := s.GetByID(id); err != nil {
        return nil, err
    }

    updates := make(map[string]interface{})
    if req.Status != nil {
        updates["status"] = *req.Status
    }
    if req.Remark != nil {
        updates["remark"] = *req.Remark
    }

    if len(updates) == 0 {
        return s.GetByID(id)
    }

    if err := s.repo.UpdateFields(id, updates); err != nil {
        logger.Error("更新订单失败", logger.Err(err), logger.Uint("order_id", id))
        return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "更新订单失败", err)
    }

    return s.GetByID(id)
}
```

### 步骤7: 实现Delete方法

```go
// Delete 删除订单
func (s *OrderService) Delete(id uint) error {
    if _, err := s.GetByID(id); err != nil {
        return err
    }

    if err := s.repo.Delete(id); err != nil {
        logger.Error("删除订单失败", logger.Err(err), logger.Uint("order_id", id))
        return apperrors.Wrap(apperrors.CodeDatabaseError, "删除订单失败", err)
    }

    logger.Info("订单删除成功", logger.Uint("order_id", id))
    return nil
}
```

---

## 三、错误处理规范

```go
// 转换gorm错误为业务错误
if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, apperrors.ErrNotFound
}

// 包装数据库错误
return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "操作失败", err)

// 返回业务错误
return nil, apperrors.New(apperrors.CodeConflict, "订单号已存在")
```

---

## 四、日志记录规范

| 场景 | 日志级别 | 示例 |
|------|---------|------|
| 方法开始 | Debug | `logger.Debug("开始创建订单", ...)` |
| 业务校验失败 | Warn | `logger.Warn("订单号已存在", ...)` |
| 数据库错误 | Error | `logger.Error("创建订单失败", logger.Err(err), ...)` |
| 操作成功 | Info | `logger.Info("订单创建成功", ...)` |

---

## 五、禁止行为

| 禁止 | 正确做法 |
|------|---------|
| 直接操作数据库 | 通过Repository操作 |
| 接收gin.Context参数 | 接收业务参数 |
| 直接返回Repository错误 | 转换为业务错误 |
| 使用fmt.Println | 使用logger包 |
| 使用装饰性分隔线注释 | 使用简洁单行注释 |

**注意**: 推荐使用 `New*` 构造函数配合Wire依赖注入，`Get*` 单例方法保留用于向后兼容

---

## 六、方法命名规范

| 操作类型 | 命名规范 | 示例 |
|---------|---------|------|
| 创建 | `Create` | `Create(req)` |
| 更新 | `Update` | `Update(id, req)` |
| 删除 | `Delete` | `Delete(id)` |
| 按ID查询 | `GetByID` | `GetByID(id)` |
| 列表查询 | `List` | `List(query)` |
| 业务操作 | 动词 + 名词 | `CancelOrder`, `PayOrder` |

---

## 七、完整文件模板

```go
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
    xxxServiceInstance *XxxService
    xxxServiceOnce     sync.Once
)

// XxxService Xxx业务逻辑层
type XxxService struct {
    repo *repository.XxxRepository
}

// NewXxxService 创建XxxService实例（用于Wire依赖注入）
func NewXxxService(repo *repository.XxxRepository) *XxxService {
    return &XxxService{repo: repo}
}

// GetXxxService 获取XxxService单例（兼容旧代码）
func GetXxxService() *XxxService {
    xxxServiceOnce.Do(func() {
        xxxServiceInstance = &XxxService{
            repo: repository.GetXxxRepository(),
        }
    })
    return xxxServiceInstance
}

// Create 创建记录
func (s *XxxService) Create(req *model.CreateXxxRequest) (*model.Xxx, error) {
    logger.Debug("开始创建Xxx", logger.String("field", req.Field))

    xxx := &model.Xxx{
        Field:  req.Field,
        Status: model.StatusEnabled,
    }

    if err := s.repo.Create(xxx); err != nil {
        logger.Error("创建Xxx失败", logger.Err(err))
        return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "创建失败", err)
    }

    logger.Info("Xxx创建成功", logger.Uint("id", xxx.ID))
    return xxx, nil
}

// GetByID 根据ID获取记录
func (s *XxxService) GetByID(id uint) (*model.Xxx, error) {
    xxx, err := s.repo.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, apperrors.ErrNotFound
        }
        return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "获取失败", err)
    }
    return xxx, nil
}

// List 获取列表
func (s *XxxService) List(query *model.XxxListQuery) ([]*model.Xxx, int64, error) {
    query.PageQuery.SetDefaults()
    items, total, err := s.repo.List(query)
    if err != nil {
        return nil, 0, apperrors.Wrap(apperrors.CodeDatabaseError, "获取列表失败", err)
    }
    return items, total, nil
}

// Update 更新记录
func (s *XxxService) Update(id uint, req *model.UpdateXxxRequest) (*model.Xxx, error) {
    if _, err := s.GetByID(id); err != nil {
        return nil, err
    }

    updates := make(map[string]interface{})
    if req.Field != nil {
        updates["field"] = *req.Field
    }

    if len(updates) == 0 {
        return s.GetByID(id)
    }

    if err := s.repo.UpdateFields(id, updates); err != nil {
        return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "更新失败", err)
    }

    return s.GetByID(id)
}

// Delete 删除记录
func (s *XxxService) Delete(id uint) error {
    if _, err := s.GetByID(id); err != nil {
        return err
    }

    if err := s.repo.Delete(id); err != nil {
        return apperrors.Wrap(apperrors.CodeDatabaseError, "删除失败", err)
    }

    logger.Info("Xxx删除成功", logger.Uint("id", id))
    return nil
}
```
