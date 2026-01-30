// Package router 路由配置
package router

import (
	"github.com/gin-gonic/gin"
	"ruleback/internal/handler"
	"ruleback/internal/middleware"
	"ruleback/internal/wire"
)

// Setup 初始化并配置路由（使用依赖注入）
func Setup(handlers *wire.Handlers) *gin.Engine {
	r := gin.New()

	registerGlobalMiddleware(r)
	registerHealthRoutes(r)
	registerAPIRoutes(r, handlers)

	return r
}

// SetupLegacy 初始化并配置路由（保留向后兼容）
func SetupLegacy() *gin.Engine {
	r := gin.New()

	registerGlobalMiddleware(r)
	registerHealthRoutes(r)
	registerAPIRoutesLegacy(r)

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

// registerAPIRoutes 注册API路由（使用依赖注入）
func registerAPIRoutes(r *gin.Engine, handlers *wire.Handlers) {
	v1 := r.Group("/api/v1")
	{
		registerPublicRoutes(v1)
		registerAuthenticatedRoutes(v1, handlers)
	}
}

// registerAPIRoutesLegacy 注册API路由（保留向后兼容）
func registerAPIRoutesLegacy(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		registerPublicRoutes(v1)
		registerAuthenticatedRoutesLegacy(v1)
	}
}

// registerPublicRoutes 注册公开路由（无需认证）
func registerPublicRoutes(rg *gin.RouterGroup) {
	// 认证相关路由
	// auth := rg.Group("/auth")
	// {
	//     authHandler := handler.GetAuthHandler()
	//     auth.POST("/login", authHandler.Login)
	//     auth.POST("/register", authHandler.Register)
	// }
}

// registerAuthenticatedRoutes 注册需要认证的路由（使用依赖注入）
func registerAuthenticatedRoutes(rg *gin.RouterGroup, handlers *wire.Handlers) {
	// 临时：不需要认证的用户路由（仅用于开发测试）
	registerUserRoutes(rg, handlers.UserHandler)
}

// registerAuthenticatedRoutesLegacy 注册需要认证的路由（保留向后兼容）
func registerAuthenticatedRoutesLegacy(rg *gin.RouterGroup) {
	registerUserRoutesLegacy(rg)
}

// registerUserRoutes 注册用户相关路由（使用依赖注入）
func registerUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	users := rg.Group("/users")
	{
		users.GET("", userHandler.List)
		users.POST("", userHandler.Create)
		users.GET("/:id", userHandler.GetByID)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
	}
}

// registerUserRoutesLegacy 注册用户相关路由（保留向后兼容）
func registerUserRoutesLegacy(rg *gin.RouterGroup) {
	userHandler := handler.GetUserHandler()

	users := rg.Group("/users")
	{
		users.GET("", userHandler.List)
		users.POST("", userHandler.Create)
		users.GET("/:id", userHandler.GetByID)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
	}
}
