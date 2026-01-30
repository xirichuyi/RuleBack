# internal/router 模块 AI 代码生成规则

> **模块职责**: 路由配置层，负责URL路径与Handler的映射关系

---

## 一、本模块的文件结构

```
internal/router/
├── router.go   # 路由配置
└── RULE.md    # 本规则文件
```

**规则**: 所有路由配置统一在 router.go 文件中

---

## 二、添加新模块路由的完整流程

### 步骤1: 创建路由注册函数

```go
// registerOrderRoutes 注册订单相关路由
func registerOrderRoutes(rg *gin.RouterGroup, handlers *wire.Handlers) {
    orders := rg.Group("/orders")
    {
        orders.GET("", handlers.OrderHandler.List)
        orders.POST("", handlers.OrderHandler.Create)
        orders.GET("/:id", handlers.OrderHandler.GetByID)
        orders.PUT("/:id", handlers.OrderHandler.Update)
        orders.DELETE("/:id", handlers.OrderHandler.Delete)
    }
}
```

### 步骤2: 在Setup中注册自定义路由

```go
// 在 cmd/server/bootstrap.go 中
r := router.Setup(handlers, registerOrderRoutes)
```

### 步骤3: 如需认证，使用 RegisterAuthenticatedRoutes 包装

```go
r := router.Setup(handlers,
    router.RegisterAuthenticatedRoutes(registerOrderRoutes),
)
```

---

## 三、路由结构规范

### 3.1 标准路由层次

```
/
├── /health                    # 健康检查（公开）
├── /ping                      # 存活检查（公开）
└── /api/v1                    # API版本1
    ├── /auth                  # 认证相关（公开）
    │   ├── POST /login
    │   └── POST /register
    └── [业务路由]
        ├── /users             # 用户管理
        └── /orders            # 订单管理
```

### 3.2 标准CRUD路由

| 操作 | HTTP方法 | 路径 | Handler方法 |
|------|----------|------|-------------|
| 列表 | GET | /xxxs | List |
| 创建 | POST | /xxxs | Create |
| 详情 | GET | /xxxs/:id | GetByID |
| 更新 | PUT | /xxxs/:id | Update |
| 删除 | DELETE | /xxxs/:id | Delete |

---

## 四、特殊路由规范

### 4.1 嵌套资源路由

```go
// GET /users/:user_id/orders
users.GET("/:user_id/orders", handlers.OrderHandler.ListByUser)
```

### 4.2 操作型路由

```go
// 用POST + 动词路径
orders.POST("/:id/cancel", handlers.OrderHandler.Cancel)
orders.POST("/:id/pay", handlers.OrderHandler.Pay)
```

---

## 五、路由命名规范

| 规则 | 正确示例 | 错误示例 |
|------|---------|---------|
| 使用小写字母 | `/users` | `/Users` |
| 使用复数形式 | `/orders` | `/order` |
| 多词用连字符 | `/user-profiles` | `/user_profiles` |
| 避免动词 | `/orders/:id` | `/getOrder/:id` |

---

## 六、禁止行为

| 禁止 | 正确做法 |
|------|----------|
| 在router中编写业务逻辑 | 只做路由映射 |
| 使用匿名函数作为Handler | 使用Handler方法 |
| 使用装饰性分隔线注释 | 使用简洁单行注释 |

---

## 七、已存在的路由

### 系统路由
| 方法 | 路径 | 功能 |
|------|------|------|
| GET | /health | 健康检查 |
| GET | /ping | 存活检查 |

---

## 八、完整使用示例

### 定义路由注册函数

```go
// 在 internal/router/router.go 或单独文件中定义
func registerUserRoutes(rg *gin.RouterGroup, handlers *wire.Handlers) {
    users := rg.Group("/users")
    {
        users.GET("", handlers.UserHandler.List)
        users.POST("", handlers.UserHandler.Create)
        users.GET("/:id", handlers.UserHandler.GetByID)
        users.PUT("/:id", handlers.UserHandler.Update)
        users.DELETE("/:id", handlers.UserHandler.Delete)
    }
}

func registerProductRoutes(rg *gin.RouterGroup, handlers *wire.Handlers) {
    products := rg.Group("/products")
    {
        products.GET("", handlers.ProductHandler.List)
        products.POST("", handlers.ProductHandler.Create)
        products.GET("/:id", handlers.ProductHandler.GetByID)
        products.PUT("/:id", handlers.ProductHandler.Update)
        products.DELETE("/:id", handlers.ProductHandler.Delete)
    }
}
```

### 在启动时注册

```go
// 在 cmd/server/bootstrap.go 的 startServer 函数中
r := router.Setup(handlers,
    registerUserRoutes,      // 公开路由
    router.RegisterAuthenticatedRoutes(registerProductRoutes),  // 需要认证
)
```
