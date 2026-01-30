# 示例代码

本目录包含使用 RuleBack 框架创建业务模块的完整示例代码。

## 文件说明

| 文件 | 说明 |
|------|------|
| `user_model.go.example` | User 模型定义示例 |
| `user_repository.go.example` | User Repository 实现示例 |
| `repository_interfaces.go.example` | Repository 接口定义示例 |
| `user_service.go.example` | User Service 实现示例 |
| `service_interfaces.go.example` | Service 接口定义示例 |
| `user_handler.go.example` | User Handler 实现示例 |

## 使用方法

1. **复制示例文件到对应目录**

```bash
# 复制并重命名（以 Product 模块为例）
cp examples/user_model.go.example internal/model/product.go
cp examples/user_repository.go.example internal/repository/product_repository.go
cp examples/user_service.go.example internal/service/product_service.go
cp examples/user_handler.go.example internal/handler/product_handler.go
```

2. **修改文件内容**

将所有 `User` 替换为你的模型名（如 `Product`），并根据业务需求调整字段和方法。

3. **更新 Wire 配置**

在 `internal/wire/providers.go` 中添加对应的 Provider 函数。

4. **更新路由配置**

在 `internal/router/router.go` 中注册新的路由。

5. **重新生成 Wire 代码**

```bash
~/go/bin/wire ./internal/wire/...
```

## 推荐方式

使用 AI（如 Claude）配合框架的 RULE.md 规则文件自动生成代码，而不是手动复制和修改示例文件。

详见 [docs/AI_USAGE.md](../docs/AI_USAGE.md)。
