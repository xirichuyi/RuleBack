# pkg/response 模块 AI 代码生成规则

> **模块职责**: 统一HTTP响应格式，确保所有API返回一致的数据结构

---

## 一、本模块的文件结构

```
pkg/response/
├── response.go   # 响应函数定义
└── RULE.md      # 本规则文件
```

**核心原则**: 所有HTTP响应必须使用本模块函数，禁止直接使用 `c.JSON()`

---

## 二、响应格式标准

### 2.1 基础响应结构

```json
{
    "code": 0,           // 业务状态码，0=成功，非0=失败
    "message": "success", // 响应消息
    "data": {}           // 响应数据（可选）
}
```

### 2.2 分页响应结构

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [],        // 数据列表
        "total": 100,      // 总记录数
        "page": 1,         // 当前页码
        "page_size": 10,   // 每页记录数
        "total_pages": 10  // 总页数
    }
}
```

---

## 三、响应函数使用指南

### 3.1 成功响应函数

| 场景 | 函数 | 示例代码 |
|------|------|---------|
| 无数据返回 | `Success(c)` | 删除成功后 |
| 返回数据 | `SuccessWithData(c, data)` | 查询详情、创建成功 |
| 自定义消息 | `SuccessWithMessage(c, msg)` | 特殊提示 |
| 数据+消息 | `SuccessWithDataAndMessage(c, data, msg)` | 特殊场景 |
| 分页数据 | `SuccessWithPage(c, list, total, page, pageSize)` | 列表查询 |

### 3.2 失败响应函数

| 场景 | 函数 | 示例代码 |
|------|------|---------|
| 业务错误 | `Fail(c, code, msg)` | 参数错误、业务校验失败 |
| 带详情错误 | `FailWithData(c, code, msg, data)` | 验证错误详情 |
| 400错误 | `BadRequest(c, msg)` | 请求格式错误 |
| 401错误 | `Unauthorized(c, msg)` | 未认证 |
| 403错误 | `Forbidden(c, msg)` | 无权限 |
| 404错误 | `NotFound(c, msg)` | 资源不存在 |
| 500错误 | `InternalServerError(c, msg)` | 服务器内部错误 |

---

## 四、在Handler中使用响应

### 4.1 创建操作

```go
func (h *UserHandler) Create(c *gin.Context) {
    var req model.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Fail(c, errors.CodeInvalidParams, "参数错误: "+err.Error())
        return
    }

    user, err := h.service.Create(&req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithData(c, user)
}
```

### 4.2 查询详情

```go
func (h *UserHandler) GetByID(c *gin.Context) {
    id, err := h.parseIDParam(c)
    if err != nil {
        response.Fail(c, errors.CodeInvalidParams, "无效的ID")
        return
    }

    user, err := h.service.GetByID(id)
    if err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithData(c, user)
}
```

### 4.3 分页查询

```go
func (h *UserHandler) List(c *gin.Context) {
    var query model.UserListQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        response.Fail(c, errors.CodeInvalidParams, "参数错误: "+err.Error())
        return
    }

    users, total, err := h.service.List(&query)
    if err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithPage(c, users, total, query.Page, query.PageSize)
}
```

### 4.4 更新操作

```go
func (h *UserHandler) Update(c *gin.Context) {
    id, err := h.parseIDParam(c)
    if err != nil {
        response.Fail(c, errors.CodeInvalidParams, "无效的ID")
        return
    }

    var req model.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Fail(c, errors.CodeInvalidParams, "参数错误: "+err.Error())
        return
    }

    user, err := h.service.Update(id, &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    response.SuccessWithData(c, user)
}
```

### 4.5 删除操作

```go
func (h *UserHandler) Delete(c *gin.Context) {
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
```

---

## 五、错误处理规范

### 5.1 Handler统一错误处理方法

```go
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

### 5.2 使用预定义HTTP状态码响应

```go
// 认证失败
response.Unauthorized(c, "请先登录")

// 权限不足
response.Forbidden(c, "无权访问该资源")

// 资源不存在
response.NotFound(c, "用户不存在")

// 服务器错误（不暴露内部细节）
response.InternalServerError(c, "服务器繁忙，请稍后重试")
```

---

## 六、禁止行为

| 禁止 | 正确做法 |
|------|----------|
| ❌ 直接使用c.JSON() | ✅ 使用response包函数 |
| ❌ 修改Response结构体 | ✅ 如需扩展请联系架构组 |
| ❌ 在message中返回敏感信息 | ✅ 使用友好的错误提示 |
| ❌ 暴露系统内部错误详情 | ✅ 记录日志，返回友好消息 |
| ❌ 硬编码错误码 | ✅ 使用errors包的常量 |

**错误示例**:
```go
// ❌ 错误 - 直接使用c.JSON
func (h *UserHandler) GetUser(c *gin.Context) {
    user, _ := h.service.GetByID(id)
    c.JSON(200, gin.H{"user": user})
}

// ❌ 错误 - 暴露内部错误
func (h *UserHandler) GetUser(c *gin.Context) {
    user, err := h.service.GetByID(id)
    if err != nil {
        response.Fail(c, 500, "SQL Error: "+err.Error())  // 暴露了SQL错误
        return
    }
}

// ✅ 正确
func (h *UserHandler) GetUser(c *gin.Context) {
    user, err := h.service.GetByID(id)
    if err != nil {
        h.handleError(c, err)  // 统一处理，不暴露内部细节
        return
    }
    response.SuccessWithData(c, user)
}
```

---

## 七、已存在的响应函数

### 成功响应
| 函数 | 签名 |
|------|------|
| Success | `Success(c *gin.Context)` |
| SuccessWithData | `SuccessWithData(c *gin.Context, data interface{})` |
| SuccessWithMessage | `SuccessWithMessage(c *gin.Context, message string)` |
| SuccessWithDataAndMessage | `SuccessWithDataAndMessage(c *gin.Context, data interface{}, message string)` |
| SuccessWithPage | `SuccessWithPage(c *gin.Context, list interface{}, total int64, page, pageSize int)` |

### 失败响应
| 函数 | 签名 |
|------|------|
| Fail | `Fail(c *gin.Context, code int, message string)` |
| FailWithData | `FailWithData(c *gin.Context, code int, message string, data interface{})` |

### HTTP状态码响应
| 函数 | HTTP状态码 |
|------|-----------|
| BadRequest | 400 |
| Unauthorized | 401 |
| Forbidden | 403 |
| NotFound | 404 |
| InternalServerError | 500 |

---

## 八、扩展规范

如需添加新的响应函数，必须遵循：

1. 函数签名必须以 `c *gin.Context` 作为第一个参数
2. 函数名必须清晰表达用途（Success/Fail 前缀）
3. 必须使用 Response 结构体
4. 必须添加完整的注释文档
5. 必须在本 RULE.md 中更新函数列表
