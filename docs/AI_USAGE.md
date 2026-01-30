# AI 代码生成使用指南

本文档介绍如何配合 AI（如 Claude、ChatGPT）使用 RuleBack 框架进行高效的代码生成。

## 核心理念

RuleBack 框架通过 `RULE.md` 规则文件来指导 AI 生成一致、高质量的代码。每个模块都有对应的规则文件，定义了：

- 文件结构和命名规范
- 代码模板和实现模式
- 禁止行为和最佳实践

## 快速开始

### 第一步：让 AI 了解项目

在开始生成代码前，让 AI 阅读主规则文件：

```
请阅读项目根目录的 CLAUDE.md 文件，了解项目的架构和代码规范。
```

### 第二步：描述需求

清晰地描述你需要的功能：

```
帮我创建一个订单(Order)模块，包含以下字段：
- order_no: 订单号 (string, 唯一)
- user_id: 用户ID (uint, 必填)
- amount: 金额 (decimal)
- status: 状态 (待支付/已支付/已取消)
- remark: 备注 (string, 可选)

需要完整的 CRUD 接口。
```

### 第三步：AI 自动生成

AI 会按照框架规范自动生成完整的模块代码，包括：

1. 数据模型 (`internal/model/order.go`)
2. 数据访问层 (`internal/repository/order_repository.go`)
3. 业务逻辑层 (`internal/service/order_service.go`)
4. HTTP处理层 (`internal/handler/order_handler.go`)
5. Wire 依赖注入配置更新
6. 路由注册更新
7. 数据库迁移配置

### 第四步：生成 Wire 代码

代码生成完成后，运行 Wire 命令：

```bash
~/go/bin/wire ./internal/wire/...
```

## 常用提示词模板

### 创建新模块

```
请先阅读 CLAUDE.md 了解项目规范。
帮我创建一个[模块名]模块，包含以下字段：
- 字段1: 类型 (约束说明)
- 字段2: 类型 (约束说明)
- ...

需要完整的 CRUD 接口。
```

### 添加业务功能

```
在[模块名]模块中添加一个"[功能名]"功能：
- 功能描述
- 业务规则
- 预期行为
```

### 添加查询接口

```
给[模块名]模块添加一个按[字段]查询的功能。
```

### 添加批量操作

```
给[模块名]模块添加批量删除功能，支持传入多个ID。
```

## 示例对话

### 示例1：创建商品模块

**用户：**
```
请先阅读 CLAUDE.md 了解项目规范。
帮我创建一个商品(Product)模块，包含以下字段：
- name: 商品名称 (string, 必填, 最大100字符)
- description: 商品描述 (string, 可选, 最大1000字符)
- price: 价格 (decimal(10,2), 必填)
- stock: 库存数量 (int, 默认0)
- category_id: 分类ID (uint, 可选)
- status: 状态 (上架/下架, 默认下架)

需要完整的 CRUD 接口，列表接口支持按名称模糊搜索和分类筛选。
```

**AI 会生成：**
- 商品模型定义（包含枚举、请求结构体、查询结构体）
- Repository 层（包含按分类查询、模糊搜索）
- Service 层（包含业务校验）
- Handler 层（完整的 CRUD 接口）
- Wire 配置更新
- 路由注册

### 示例2：添加业务功能

**用户：**
```
在订单模块中添加一个"取消订单"功能：
- 只有待支付状态的订单可以取消
- 取消后状态变为"已取消"
- 需要记录取消时间
- 提供 POST /api/v1/orders/:id/cancel 接口
```

**AI 会生成：**
- 在 Model 中添加 CanceledAt 字段
- 在 Repository 中添加 UpdateStatus 方法
- 在 Service 中添加 Cancel 方法（包含状态校验）
- 在 Handler 中添加 Cancel 接口
- 更新路由注册

## 规则文件索引

| 模块 | 规则文件 | 说明 |
|------|---------|------|
| 全局 | `CLAUDE.md` | 项目全局规范 |
| 模型 | `internal/model/RULE.md` | 数据模型定义规范 |
| 数据层 | `internal/repository/RULE.md` | Repository 实现规范 |
| 业务层 | `internal/service/RULE.md` | Service 实现规范 |
| 接口层 | `internal/handler/RULE.md` | Handler 实现规范 |
| 路由 | `internal/router/RULE.md` | 路由注册规范 |
| 依赖注入 | `internal/wire/RULE.md` | Wire 配置规范 |
| 中间件 | `internal/middleware/RULE.md` | 中间件实现规范 |
| 配置 | `internal/config/RULE.md` | 配置管理规范 |
| 错误 | `pkg/errors/RULE.md` | 错误码定义规范 |
| 响应 | `pkg/response/RULE.md` | 响应格式规范 |
| 日志 | `pkg/logger/RULE.md` | 日志记录规范 |
| 入口 | `cmd/server/RULE.md` | 程序入口规范 |

## 最佳实践

### 1. 提供完整的字段信息

```
# 好的描述
- price: 价格 (decimal(10,2), 必填, 大于0)

# 不够详细
- price: 价格
```

### 2. 说明业务规则

```
# 好的描述
创建订单时自动生成订单号，格式：ORD + 年月日 + 6位随机数

# 不够详细
需要订单号
```

### 3. 指定接口路径（如有特殊要求）

```
# 好的描述
提供 POST /api/v1/orders/:id/pay 接口用于支付

# 让 AI 自动处理
需要支付功能
```

### 4. 分步骤进行复杂功能

对于复杂功能，建议分步骤进行：

```
第一步：先创建基础的订单模块，包含基本的 CRUD
第二步：添加订单状态流转功能
第三步：添加订单统计功能
```

## 常见问题

### Q: AI 生成的代码不符合规范怎么办？

让 AI 重新阅读相关的 RULE.md 文件：

```
请重新阅读 internal/service/RULE.md，按照规范修改刚才生成的 Service 代码。
```

### Q: 如何让 AI 只生成部分代码？

明确指定需要生成的内容：

```
只需要生成 Model 和 Repository 层，Service 和 Handler 稍后再做。
```

### Q: 生成代码后编译报错怎么办？

1. 确保运行了 `go mod tidy`
2. 确保运行了 Wire 命令
3. 让 AI 检查并修复错误：

```
运行时出现以下错误，请帮我修复：
[粘贴错误信息]
```

## 注意事项

1. **先阅读规则**：每次新对话都让 AI 先阅读 CLAUDE.md
2. **检查生成的代码**：AI 生成的代码需要人工审核
3. **运行 Wire**：添加新模块后必须重新生成 Wire 代码
4. **数据库迁移**：新模型需要在 bootstrap.go 中添加迁移
5. **测试验证**：生成代码后进行测试验证
