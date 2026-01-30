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

## 用户模块

### 获取用户列表

```
GET /api/v1/users
```

**查询参数:**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认10 |
| status | int | 否 | 状态筛选 |

**响应:**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [
            {
                "id": 1,
                "username": "admin",
                "email": "admin@example.com",
                "status": 1,
                "created_at": "2024-01-01T00:00:00Z"
            }
        ],
        "total": 100,
        "page": 1,
        "page_size": 10,
        "total_pages": 10
    }
}
```

### 创建用户

```
POST /api/v1/users
```

**请求体:**
```json
{
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "password123"
}
```

**响应:**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 2,
        "username": "newuser",
        "email": "newuser@example.com",
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z"
    }
}
```

### 获取用户详情

```
GET /api/v1/users/:id
```

**路径参数:**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 用户ID |

**响应:**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "status": 1,
        "created_at": "2024-01-01T00:00:00Z"
    }
}
```

### 更新用户

```
PUT /api/v1/users/:id
```

**路径参数:**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 用户ID |

**请求体:**
```json
{
    "email": "updated@example.com",
    "status": 1
}
```

**响应:**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "username": "admin",
        "email": "updated@example.com",
        "status": 1,
        "updated_at": "2024-01-02T00:00:00Z"
    }
}
```

### 删除用户

```
DELETE /api/v1/users/:id
```

**路径参数:**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | int | 用户ID |

**响应:**
```json
{
    "code": 0,
    "message": "删除成功"
}
```

---

## 更新日志

| 日期 | 版本 | 变更内容 |
|------|------|---------|
| 2024-01-01 | v1.0 | 初始版本，包含用户模块 |

<!-- 新增接口时在此处添加 -->
