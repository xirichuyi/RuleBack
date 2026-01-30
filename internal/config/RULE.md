# internal/config 模块 AI 代码生成规则

> **模块职责**: 应用程序配置管理，从配置文件和环境变量加载配置

---

## 一、本模块的文件

```
internal/config/
├── config.go   # 配置定义和加载
└── RULE.md     # 本规则文件
```

---

## 二、添加新配置项的完整流程

### 场景A: 在现有配置分组中添加字段

```go
// ServerConfig HTTP服务器配置
type ServerConfig struct {
    // ... 现有字段 ...
    MaxHeaderBytes int `mapstructure:"max_header_bytes"`
}

// 在 setDefaults 中设置默认值
func setDefaults(cfg *Config) {
    if cfg.Server.MaxHeaderBytes == 0 {
        cfg.Server.MaxHeaderBytes = 1 << 20 // 1MB
    }
}
```

### 场景B: 添加新的配置分组

```go
// RedisConfig Redis配置
type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
    PoolSize int    `mapstructure:"pool_size"`
}

// 在 Config 根结构体中添加
type Config struct {
    // ... 现有分组 ...
    Redis RedisConfig `mapstructure:"redis"`
}

// 添加辅助方法
func (c *RedisConfig) GetAddress() string {
    return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
```

---

## 三、命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 配置结构体 | PascalCase + Config | `DatabaseConfig` |
| 字段名 | PascalCase | `MaxOpenConns` |
| mapstructure标签 | snake_case | `max_open_conns` |
| 环境变量 | APP_分组_字段 | `APP_DATABASE_HOST` |

---

## 四、辅助方法规范

| 方法类型 | 命名规范 | 示例 |
|---------|---------|------|
| 获取值 | Get + 描述 | `GetDSN()`, `GetAddress()` |
| 判断 | Is + 描述 | `IsDevelopment()`, `IsProduction()` |

---

## 五、禁止行为

| 禁止 | 原因 |
|------|------|
| 在配置模块中引入业务逻辑 | 配置模块只负责配置管理 |
| 硬编码配置值 | 应通过配置文件或默认值 |
| 修改已存在字段的类型 | 破坏兼容性 |
| 使用装饰性分隔线注释 | 使用简洁单行注释 |

---

## 六、已存在的配置分组

| 分组 | 用途 |
|------|------|
| `App` | 应用基础配置 (Name, Env, Debug) |
| `Server` | HTTP服务器配置 (Host, Port, Timeout) |
| `Database` | 数据库配置 (Driver, Host, 连接池) |
| `Log` | 日志配置 (Level, Format, Output) |
| `JWT` | JWT认证配置 (Secret, ExpireTime) |

---

## 七、环境变量覆盖

```bash
# 格式: APP_分组_字段（全大写，下划线分隔）
APP_SERVER_PORT=9090
APP_DATABASE_HOST=192.168.1.100
APP_LOG_LEVEL=debug
APP_JWT_SECRET=my-secret-key
```
