package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 返回跨域中间件（基于社区标准库 gin-contrib/cors）。
// 配置保持与历史行为一致：允许任意来源、不携带凭证、预检缓存 86400 秒。
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "UPDATE"},
		AllowHeaders:     []string{"Token", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "X-Max"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour, // 86400 秒
	})
}
