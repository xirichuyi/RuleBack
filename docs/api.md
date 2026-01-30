# RuleBack API 文档

> 本文档记录所有API接口，每次接口变更后需要更新此文档

---

## 通用说明

### 基础URL

```
开发环境: http://localhost:8080/api/v1
生产环境: https://api.example.com/api/v1
```

### 响应格式

所有接口返回统一的JSON格式：

```json
{
    "code": 0,
    "message": "success",
    "data": {}
}
```

### 状态码说明

| code | 说明 |
|------|------|
| 0 | 成功 |
| 10001 | 参数错误 |
| 10002 | 未认证 |
| 10003 | 无权限 |
| 10004 | 资源不存在 |
| 10006 | 内部错误 |
| 10007 | 数据库错误 |

### 分页响应

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [],
        "total": 100,
        "page": 1,
        "page_size": 10,
        "total_pages": 10
    }
}
```

### 认证方式

需要认证的接口在Header中携带Token：

```
Authorization: Bearer <token>
```

---

## 系统接口

### 健康检查

```
GET /health
```

**响应:**
```json
{
    "status": "ok"
}
```

### 存活检查

```
GET /ping
```

**响应:**
```
pong
```

---

## 业务接口

> 框架初始状态不包含业务接口，请根据业务需求添加。
>
> 添加新模块时，请参考下方模板和 `docs/RULE.md` 中的文档格式规范。

### 接口文档模板

以下是标准CRUD接口的文档模板，添加新模块时可参考：

```
## 模块名称

### 获取列表

GET /api/v1/xxxs

**查询参数:**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认10 |

### 创建记录

POST /api/v1/xxxs

### 获取详情

GET /api/v1/xxxs/:id

### 更新记录

PUT /api/v1/xxxs/:id

### 删除记录

DELETE /api/v1/xxxs/:id
```

---

## 更新日志

| 日期 | 版本 | 变更内容 |
|------|------|---------|
| 2024-01-01 | v1.0 | 框架初始版本 |

<!-- 新增接口时在此处添加 -->
