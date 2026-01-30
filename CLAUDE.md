# RuleBack - AI代码生成规范

> **重要**: 这是AI生成代码的主规则文件。在生成任何代码之前，必须完整阅读本文件。

---

## 第一章：核心行为规范

### 1.1 AI必须遵守的行为准则

| 规则 | 说明 |
|------|------|
| 不生成表情符号 | 代码、注释、文档中禁止出现任何emoji |
| 不自动生成文档 | 完成代码后不要自动创建README或说明文档 |
| 不使用硬编码 | 关键配置必须写入配置文件，禁止在代码中硬编码 |
| 不生成测试用例 | 除非用户明确要求，否则不生成测试代码 |
| 更新API文档 | 每次修改接口后，更新docs目录下的API文档（不创建新文件） |
| 注释使用中文 | 默认使用中文编写注释，保持简洁明了 |
| 注释不要过多 | 注释简洁扼要，不使用装饰性分隔线 |
| 复用现有实例 | 不重复创建实例，复用已存在的单例 |
| 使用成熟框架 | 使用Gin、GORM等成熟框架，禁止造轮子 |

### 1.2 实例复用原则（重要）

**禁止重复创建实例，必须复用现有实例：**

```go
// 禁止：每次都创建新实例
func NewUserHandler() *UserHandler {
    return &UserHandler{
        service: service.NewUserService(),  // 每次都new，禁止
    }
}

// 正确：使用单例模式
var (
    userServiceInstance *UserService
    userServiceOnce     sync.Once
)

func GetUserService() *UserService {
    userServiceOnce.Do(func() {
        userServiceInstance = &UserService{
            repo: GetUserRepository(),
        }
    })
    return userServiceInstance
}
```

**全局实例必须使用sync.Once保护：**

```go
var (
    globalDB *gorm.DB
    dbOnce   sync.Once
)

func GetDB() *gorm.DB {
    dbOnce.Do(func() {
        // 初始化
    })
    return globalDB
}
```

### 1.3 禁止造轮子

**必须使用成熟的第三方库：**

| 功能 | 使用的库 | 禁止行为 |
|------|---------|---------|
| Web框架 | `github.com/gin-gonic/gin` | 自己实现HTTP路由 |
| ORM | `gorm.io/gorm` | 自己拼接SQL |
| 日志 | `go.uber.org/zap` | 使用fmt.Println |
| 配置 | `github.com/spf13/viper` | 自己解析配置文件 |
| UUID | `github.com/google/uuid` | 自己拼接时间戳生成ID |
| 密码加密 | `golang.org/x/crypto/bcrypt` | 自己实现加密算法 |
| JWT | `github.com/golang-jwt/jwt` | 自己实现token |
| 依赖注入 | `github.com/google/wire` | 手动管理复杂依赖 |

### 1.4 注释规范

**禁止使用装饰性分隔线注释：**
```go
// 禁止这样的注释
// =============================================================================
// 初始化函数
// =============================================================================

// 正确的注释方式
// UserService 用户业务逻辑层
type UserService struct {
    repo *repository.UserRepository
}
```

### 1.5 运行模式

项目支持两种运行模式：

| 模式 | 配置值 | 特点 |
|------|--------|------|
| 开发模式 | `app.env: development` | 详细日志、debug信息 |
| 生产模式 | `app.env: production` | 精简日志、性能优化 |

### 1.6 配置管理原则

**关键配置必须写入配置文件，禁止硬编码：**

```go
// 禁止硬编码
const JWTSecret = "my-secret-key"  // 禁止

// 正确方式
secret := config.Get().JWT.Secret
```

---

## 第二章：需求处理流程

### 2.1 接到需求后的处理步骤

```
步骤1: 理解需求 → 分析用户需求，明确要实现的功能
步骤2: 拆分任务 → 将需求拆分为具体的实现步骤
步骤3: 检查复用 → 检查是否有可复用的现有代码和实例
步骤4: 按模块实现 → Model → Repository → Service → Handler → Router
步骤5: 更新文档 → 在docs目录更新API文档
```

### 2.2 任务拆分示例

**需求**: 添加订单管理功能

**拆分后的任务：**
1. 在 `internal/model/order.go` 创建Order模型
2. 在 `internal/repository/order_repository.go` 创建数据访问层
3. 在 `internal/service/order_service.go` 创建业务逻辑层
4. 在 `internal/handler/order_handler.go` 创建HTTP处理器
5. 在 `internal/router/router.go` 注册路由
6. 在 `cmd/server/main.go` 添加数据库迁移
7. 更新 `docs/api.md` 的订单接口文档

---

## 第三章：绝对禁止的行为

| 禁止行为 | 正确做法 |
|---------|---------|
| 使用 `fmt.Println` 记录日志 | 使用 `logger.Info/Debug/Error` |
| 使用 `c.JSON` 返回响应 | 使用 `response.SuccessWithData/Fail` |
| 硬编码错误码数字 | 使用 `errors.CodeInvalidParams` |
| 硬编码配置值 | 从 `config.Get()` 读取 |
| 在Handler中写业务逻辑 | 业务逻辑放在Service层 |
| 在Service中直接操作数据库 | 通过Repository操作数据库 |
| 每次都New新实例 | 使用单例模式复用实例 |
| 自己实现已有功能 | 使用成熟的第三方库 |
| 自己拼接时间戳做ID | 使用uuid库 |
| 明文存储密码 | 使用bcrypt加密 |
| 使用装饰性分隔线注释 | 使用简洁的单行注释 |
| 自动生成测试用例 | 除非用户明确要求 |
| 自动创建文档文件 | 只更新现有文档 |

---

## 第四章：项目架构

### 4.1 目录结构

```
ruleback/
├── cmd/server/main.go          # 程序入口
├── configs/config.yaml         # 配置文件
├── docs/api.md                 # API文档
├── internal/
│   ├── config/                # 配置管理
│   ├── model/                 # 数据模型
│   ├── repository/            # 数据访问
│   ├── service/               # 业务逻辑
│   ├── handler/               # HTTP处理
│   ├── middleware/            # 中间件
│   ├── router/                # 路由配置
│   └── wire/                  # 依赖注入（Wire）
└── pkg/                       # 公共包
    ├── response/              # 统一响应
    ├── errors/                # 错误处理
    ├── logger/                # 日志记录
    └── database/              # 数据库连接
```

### 4.2 代码分层

```
Router → Handler → Service → Repository → Model
```

每层职责：
- **Router**: 路由注册、中间件应用
- **Handler**: 接收请求、参数验证、调用Service、返回响应
- **Service**: 业务逻辑处理、错误转换
- **Repository**: 数据库CRUD操作
- **Model**: 数据结构定义

---

## 第五章：实例管理规范

### 5.1 Wire依赖注入（推荐）

项目使用Google Wire进行编译时依赖注入，在 `internal/wire/` 目录管理：

```go
// 1. 每个模块提供New构造函数
func NewUserRepository(base *BaseRepository) *UserRepository {
    return &UserRepository{BaseRepository: base}
}

func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

func NewUserHandler(svc *service.UserService) *UserHandler {
    return &UserHandler{service: svc}
}

// 2. 在wire/providers.go定义Provider
func ProvideUserRepository(base *repository.BaseRepository) *repository.UserRepository {
    return repository.NewUserRepository(base)
}

// 3. 在bootstrap.go使用Wire初始化
handlers, _ := wire.InitializeHandlers(database.GetDB())
r := router.Setup(handlers)
```

**添加新模块后必须重新生成Wire代码：**
```bash
~/go/bin/wire ./internal/wire/...
```

### 5.2 单例模式（向后兼容）

保留Get*方法用于向后兼容和简单场景：

```go
var (
    userRepoInstance *UserRepository
    userRepoOnce     sync.Once
)

func GetUserRepository() *UserRepository {
    userRepoOnce.Do(func() {
        userRepoInstance = &UserRepository{
            BaseRepository: GetBaseRepository(),
        }
    })
    return userRepoInstance
}
```

### 5.3 在路由中使用

```go
// 推荐：使用Wire注入的Handler
func registerUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
    users := rg.Group("/users")
    {
        users.GET("", userHandler.List)
        users.POST("", userHandler.Create)
    }
}

// 向后兼容：使用Get单例
func registerUserRoutesLegacy(rg *gin.RouterGroup) {
    h := handler.GetUserHandler()
    // ...
}
```

---

## 第六章：已存在的功能（直接使用）

### 6.1 响应函数 (pkg/response)

```go
response.Success(c)
response.SuccessWithData(c, data)
response.SuccessWithPage(c, list, total, page, size)
response.Fail(c, errors.CodeXxx, "错误消息")
```

### 6.2 错误码 (pkg/errors)

```go
errors.CodeSuccess         // 0
errors.CodeInvalidParams   // 10001
errors.CodeUnauthorized    // 10002
errors.CodeNotFound        // 10004
errors.CodeDatabaseError   // 10007
```

### 6.3 日志函数 (pkg/logger)

```go
logger.Debug("消息", logger.String("key", "value"))
logger.Info("消息", logger.Uint("user_id", 1))
logger.Error("消息", logger.Err(err))
```

### 6.4 配置获取 (internal/config)

```go
cfg := config.Get()
cfg.App.Env
cfg.Server.Port
cfg.Database.Host
cfg.JWT.Secret
```

---

## 第七章：第三方库使用规范

### 7.1 必须使用的库

| 用途 | 库 | 导入路径 |
|------|-----|---------|
| Web框架 | Gin | `github.com/gin-gonic/gin` |
| ORM | GORM | `gorm.io/gorm` |
| 日志 | Zap | `go.uber.org/zap` |
| 配置 | Viper | `github.com/spf13/viper` |
| UUID | Google UUID | `github.com/google/uuid` |
| 密码 | Bcrypt | `golang.org/x/crypto/bcrypt` |
| 依赖注入 | Wire | `github.com/google/wire` |

### 7.2 使用示例

```go
// UUID生成
import "github.com/google/uuid"
id := uuid.New().String()

// 密码加密
import "golang.org/x/crypto/bcrypt"
hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// 密码验证
err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
```

---

## 第八章：规则文件索引

| 模块 | 规则文件 |
|------|---------|
| 配置 | `internal/config/RULE.md` |
| 模型 | `internal/model/RULE.md` |
| Repository | `internal/repository/RULE.md` |
| Service | `internal/service/RULE.md` |
| Handler | `internal/handler/RULE.md` |
| 中间件 | `internal/middleware/RULE.md` |
| 路由 | `internal/router/RULE.md` |
| Wire依赖注入 | `internal/wire/RULE.md` |
| 响应 | `pkg/response/RULE.md` |
| 错误 | `pkg/errors/RULE.md` |
| 日志 | `pkg/logger/RULE.md` |
| 文档 | `docs/RULE.md` |

---

## 第九章：检查清单

### 生成代码前

- [ ] 已理解用户需求
- [ ] 已拆分为具体步骤
- [ ] 已检查是否有可复用的代码
- [ ] 已检查是否有成熟的库可用

### 生成代码时

- [ ] 使用单例模式管理实例（Get而不是New）
- [ ] 使用成熟的第三方库
- [ ] 使用response包返回响应
- [ ] 使用errors包的错误码
- [ ] 使用logger包记录日志
- [ ] 从config读取配置
- [ ] 注释简洁，使用中文
- [ ] 不使用装饰性分隔线

### 生成代码后

- [ ] 更新docs/api.md
- [ ] 不自动创建新文档
- [ ] 不自动生成测试用例

---

**记住：复用、简洁、规范。**
