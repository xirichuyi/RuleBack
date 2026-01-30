# RuleBack

RuleBack 是一个面向 AI 代码生成优化的 Go 后端框架，采用清晰的分层架构和规范化的代码生成规则，使 AI 能够生成一致、高质量的代码。

## 特性

- **分层架构**: Router → Handler → Service → Repository → Model
- **依赖注入**: 使用 Google Wire 进行编译时依赖注入
- **统一响应**: 标准化的 API 响应格式
- **错误处理**: 分段式错误码和统一错误处理
- **结构化日志**: 基于 Zap 的高性能结构化日志
- **配置管理**: 基于 Viper 的配置管理，支持环境变量覆盖
- **AI 友好**: 每个模块都有 RULE.md 规则文件指导代码生成

## 技术栈

| 功能 | 库 |
|------|-----|
| Web框架 | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io/) |
| 日志 | [Zap](https://github.com/uber-go/zap) |
| 配置 | [Viper](https://github.com/spf13/viper) |
| 依赖注入 | [Wire](https://github.com/google/wire) |

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+ 或 PostgreSQL 14+

### 安装

```bash
# 克隆项目
git clone https://github.com/your-repo/ruleback.git
cd ruleback

# 安装依赖
go mod tidy

# 安装 Wire (用于依赖注入代码生成)
go install github.com/google/wire/cmd/wire@latest
```

### 配置

复制配置文件并修改：

```bash
cp configs/config.yaml.example configs/config.yaml
```

主要配置项：

```yaml
app:
  name: "ruleback"
  env: "development"    # development, staging, production

server:
  host: "0.0.0.0"
  port: 8080

database:
  driver: "mysql"       # mysql 或 postgres
  host: "localhost"
  port: 3306
  database: "ruleback"
  username: "root"
  password: ""
```

### 运行

```bash
# 开发模式运行
go run cmd/server/main.go cmd/server/bootstrap.go

# 编译并运行
go build -o ruleback ./cmd/server/...
./ruleback
```

### 环境变量覆盖

配置可通过环境变量覆盖，格式：`APP_分组_字段`（大写下划线）

```bash
export APP_SERVER_PORT=9090
export APP_DATABASE_HOST=192.168.1.100
export APP_JWT_SECRET=your-production-secret
```

## 项目结构

```
ruleback/
├── cmd/
│   └── server/
│       ├── main.go           # 程序入口
│       ├── bootstrap.go      # 初始化和服务器管理
│       └── RULE.md          # 入口模块规则
├── configs/
│   └── config.yaml          # 配置文件
├── docs/
│   ├── api.md               # API文档
│   └── RULE.md             # 文档规则
├── internal/
│   ├── config/              # 配置管理
│   │   ├── config.go
│   │   └── RULE.md
│   ├── model/               # 数据模型
│   │   ├── base.go
│   │   ├── user.go
│   │   └── RULE.md
│   ├── repository/          # 数据访问层
│   │   ├── base.go
│   │   ├── user_repository.go
│   │   ├── interfaces.go
│   │   └── RULE.md
│   ├── service/             # 业务逻辑层
│   │   ├── user_service.go
│   │   ├── interfaces.go
│   │   └── RULE.md
│   ├── handler/             # HTTP处理层
│   │   ├── user_handler.go
│   │   └── RULE.md
│   ├── middleware/          # 中间件
│   │   ├── middleware.go
│   │   └── RULE.md
│   ├── router/              # 路由配置
│   │   ├── router.go
│   │   └── RULE.md
│   └── wire/                # 依赖注入
│       ├── providers.go
│       ├── wire.go
│       ├── wire_gen.go
│       └── RULE.md
├── pkg/
│   ├── database/            # 数据库连接
│   │   └── database.go
│   ├── errors/              # 错误处理
│   │   ├── errors.go
│   │   └── RULE.md
│   ├── logger/              # 日志记录
│   │   ├── logger.go
│   │   └── RULE.md
│   └── response/            # 统一响应
│       ├── response.go
│       └── RULE.md
├── CLAUDE.md                # AI代码生成主规则
├── go.mod
└── README.md
```

## 架构说明

### 分层架构

```
请求 → Router → Middleware → Handler → Service → Repository → Database
                                ↓
响应 ← Router ← Middleware ← Handler ← Service ← Repository ← Database
```

| 层级 | 职责 |
|------|------|
| Router | 路由注册、中间件应用 |
| Middleware | 认证、日志、CORS、错误处理等横切关注点 |
| Handler | 接收请求、参数验证、调用Service、返回响应 |
| Service | 业务逻辑处理、错误转换 |
| Repository | 数据库CRUD操作 |
| Model | 数据结构定义 |

### 依赖注入

项目使用 Google Wire 进行编译时依赖注入：

```
database.GetDB()
       │
       ▼
BaseRepository → UserRepository → UserService → UserHandler → Handlers → Router
```

添加新模块后需重新生成 Wire 代码：

```bash
~/go/bin/wire ./internal/wire/...
```

## API 规范

### 响应格式

**成功响应：**
```json
{
    "code": 0,
    "message": "success",
    "data": {}
}
```

**分页响应：**
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

**错误响应：**
```json
{
    "code": 10001,
    "message": "参数错误"
}
```

### 错误码规范

| 范围 | 用途 |
|------|------|
| 0 | 成功 |
| 10000-10999 | 通用错误 |
| 20000-20999 | 用户模块 |
| 30000-30999 | 其他业务模块 |

## 开发指南

### 添加新模块

以添加订单模块为例：

1. **创建模型** - `internal/model/order.go`
2. **创建Repository** - `internal/repository/order_repository.go`
3. **创建Service** - `internal/service/order_service.go`
4. **创建Handler** - `internal/handler/order_handler.go`
5. **更新Wire** - `internal/wire/providers.go`
6. **注册路由** - `internal/router/router.go`
7. **数据库迁移** - `cmd/server/bootstrap.go`
8. **更新API文档** - `docs/api.md`

详细规则请参考各模块的 `RULE.md` 文件。

### 代码规范

- 使用 `response` 包返回响应，禁止直接使用 `c.JSON()`
- 使用 `errors` 包的错误码，禁止硬编码数字
- 使用 `logger` 包记录日志，禁止使用 `fmt.Println`
- 从 `config` 读取配置，禁止硬编码配置值
- 注释使用中文，保持简洁
- 禁止使用装饰性分隔线注释（如 `// ===`）

### AI 代码生成

项目包含完整的 AI 代码生成规则：

- `CLAUDE.md` - 主规则文件，定义全局行为准则
- 各模块 `RULE.md` - 模块级规则，定义具体实现规范

AI 在生成代码前应阅读相关规则文件。

## 常用命令

```bash
# 运行项目
go run cmd/server/main.go cmd/server/bootstrap.go

# 编译项目
go build -o ruleback ./cmd/server/...

# 运行测试
go test ./...

# 格式化代码
go fmt ./...

# 检查代码
go vet ./...

# 更新依赖
go mod tidy

# 生成Wire代码
~/go/bin/wire ./internal/wire/...
```

## 许可证

MIT License
