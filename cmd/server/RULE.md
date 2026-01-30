# cmd/server 模块 AI 代码生成规则

> **模块职责**: 应用程序入口，负责初始化和启动

---

## 一、本模块的文件结构

```
cmd/server/
├── main.go       # 程序入口（简洁展示主流程）
├── bootstrap.go  # 初始化和服务器管理函数
└── RULE.md       # 本规则文件
```

---

## 二、文件职责划分

### main.go - 程序入口（保持简洁）

只包含常量、全局变量、main函数：

```go
package main

import (
    "fmt"
    "os"
    "time"

    "ruleback/internal/config"
)

const (
    ConfigPath      = "configs/config.yaml"
    ShutdownTimeout = 5 * time.Second
)

var cfg *config.Config

func main() {
    // 1. 初始化应用
    if err := initApp(); err != nil {
        fmt.Printf("应用初始化失败: %v\n", err)
        os.Exit(1)
    }

    // 2. 启动服务器
    srv := startServer()

    // 3. 等待关闭信号
    gracefulShutdown(srv)
}
```

### bootstrap.go - 初始化函数

包含所有初始化和服务器管理函数：
- `initApp()` - 初始化应用程序
- `initLogger()` - 初始化日志系统
- `initDatabase()` - 初始化数据库连接
- `migrateDatabase()` - 执行数据库迁移
- `startServer()` - 启动HTTP服务器
- `gracefulShutdown()` - 优雅关闭服务器

---

## 三、添加新初始化步骤

### 步骤1: 在bootstrap.go中添加初始化函数

```go
// initRedis 初始化Redis连接
func initRedis() error {
    if err := redis.Init(&cfg.Redis); err != nil {
        return err
    }
    logger.Info("Redis初始化完成")
    return nil
}
```

### 步骤2: 在initApp中调用

```go
func initApp() error {
    // ...现有初始化

    if err = initRedis(); err != nil {
        return logger.Errorf("初始化Redis失败: %w", err)
    }

    return nil
}
```

### 步骤3: 在gracefulShutdown中添加清理

```go
func gracefulShutdown(srv *http.Server) {
    // ...现有清理

    if err := redis.Close(); err != nil {
        logger.Error("Redis关闭异常", logger.Err(err))
    }
}
```

---

## 四、添加新数据库模型迁移

在bootstrap.go的migrateDatabase函数中添加：

```go
func migrateDatabase() error {
    models := []interface{}{
        &model.User{},
        &model.Order{},  // 添加新模型
    }
    return database.AutoMigrate(models...)
}
```

---

## 五、初始化顺序规范

必须按以下顺序初始化：

```
1. 配置 (config)
2. 日志 (logger)
3. 数据库 (database)
4. 缓存 (redis等，如需要)
5. 数据库迁移 (migrate)
```

关闭顺序与初始化相反。

---

## 六、禁止行为

| 禁止 | 正确做法 |
|------|----------|
| 在main.go中定义路由 | 在router模块定义 |
| 在main.go中写业务逻辑 | 业务逻辑放在service层 |
| 硬编码配置值 | 从config读取 |
| 使用装饰性分隔线注释 | 使用简洁单行注释 |
| 使用fmt.Println记录日志 | 使用logger包 |
| 在main.go中添加过多函数 | 函数放到bootstrap.go |
