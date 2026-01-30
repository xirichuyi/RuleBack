# internal/wire 模块 AI 代码生成规则

> **模块职责**: 依赖注入配置，使用Google Wire管理组件依赖关系

---

## 一、本模块的文件结构

```
internal/wire/
├── providers.go   # Provider函数定义
├── wire.go        # Wire injector定义（编译时忽略）
├── wire_gen.go    # Wire生成的代码（勿手动修改）
└── RULE.md        # 本规则文件
```

---

## 二、Wire核心概念

| 概念 | 说明 |
|------|------|
| Provider | 提供依赖实例的函数 |
| Injector | 声明依赖注入入口的函数 |
| ProviderSet | Provider函数的集合 |
| wire_gen.go | Wire自动生成的依赖注入代码 |

---

## 三、添加新模块的完整流程

### 步骤1: 在对应层添加New构造函数

```go
// repository/order_repository.go
func NewOrderRepository(base *BaseRepository) *OrderRepository {
    return &OrderRepository{BaseRepository: base}
}

// service/order_service.go
func NewOrderService(repo *repository.OrderRepository) *OrderService {
    return &OrderService{repo: repo}
}

// handler/order_handler.go
func NewOrderHandler(svc *service.OrderService) *OrderHandler {
    return &OrderHandler{service: svc}
}
```

### 步骤2: 在providers.go添加Provider函数

```go
// ProvideOrderRepository 提供OrderRepository实例
func ProvideOrderRepository(base *repository.BaseRepository) *repository.OrderRepository {
    return repository.NewOrderRepository(base)
}

// ProvideOrderService 提供OrderService实例
func ProvideOrderService(repo *repository.OrderRepository) *service.OrderService {
    return service.NewOrderService(repo)
}

// ProvideOrderHandler 提供OrderHandler实例
func ProvideOrderHandler(svc *service.OrderService) *handler.OrderHandler {
    return handler.NewOrderHandler(svc)
}
```

### 步骤3: 更新Handlers结构体

```go
type Handlers struct {
    OrderHandler *handler.OrderHandler  // 新增Handler字段
}

func ProvideHandlers(orderHandler *handler.OrderHandler) *Handlers {
    return &Handlers{
        OrderHandler: orderHandler,
    }
}
```

### 步骤4: 更新wire.go的ProviderSet

```go
var ProviderSet = wire.NewSet(
    ProvideBaseRepository,
    ProvideOrderRepository,  // 新增
    ProvideOrderService,     // 新增
    ProvideOrderHandler,     // 新增
    ProvideHandlers,
)
```

### 步骤5: 重新生成wire_gen.go

```bash
~/go/bin/wire ./internal/wire/...
```

### 步骤6: 在router中注册路由

```go
func registerOrderRoutes(rg *gin.RouterGroup, orderHandler *handler.OrderHandler) {
    orders := rg.Group("/orders")
    {
        orders.GET("", orderHandler.List)
        orders.POST("", orderHandler.Create)
        // ...
    }
}
```

---

## 四、Provider函数命名规范

| 类型 | 命名规范 | 示例 |
|------|---------|------|
| Repository | Provide + 模型 + Repository | `ProvideOrderRepository` |
| Service | Provide + 模型 + Service | `ProvideOrderService` |
| Handler | Provide + 模型 + Handler | `ProvideOrderHandler` |

---

## 五、禁止行为

| 禁止 | 正确做法 |
|------|----------|
| 手动修改wire_gen.go | 修改wire.go后重新运行wire生成 |
| 在Provider中包含业务逻辑 | Provider只负责创建实例 |
| 循环依赖 | 重新设计依赖关系 |
| 忘记更新ProviderSet | 添加新Provider后必须更新 |
| 使用装饰性分隔线注释 | 使用简洁单行注释 |

---

## 六、重新生成代码

每次修改wire.go或providers.go后，必须重新生成：

```bash
# 使用wire命令
~/go/bin/wire ./internal/wire/...

# 或使用go generate
go generate ./internal/wire/...
```

---

## 七、依赖关系图

```
database.GetDB()
       │
       ▼
BaseRepository
       │
       ├──────────────────────────────────────┐
       ▼                                      ▼
XxxRepository ─────► XxxService ─────► XxxHandler
                                              │
                                              ▼
                                          Handlers
                                              │
                                              ▼
                                           Router
```
