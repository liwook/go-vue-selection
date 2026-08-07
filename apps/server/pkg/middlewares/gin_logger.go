package middlewares

import (
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
)

// GinLogger 请求日志中间件：记录每个请求的方法、路径、状态码、耗时等。
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)
		slog.Info(path,
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"ip", c.ClientIP(),
			"user-agent", c.Request.UserAgent(),
			"errors", c.Errors.ByType(gin.ErrorTypePrivate).String(),
			"cost", cost,
		)
	}
}
