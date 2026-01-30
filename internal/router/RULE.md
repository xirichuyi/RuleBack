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

### 步骤2: 在认证路由组中注册

```go
func registerAuthenticatedRoutes(rg *gin.RouterGroup) {
    // 临时：不需要认证的路由（仅用于开发测试）
    registerUserRoutes(rg)
    registerOrderRoutes(rg)  // 添加新模块
}
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
    └── [需要认证的路由]
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
users.GET("/:user_id/orders", orderHandler.ListByUser)
```

### 4.2 操作型路由

```go
// 用POST + 动词路径
orders.POST("/:id/cancel", orderHandler.Cancel)
orders.POST("/:id/pay", orderHandler.Pay)
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
| 使用 `New*` 创建Handler | 使用 `Get*` 单例方法 |
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

### 用户模块 (/api/v1/users)
| 方法 | 路径 | 功能 |
|------|------|------|
| GET | /api/v1/users | 获取用户列表 |
| POST | /api/v1/users | 创建用户 |
| GET | /api/v1/users/:id | 获取用户详情 |
| PUT | /api/v1/users/:id | 更新用户 |
| DELETE | /api/v1/users/:id | 删除用户 |

---

## 八、完整router.go模板

```go
package router

import (
    "github.com/gin-gonic/gin"
    "ruleback/internal/handler"
    "ruleback/internal/middleware"
)

// Setup 初始化并配置路由
func Setup() *gin.Engine {
    r := gin.New()

    registerGlobalMiddleware(r)
    registerHealthRoutes(r)
    registerAPIRoutes(r)

    return r
}

// registerGlobalMiddleware 注册全局中间件
func registerGlobalMiddleware(r *gin.Engine) {
    r.Use(middleware.Logger())
    r.Use(middleware.Recovery())
    r.Use(middleware.CORS())
    r.Use(middleware.RequestID())
}

// registerHealthRoutes 注册健康检查路由
func registerHealthRoutes(r *gin.Engine) {
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })
    r.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })
}

// registerAPIRoutes 注册API路由
func registerAPIRoutes(r *gin.Engine) {
    v1 := r.Group("/api/v1")
    {
        registerPublicRoutes(v1)
        registerAuthenticatedRoutes(v1)
    }
}

// registerPublicRoutes 注册公开路由（无需认证）
func registerPublicRoutes(rg *gin.RouterGroup) {
    // auth := rg.Group("/auth")
    // {
    //     h := handler.GetAuthHandler()
    //     auth.POST("/login", h.Login)
    //     auth.POST("/register", h.Register)
    // }
}

// registerAuthenticatedRoutes 注册需要认证的路由
func registerAuthenticatedRoutes(rg *gin.RouterGroup) {
    // 临时：不需要认证的路由（仅用于开发测试）
    registerUserRoutes(rg)
}

// registerUserRoutes 注册用户相关路由
func registerUserRoutes(rg *gin.RouterGroup) {
    h := handler.GetUserHandler()  // 使用Get单例方法

    users := rg.Group("/users")
    {
        users.GET("", h.List)
        users.POST("", h.Create)
        users.GET("/:id", h.GetByID)
        users.PUT("/:id", h.Update)
        users.DELETE("/:id", h.Delete)
    }
}
```
