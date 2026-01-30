// Package router 路由配置
package router

import (
	"github.com/gin-gonic/gin"
	"ruleback/internal/middleware"
	"ruleback/internal/wire"
)

// RouteRegister 路由注册函数类型
type RouteRegister func(rg *gin.RouterGroup, handlers *wire.Handlers)

// Setup 初始化并配置路由
// handlers: Wire注入的Handler实例
// customRoutes: 自定义路由注册函数（可选）
func Setup(handlers *wire.Handlers, customRoutes ...RouteRegister) *gin.Engine {
	r := gin.New()

	registerGlobalMiddleware(r)
	registerHealthRoutes(r)
	registerAPIRoutes(r, handlers, customRoutes...)

	return r
}

// registerGlobalMiddleware 注册全局中间件
func registerGlobalMiddleware(r *gin.Engine) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
}

// registerHealthRoutes 注册健康检查路由
func registerHealthRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}

// registerAPIRoutes 注册API路由
func registerAPIRoutes(r *gin.Engine, handlers *wire.Handlers, customRoutes ...RouteRegister) {
	v1 := r.Group("/api/v1")
	{
		// 注册公开路由（无需认证）
		registerPublicRoutes(v1, handlers)

		// 注册自定义路由
		for _, register := range customRoutes {
			register(v1, handlers)
		}
	}
}

// registerPublicRoutes 注册公开路由（无需认证）
// 使用框架时，可以在此处添加不需要认证的路由
// 示例:
//
//	auth := rg.Group("/auth")
//	{
//	    auth.POST("/login", handlers.AuthHandler.Login)
//	    auth.POST("/register", handlers.AuthHandler.Register)
//	}
func registerPublicRoutes(rg *gin.RouterGroup, handlers *wire.Handlers) {
	// 在此添加公开路由
}

// RegisterAuthenticatedRoutes 返回一个带认证中间件的路由注册函数
// 使用示例:
//
//	router.Setup(handlers, router.RegisterAuthenticatedRoutes(func(rg *gin.RouterGroup, h *wire.Handlers) {
//	    users := rg.Group("/users")
//	    users.GET("", h.UserHandler.List)
//	}))
func RegisterAuthenticatedRoutes(register RouteRegister) RouteRegister {
	return func(rg *gin.RouterGroup, handlers *wire.Handlers) {
		authenticated := rg.Group("")
		authenticated.Use(middleware.Auth())
		register(authenticated, handlers)
	}
}
