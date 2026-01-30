# internal/handler 模块 AI 代码生成规则

> **模块职责**: HTTP处理层，接收请求、调用Service、返回响应

---

## 一、本模块的文件结构

```
internal/handler/
├── xxx_handler.go      # 业务Handler（按需创建）
└── RULE.md            # 本规则文件
```

**规则**: 每个业务实体对应一个 Handler 文件

---

## 二、创建新Handler的完整流程

### 步骤1: 创建文件

文件命名: `{实体名小写}_handler.go`，如 `order_handler.go`

### 步骤2: 定义结构体和构造函数

**推荐方式（Wire依赖注入）：**

```go
package handler

import (
    "strconv"

    "github.com/gin-gonic/gin"
    "ruleback/internal/model"
    "ruleback/internal/service"
    "ruleback/pkg/errors"
    "ruleback/pkg/logger"
    "ruleback/pkg/response"
)

// OrderHandler 订单HTTP处理器
type OrderHandler struct {
    service *service.OrderService
}

// NewOrderHandler 创建OrderHandler实例（用于Wire依赖注入）
func NewOrderHandler(svc *service.OrderService) *OrderHandler {
    return &OrderHandler{service: svc}
}
```

**兼容方式（单例模式，保留向后兼容）：**

```go
package handler

import (
    "strconv"
    "sync"

    "github.com/gin-gonic/gin"
    "ruleback/internal/model"
    "ruleback/internal/service"
    "ruleback/pkg/errors"
    "ruleback/pkg/logger"
    "ruleback/pkg/response"
)

var (
    orderHandlerInstance *OrderHandler
    orderHandlerOnce     sync.Once
)

// OrderHandler 订单HTTP处理器
type OrderHandler struct {
    service *service.OrderService
}

// NewOrderHandler 创建OrderHandler实例（用于Wire依赖注入）
func NewOrderHandler(svc *service.OrderService) *OrderHandler {
    return &OrderHandler{service: svc}
}

// GetOrderHandler 获取OrderHandler单例（兼容旧代码）
func GetOrderHandler() *OrderHandler {
    orderHandlerOnce.Do(func() {
        orderHandlerInstance = &OrderHandler{
            service: service.GetOrderService(),
        }
    })
    return orderHandlerInstance
}
```

### 步骤3: 实现Create接口

```go
// Create 创建订单
func (h *OrderHandler) Create(c *gin.Context) {
    var req model.CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.Warn("创建订单参数错误", logger.Err(err), logger.String("ip", c.ClientIP()))
        response.Fail(c, errors.CodeInvalidParams, "参数错误: "+err.Error())
        return
    }

    userID := c.GetUint("user_id")

    order, err := h.service.Create(userID, &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithData(c, order)
}
```

### 步骤4: 实现GetByID接口

```go
// GetByID 获取订单详情
func (h *OrderHandler) GetByID(c *gin.Context) {
    id, err := h.parseIDParam(c)
    if err != nil {
        response.Fail(c, errors.CodeInvalidParams, "无效的订单ID")
        return
    }

    order, err := h.service.GetByID(id)
    if err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithData(c, order)
}
```

### 步骤5: 实现List接口

```go
// List 获取订单列表
func (h *OrderHandler) List(c *gin.Context) {
    var query model.OrderListQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        logger.Warn("获取订单列表参数错误", logger.Err(err))
        response.Fail(c, errors.CodeInvalidParams, "参数错误: "+err.Error())
        return
    }

    orders, total, err := h.service.List(&query)
    if err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithPage(c, orders, total, query.Page, query.PageSize)
}
```

### 步骤6: 实现Update接口

```go
// Update 更新订单
func (h *OrderHandler) Update(c *gin.Context) {
    id, err := h.parseIDParam(c)
    if err != nil {
        response.Fail(c, errors.CodeInvalidParams, "无效的订单ID")
        return
    }

    var req model.UpdateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.Warn("更新订单参数错误", logger.Err(err), logger.Uint("order_id", id))
        response.Fail(c, errors.CodeInvalidParams, "参数错误: "+err.Error())
        return
    }

    order, err := h.service.Update(id, &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithData(c, order)
}
```

### 步骤7: 实现Delete接口

```go
// Delete 删除订单
func (h *OrderHandler) Delete(c *gin.Context) {
    id, err := h.parseIDParam(c)
    if err != nil {
        response.Fail(c, errors.CodeInvalidParams, "无效的订单ID")
        return
    }

    if err := h.service.Delete(id); err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithMessage(c, "删除成功")
}
```

### 步骤8: 实现辅助方法

```go
// parseIDParam 解析路径中的ID参数
func (h *OrderHandler) parseIDParam(c *gin.Context) (uint, error) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        return 0, err
    }
    return uint(id), nil
}

// handleError 统一处理错误响应
func (h *OrderHandler) handleError(c *gin.Context, err error) {
    appErr := errors.GetAppError(err)
    if appErr != nil {
        response.Fail(c, appErr.Code, appErr.Message)
        return
    }
    logger.Error("处理请求时发生未知错误", logger.Err(err))
    response.Fail(c, errors.CodeInternalError, "服务器内部错误")
}
```

---

## 三、响应函数使用规范

| 场景 | 使用函数 |
|------|---------|
| 创建成功 | `response.SuccessWithData(c, data)` |
| 查询详情 | `response.SuccessWithData(c, data)` |
| 列表查询 | `response.SuccessWithPage(c, list, total, page, pageSize)` |
| 更新成功 | `response.SuccessWithData(c, data)` |
| 删除成功 | `response.SuccessWithMessage(c, "删除成功")` |
| 参数错误 | `response.Fail(c, errors.CodeInvalidParams, msg)` |

---

## 四、禁止行为

| 禁止 | 正确做法 |
|------|---------|
| 在Handler中编写业务逻辑 | 业务逻辑放在Service层 |
| 直接调用Repository | 通过Service调用 |
| 使用c.JSON返回自定义格式 | 使用response包 |
| 使用装饰性分隔线注释 | 使用简洁单行注释 |

**注意**: 推荐使用 `New*` 构造函数配合Wire依赖注入，`Get*` 单例方法保留用于向后兼容

---

## 五、参数获取方式

| 来源 | 方法 | 示例 |
|------|------|------|
| 路径参数 | `c.Param("name")` | `/orders/:id` |
| 查询参数 | `c.Query("name")` | `?page=1` |
| 请求体JSON | `c.ShouldBindJSON(&req)` | POST/PUT body |
| Header | `c.GetHeader("name")` | Authorization |
| 上下文值 | `c.GetUint("name")` | user_id |

---

## 六、完整文件模板

```go
package handler

import (
    "strconv"
    "sync"

    "github.com/gin-gonic/gin"
    "ruleback/internal/model"
    "ruleback/internal/service"
    "ruleback/pkg/errors"
    "ruleback/pkg/logger"
    "ruleback/pkg/response"
)

var (
    xxxHandlerInstance *XxxHandler
    xxxHandlerOnce     sync.Once
)

// XxxHandler Xxx HTTP处理器
type XxxHandler struct {
    service *service.XxxService
}

// NewXxxHandler 创建XxxHandler实例（用于Wire依赖注入）
func NewXxxHandler(svc *service.XxxService) *XxxHandler {
    return &XxxHandler{service: svc}
}

// GetXxxHandler 获取XxxHandler单例（兼容旧代码）
func GetXxxHandler() *XxxHandler {
    xxxHandlerOnce.Do(func() {
        xxxHandlerInstance = &XxxHandler{
            service: service.GetXxxService(),
        }
    })
    return xxxHandlerInstance
}

// Create 创建记录
func (h *XxxHandler) Create(c *gin.Context) {
    var req model.CreateXxxRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Fail(c, errors.CodeInvalidParams, "参数错误: "+err.Error())
        return
    }

    xxx, err := h.service.Create(&req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithData(c, xxx)
}

// GetByID 获取详情
func (h *XxxHandler) GetByID(c *gin.Context) {
    id, err := h.parseIDParam(c)
    if err != nil {
        response.Fail(c, errors.CodeInvalidParams, "无效的ID")
        return
    }

    xxx, err := h.service.GetByID(id)
    if err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithData(c, xxx)
}

// List 获取列表
func (h *XxxHandler) List(c *gin.Context) {
    var query model.XxxListQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        response.Fail(c, errors.CodeInvalidParams, "参数错误: "+err.Error())
        return
    }

    items, total, err := h.service.List(&query)
    if err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithPage(c, items, total, query.Page, query.PageSize)
}

// Update 更新记录
func (h *XxxHandler) Update(c *gin.Context) {
    id, err := h.parseIDParam(c)
    if err != nil {
        response.Fail(c, errors.CodeInvalidParams, "无效的ID")
        return
    }

    var req model.UpdateXxxRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Fail(c, errors.CodeInvalidParams, "参数错误: "+err.Error())
        return
    }

    xxx, err := h.service.Update(id, &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithData(c, xxx)
}

// Delete 删除记录
func (h *XxxHandler) Delete(c *gin.Context) {
    id, err := h.parseIDParam(c)
    if err != nil {
        response.Fail(c, errors.CodeInvalidParams, "无效的ID")
        return
    }

    if err := h.service.Delete(id); err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithMessage(c, "删除成功")
}

// parseIDParam 解析路径中的ID参数
func (h *XxxHandler) parseIDParam(c *gin.Context) (uint, error) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    return uint(id), err
}

// handleError 统一处理错误响应
func (h *XxxHandler) handleError(c *gin.Context, err error) {
    appErr := errors.GetAppError(err)
    if appErr != nil {
        response.Fail(c, appErr.Code, appErr.Message)
        return
    }
    logger.Error("处理请求时发生未知错误", logger.Err(err))
    response.Fail(c, errors.CodeInternalError, "服务器内部错误")
}
```

---

## 七、创建Handler后的后续步骤

**推荐方式（Wire依赖注入）：**

1. 在 `internal/wire/providers.go` 中添加Provider:

```go
func ProvideOrderRepository(base *repository.BaseRepository) *repository.OrderRepository {
    return repository.NewOrderRepository(base)
}

func ProvideOrderService(repo *repository.OrderRepository) *service.OrderService {
    return service.NewOrderService(repo)
}

func ProvideOrderHandler(svc *service.OrderService) *handler.OrderHandler {
    return handler.NewOrderHandler(svc)
}
```

2. 在 `internal/wire/providers.go` 的 Handlers 结构体中添加字段:

```go
type Handlers struct {
    UserHandler  *handler.UserHandler
    OrderHandler *handler.OrderHandler  // 新增
}
```

3. 在路由中注册（使用Wire注入的Handler）:

```go
func registerOrderRoutes(rg *gin.RouterGroup, orderHandler *handler.OrderHandler) {
    orders := rg.Group("/orders")
    {
        orders.GET("", orderHandler.List)
        orders.POST("", orderHandler.Create)
        orders.GET("/:id", orderHandler.GetByID)
        orders.PUT("/:id", orderHandler.Update)
        orders.DELETE("/:id", orderHandler.Delete)
    }
}
```

4. 重新生成Wire代码: `~/go/bin/wire ./internal/wire/...`

**兼容方式（单例模式）：**

```go
func registerOrderRoutes(rg *gin.RouterGroup) {
    h := handler.GetOrderHandler()  // 使用Get单例方法

    orders := rg.Group("/orders")
    {
        orders.GET("", h.List)
        orders.POST("", h.Create)
        orders.GET("/:id", h.GetByID)
        orders.PUT("/:id", h.Update)
        orders.DELETE("/:id", h.Delete)
    }
}
```
