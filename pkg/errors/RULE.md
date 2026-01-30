# pkg/errors 模块 AI 代码生成规则

> **模块职责**: 统一错误处理，定义业务错误码和自定义错误类型

---

## 一、本模块的文件结构

```
pkg/errors/
├── errors.go   # 错误码和错误类型定义
└── RULE.md    # 本规则文件
```

---

## 二、错误码分段规范

| 范围 | 用途 | 前缀 |
|------|------|------|
| 0 | 成功 | CodeSuccess |
| 10000-10999 | 通用错误 | Code + 描述 |
| 20000-20999 | 用户模块 | CodeUser + 描述 |
| 30000-30999 | 模块A | 根据业务定义 |
| 40000-40999 | 模块B | 根据业务定义 |

---

## 三、添加新错误码的完整流程

### 步骤1: 在errors.go中添加错误码常量

```go
// Order模块错误码 (30xxx)
const (
    CodeOrderNotFound = 30001
    CodeOrderExpired  = 30002
    CodeOrderPaid     = 30003
    CodeOrderCanceled = 30004
)
```

### 步骤2: 在codeMessages映射中添加默认消息

```go
var codeMessages = map[int]string{
    // ... 现有映射 ...
    CodeOrderNotFound: "订单不存在",
    CodeOrderExpired:  "订单已过期",
    CodeOrderPaid:     "订单已支付",
    CodeOrderCanceled: "订单已取消",
}
```

### 步骤3: 添加预定义错误实例（可选）

```go
var (
    ErrOrderNotFound = NewWithCode(CodeOrderNotFound)
    ErrOrderExpired  = NewWithCode(CodeOrderExpired)
)
```

---

## 四、错误创建方式

| 方式 | 函数 | 使用场景 |
|------|------|---------|
| 预定义错误 | `errors.ErrUserNotFound` | 标准错误，使用默认消息 |
| 错误码创建 | `errors.NewWithCode(code)` | 使用默认消息 |
| 自定义消息 | `errors.New(code, msg)` | 需要自定义消息 |
| 包装错误 | `errors.Wrap(code, msg, err)` | 保留原始错误 |

---

## 五、在各层使用错误

### Repository层 - 直接返回原始错误

```go
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
    var user model.User
    if err := r.DB().First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

### Service层 - 转换为业务错误

```go
func (s *UserService) GetByID(id uint) (*model.User, error) {
    user, err := s.repo.GetByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, apperrors.ErrUserNotFound
        }
        return nil, apperrors.Wrap(apperrors.CodeDatabaseError, "获取用户失败", err)
    }
    return user, nil
}
```

### Handler层 - 转换为HTTP响应

```go
func (h *UserHandler) handleError(c *gin.Context, err error) {
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

## 六、禁止行为

| 禁止 | 正确做法 |
|------|----------|
| 硬编码错误码数字 | 使用错误码常量 |
| 不同模块使用相同错误码 | 按分段规范使用 |
| 在错误消息中暴露敏感信息 | 使用Wrap保留原错误，返回友好消息 |
| 修改已存在的错误码值 | 只能新增，不能修改 |

---

## 七、已存在的错误码

### 通用错误码 (10xxx)

| 常量 | 值 | 说明 |
|------|-----|------|
| CodeSuccess | 0 | 成功 |
| CodeInvalidParams | 10001 | 参数错误 |
| CodeUnauthorized | 10002 | 未认证 |
| CodeForbidden | 10003 | 无权限 |
| CodeNotFound | 10004 | 资源不存在 |
| CodeDatabaseError | 10007 | 数据库错误 |

### 用户模块错误码 (20xxx)

| 常量 | 值 | 说明 |
|------|-----|------|
| CodeUserNotFound | 20001 | 用户不存在 |
| CodeUserExists | 20002 | 用户已存在 |
| CodePasswordIncorrect | 20004 | 密码错误 |
| CodeTokenExpired | 20005 | Token已过期 |

---

## 八、已存在的错误函数

| 函数 | 说明 |
|------|------|
| `New(code, message)` | 创建带自定义消息的错误 |
| `NewWithCode(code)` | 创建使用默认消息的错误 |
| `Wrap(code, message, err)` | 包装原始错误 |
| `GetAppError(err)` | 从error中提取AppError |
