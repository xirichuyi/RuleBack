// Package middleware HTTP中间件
package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"ruleback/pkg/errors"
	"ruleback/pkg/logger"
	"ruleback/pkg/response"
)

// Logger 请求日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		if query != "" {
			path = path + "?" + query
		}

		logger.Info("HTTP请求",
			logger.String("method", method),
			logger.String("path", path),
			logger.String("ip", clientIP),
			logger.Int("status", statusCode),
			logger.Float64("latency_ms", float64(latency.Milliseconds())),
		)
	}
}

// Recovery 错误恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("服务器内部错误",
					logger.Field("error", err),
					logger.String("path", c.Request.URL.Path),
					logger.String("method", c.Request.Method),
				)
				response.InternalServerError(c, "服务器内部错误")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORS 跨域处理中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Auth JWT认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Unauthorized(c, "未提供认证Token")
			c.Abort()
			return
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// TODO: 验证JWT Token
		// claims, err := jwt.ParseToken(token)
		// if err != nil {
		//     response.Unauthorized(c, "Token无效或已过期")
		//     c.Abort()
		//     return
		// }
		// c.Set("user_id", claims.UserID)

		c.Next()
	}
}

// RequireRole 角色权限中间件
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 从上下文获取用户角色并验证
		c.Next()
	}
}

// RequestID 请求ID中间件
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// TODO: 使用UUID库生成
	return time.Now().Format("20060102150405.000000")
}

// RateLimit 请求限流中间件
func RateLimit(limit int, window int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现基于Redis或内存的限流
		c.Next()
	}
}

// Timeout 请求超时中间件
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现请求超时处理
		c.Next()
	}
}

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			appErr := errors.GetAppError(err)
			if appErr != nil {
				response.Fail(c, appErr.Code, appErr.Message)
				return
			}

			logger.Error("未处理的错误", logger.Err(err))
			response.InternalServerError(c, "服务器内部错误")
		}
	}
}
