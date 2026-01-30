# pkg/logger 模块 AI 代码生成规则

> **模块职责**: 统一日志记录，提供结构化日志功能

---

## 一、本模块的文件结构

```
pkg/logger/
├── logger.go   # 日志函数定义
└── RULE.md    # 本规则文件
```

**核心原则**: 所有日志必须使用本模块函数，禁止使用 `fmt.Println` 或 `log` 包

---

## 二、日志级别规范

| 级别 | 用途 | 场景示例 |
|------|------|---------|
| Debug | 调试信息 | 变量值、流程跟踪、SQL语句 |
| Info | 正常信息 | 用户操作、系统事件、业务成功 |
| Warn | 警告信息 | 非致命错误、参数异常、可恢复问题 |
| Error | 错误信息 | 操作失败、异常情况、需要关注 |
| Fatal | 致命错误 | 程序无法继续运行 |

---

## 三、日志记录规范

### 3.1 结构化日志（推荐）

使用字段函数记录结构化数据：

```go
// 基础用法
logger.Info("用户登录成功",
    logger.Uint("user_id", userID),
)

// 多个字段
logger.Info("创建订单",
    logger.String("order_no", orderNo),
    logger.Uint("user_id", userID),
    logger.Float64("amount", 99.99),
)

// 记录错误
logger.Error("数据库查询失败",
    logger.Err(err),
    logger.String("table", "users"),
    logger.String("operation", "select"),
)
```

### 3.2 字段函数

| 函数 | 类型 | 示例 |
|------|------|------|
| String | string | `logger.String("name", "张三")` |
| Int | int | `logger.Int("count", 100)` |
| Int64 | int64 | `logger.Int64("id", 1234567890)` |
| Uint | uint | `logger.Uint("user_id", 1)` |
| Float64 | float64 | `logger.Float64("price", 99.99)` |
| Bool | bool | `logger.Bool("active", true)` |
| Err | error | `logger.Err(err)` |
| Field | any | `logger.Field("data", anyValue)` |

### 3.3 格式化日志（简单场景）

```go
// 简单格式化
logger.Infof("用户 %d 登录成功", userID)

// 复杂格式化
logger.Errorf("查询用户 %d 失败: %v", userID, err)
```

---

## 四、各层日志规范

### 4.1 Handler层

```go
func (h *UserHandler) Create(c *gin.Context) {
    // 参数错误 - Warn级别
    var req model.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.Warn("创建用户参数错误",
            logger.Err(err),
            logger.String("ip", c.ClientIP()),
        )
        response.Fail(c, errors.CodeInvalidParams, "参数错误")
        return
    }

    user, err := h.service.Create(&req)
    if err != nil {
        // 业务错误在Service层已记录，Handler不重复记录
        h.handleError(c, err)
        return
    }

    // 成功 - Info级别（可选）
    logger.Info("用户创建成功",
        logger.Uint("user_id", user.ID),
        logger.String("username", user.Username),
    )
    response.SuccessWithData(c, user)
}
```

### 4.2 Service层

```go
func (s *UserService) Create(req *model.CreateUserRequest) (*model.User, error) {
    // 方法开始 - Debug级别
    logger.Debug("开始创建用户",
        logger.String("username", req.Username),
    )

    // 业务校验失败 - Warn级别
    exists, err := s.repo.ExistsByUsername(req.Username)
    if err != nil {
        logger.Error("检查用户名失败",
            logger.Err(err),
            logger.String("username", req.Username),
        )
        return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "检查用户名失败", err)
    }
    if exists {
        logger.Warn("用户名已存在",
            logger.String("username", req.Username),
        )
        return nil, apperrors.ErrUserExists
    }

    // 数据库操作失败 - Error级别
    if err := s.repo.Create(user); err != nil {
        logger.Error("创建用户失败",
            logger.Err(err),
            logger.String("username", req.Username),
        )
        return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "创建用户失败", err)
    }

    // 操作成功 - Info级别
    logger.Info("用户创建成功",
        logger.Uint("user_id", user.ID),
    )
    return user, nil
}
```

### 4.3 Repository层

Repository层通常只在Debug级别记录：

```go
func (r *UserRepository) Create(user *model.User) error {
    logger.Debug("插入用户记录",
        logger.String("username", user.Username),
    )
    return r.DB().Create(user).Error
}
```

---

## 五、日志字段命名规范

### 5.1 标准字段名

| 场景 | 字段名 | 示例 |
|------|--------|------|
| 用户ID | user_id | `logger.Uint("user_id", id)` |
| 请求ID | request_id | `logger.String("request_id", rid)` |
| 错误信息 | error | `logger.Err(err)` |
| 操作类型 | operation | `logger.String("operation", "create")` |
| 耗时(毫秒) | elapsed_ms | `logger.Float64("elapsed_ms", ms)` |
| IP地址 | ip | `logger.String("ip", clientIP)` |
| 表名 | table | `logger.String("table", "users")` |
| 订单号 | order_no | `logger.String("order_no", no)` |

### 5.2 消息格式规范

- 使用中文或英文，保持项目统一
- 消息简洁明了，说明发生了什么
- **不要**在消息中包含变量值，使用Field记录

```go
// ✅ 正确
logger.Info("用户登录成功", logger.Uint("user_id", userID))

// ❌ 错误
logger.Info(fmt.Sprintf("用户 %d 登录成功", userID))
```

---

## 六、必须记录日志的场景

| 场景 | 级别 | 说明 |
|------|------|------|
| 请求参数错误 | Warn | Handler层参数绑定失败 |
| 业务校验失败 | Warn | Service层业务规则不满足 |
| 数据库操作失败 | Error | Repository返回错误 |
| 外部API调用失败 | Error | HTTP请求、RPC调用失败 |
| 重要业务操作成功 | Info | 创建、更新、删除操作 |
| 程序启动/关闭 | Info | 服务生命周期事件 |
| 未知异常 | Error | 未预期的错误 |

---

## 七、禁止行为

| 禁止 | 正确做法 |
|------|----------|
| ❌ 使用fmt.Println或log包 | ✅ 使用logger包 |
| ❌ 在日志中记录敏感信息 | ✅ 脱敏处理或不记录 |
| ❌ 在循环中记录大量日志 | ✅ 循环外汇总记录 |
| ❌ 日志级别使用错误 | ✅ 按规范选择级别 |
| ❌ 在消息中拼接变量 | ✅ 使用字段函数 |

**错误示例**:
```go
// ❌ 错误 - 使用fmt.Println
fmt.Println("用户登录成功")

// ❌ 错误 - 记录敏感信息
logger.Info("用户注册",
    logger.String("password", req.Password),  // 禁止记录密码
)

// ❌ 错误 - 循环中记录大量日志
for _, item := range items {
    logger.Debug("处理item", logger.Field("item", item))
}

// ❌ 错误 - 正常流程使用Error级别
logger.Error("用户登录成功")  // 应该用Info

// ✅ 正确
logger.Info("用户登录成功", logger.Uint("user_id", userID))

// ✅ 正确 - 不记录敏感信息
logger.Info("用户注册", logger.String("username", req.Username))

// ✅ 正确 - 汇总记录
logger.Debug("处理items完成", logger.Int("count", len(items)))
```

---

## 八、已存在的日志函数

### 结构化日志

| 函数 | 签名 |
|------|------|
| Debug | `Debug(msg string, fields ...zap.Field)` |
| Info | `Info(msg string, fields ...zap.Field)` |
| Warn | `Warn(msg string, fields ...zap.Field)` |
| Error | `Error(msg string, fields ...zap.Field)` |
| Fatal | `Fatal(msg string, fields ...zap.Field)` |

### 格式化日志

| 函数 | 签名 |
|------|------|
| Debugf | `Debugf(template string, args ...interface{})` |
| Infof | `Infof(template string, args ...interface{})` |
| Warnf | `Warnf(template string, args ...interface{})` |
| Errorf | `Errorf(template string, args ...interface{})` |
| Fatalf | `Fatalf(template string, args ...interface{})` |

### 字段函数

| 函数 | 签名 |
|------|------|
| Field | `Field(key string, value interface{}) zap.Field` |
| String | `String(key string, val string) zap.Field` |
| Int | `Int(key string, val int) zap.Field` |
| Int64 | `Int64(key string, val int64) zap.Field` |
| Uint | `Uint(key string, val uint) zap.Field` |
| Float64 | `Float64(key string, val float64) zap.Field` |
| Bool | `Bool(key string, val bool) zap.Field` |
| Err | `Err(err error) zap.Field` |

### 辅助函数

| 函数 | 签名 | 说明 |
|------|------|------|
| WithFields | `WithFields(fields ...zap.Field) *zap.Logger` | 创建带预设字段的logger |
| Sync | `Sync()` | 同步日志缓冲区 |

---

## 九、带上下文的日志

```go
// 创建带预设字段的logger
reqLogger := logger.WithFields(
    logger.String("request_id", requestID),
    logger.String("user_id", userID),
)

// 后续日志自动包含预设字段
reqLogger.Info("开始处理请求")
reqLogger.Info("处理完成")
```
