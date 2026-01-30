# internal/middleware 模块 AI 代码生成规则

> **模块职责**: HTTP中间件层，处理请求拦截、认证授权、日志记录等横切关注点

---

## 一、本模块的文件结构

```
internal/middleware/
├── middleware.go   # 所有中间件定义
└── RULE.md        # 本规则文件
```

---

## 二、创建新中间件的完整流程

### 无参数中间件

```go
// Xxx Xxx中间件
func Xxx() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置处理
        c.Next()
        // 后置处理
    }
}
```

### 有参数中间件

```go
// RateLimit 请求限流中间件
func RateLimit(limit int, window int) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 检查是否超过限制
        if !allowed {
            response.Fail(c, errors.CodeRateLimited, "请求过于频繁")
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 在router中注册

**全局中间件**:
```go
func registerGlobalMiddleware(r *gin.Engine) {
    r.Use(middleware.Logger())
    r.Use(middleware.Recovery())
    r.Use(middleware.CORS())
    r.Use(middleware.RequestID())
}
```

**路由组中间件**:
```go
authenticated := rg.Group("")
authenticated.Use(middleware.Auth())
```

---

## 三、中间件分类

| 类型 | 中间件 | 功能 |
|------|--------|------|
| 基础 | `Logger()` | 请求日志记录 |
| 基础 | `Recovery()` | panic恢复 |
| 基础 | `CORS()` | 跨域处理 |
| 基础 | `RequestID()` | 请求ID追踪 |
| 认证 | `Auth()` | JWT认证 |
| 认证 | `RequireRole(roles...)` | 角色权限 |
| 流控 | `RateLimit(limit, window)` | 请求限流 |

---

## 四、上下文数据规范

| 键名 | 类型 | 说明 | 设置者 |
|------|------|------|--------|
| `user_id` | uint | 用户ID | Auth |
| `username` | string | 用户名 | Auth |
| `request_id` | string | 请求ID | RequestID |

获取方式:
```go
userID := c.GetUint("user_id")
requestID := c.GetString("request_id")
```

---

## 五、中断请求规范

必须同时调用响应函数和 `c.Abort()`:

```go
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            response.Unauthorized(c, "未提供认证信息")
            c.Abort()
            return
        }
        c.Next()
    }
}
```

---

## 六、禁止行为

| 禁止 | 正确做法 |
|------|----------|
| 在中间件中编写业务逻辑 | 只做认证、授权、日志等横切关注点 |
| 直接使用c.JSON返回响应 | 使用response包的函数 |
| 拦截请求后忘记调用c.Abort() | response.Xxx() + c.Abort() + return |
| 使用装饰性分隔线注释 | 使用简洁单行注释 |

---

## 七、已存在的中间件

| 中间件 | 函数签名 | 功能 |
|--------|---------|------|
| 日志 | `Logger()` | 记录请求日志 |
| 恢复 | `Recovery()` | panic恢复 |
| 跨域 | `CORS()` | 跨域请求处理 |
| 请求ID | `RequestID()` | 生成请求追踪ID |
| JWT认证 | `Auth()` | JWT Token验证 |
| 角色权限 | `RequireRole(roles...)` | 角色权限检查 |
| 限流 | `RateLimit(limit, window)` | 请求频率限制 |
| 错误处理 | `ErrorHandler()` | 全局错误处理 |
