# internal/model 模块 AI 代码生成规则

> **模块职责**: 定义数据模型、请求/响应结构体、查询参数

---

## 一、本模块的文件结构

```
internal/model/
├── base.go     # 基础模型（BaseModel）和通用类型（勿修改）
├── xxx.go      # 业务模型（按需创建）
└── RULE.md     # 本规则文件
```

**规则**: 每个业务实体对应一个文件

---

## 二、创建新模型的完整流程

### 步骤1: 创建模型文件

文件命名: `{实体名小写}.go`，如 `order.go`

### 步骤2: 编写数据库模型

```go
package model

// Order 订单模型
type Order struct {
    BaseModel
    UserID  uint        `gorm:"index;not null" json:"user_id"`
    OrderNo string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_no"`
    Amount  float64     `gorm:"type:decimal(10,2);not null" json:"amount"`
    Status  OrderStatus `gorm:"type:tinyint;default:0" json:"status"`
}

// TableName 指定表名
func (Order) TableName() string {
    return "orders"
}
```

### 步骤3: 定义枚举类型（如需要）

```go
// OrderStatus 订单状态
type OrderStatus int8

const (
    OrderStatusPending   OrderStatus = 0
    OrderStatusPaid      OrderStatus = 1
    OrderStatusCancelled OrderStatus = 2
)

func (s OrderStatus) String() string {
    switch s {
    case OrderStatusPending:
        return "pending"
    case OrderStatusPaid:
        return "paid"
    case OrderStatusCancelled:
        return "cancelled"
    default:
        return "unknown"
    }
}
```

### 步骤4: 编写请求结构体

```go
// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
    ProductID uint   `json:"product_id" binding:"required"`
    Quantity  int    `json:"quantity" binding:"required,min=1"`
    Remark    string `json:"remark" binding:"max=500"`
}

// UpdateOrderRequest 更新订单请求
type UpdateOrderRequest struct {
    Status *OrderStatus `json:"status" binding:"omitempty,oneof=0 1 2"`
    Remark *string      `json:"remark" binding:"omitempty,max=500"`
}
```

### 步骤5: 编写查询结构体

```go
// OrderListQuery 订单列表查询参数
type OrderListQuery struct {
    PageQuery
    SortQuery
    UserID  *uint        `form:"user_id"`
    Status  *OrderStatus `form:"status"`
    OrderNo string       `form:"order_no"`
}
```

---

## 三、标签使用规范

### GORM 标签

| 类型 | 标签 | 示例 |
|------|------|------|
| 字符串 | `type:varchar(长度)` | `gorm:"type:varchar(100)"` |
| 小数 | `type:decimal(总位数,小数位)` | `gorm:"type:decimal(10,2)"` |
| 非空 | `not null` | `gorm:"not null"` |
| 唯一索引 | `uniqueIndex` | `gorm:"uniqueIndex"` |
| 普通索引 | `index` | `gorm:"index"` |

### Binding 验证标签

| 标签 | 说明 |
|------|------|
| `required` | 必填 |
| `omitempty` | 可选 |
| `min=n` | 最小值/长度 |
| `max=n` | 最大值/长度 |
| `email` | 邮箱格式 |
| `oneof=a b c` | 枚举值 |

---

## 四、命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 模型结构体 | PascalCase，单数 | `User`, `Order` |
| 表名 | snake_case，复数 | `users`, `orders` |
| 创建请求 | Create + 模型 + Request | `CreateOrderRequest` |
| 更新请求 | Update + 模型 + Request | `UpdateOrderRequest` |
| 查询参数 | 模型 + ListQuery | `OrderListQuery` |

---

## 五、已存在的基础类型

```go
// BaseModel 所有模型必须嵌入
type BaseModel struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Status 通用状态
type Status int8
const (
    StatusDisabled Status = 0
    StatusEnabled  Status = 1
)

// PageQuery 分页参数
// SortQuery 排序参数
```

---

## 六、禁止行为

| 禁止 | 原因 |
|------|------|
| 在模型中编写业务逻辑 | 业务逻辑应在 service 层 |
| Password 字段设为 `json:"password"` | 敏感信息泄露 |
| 修改 BaseModel 定义 | 会影响所有模型 |
| 使用装饰性分隔线注释 | 使用简洁单行注释 |

---

## 七、完整模型文件模板

```go
package model

// Xxx Xxx模型
type Xxx struct {
    BaseModel
    Field1 string `gorm:"type:varchar(100);not null" json:"field1"`
    Status Status `gorm:"type:tinyint;default:1" json:"status"`
}

func (Xxx) TableName() string {
    return "xxxs"
}

// CreateXxxRequest 创建Xxx请求
type CreateXxxRequest struct {
    Field1 string `json:"field1" binding:"required,max=100"`
}

// UpdateXxxRequest 更新Xxx请求
type UpdateXxxRequest struct {
    Field1 *string `json:"field1" binding:"omitempty,max=100"`
    Status *Status `json:"status" binding:"omitempty,oneof=0 1"`
}

// XxxListQuery Xxx列表查询参数
type XxxListQuery struct {
    PageQuery
    SortQuery
    Field1 string  `form:"field1"`
    Status *Status `form:"status"`
}
```
