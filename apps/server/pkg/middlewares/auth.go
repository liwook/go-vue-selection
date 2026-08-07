package middlewares

import (
	"strings"

	"github.com/liwook/go-vue-selection/pkg/jwt"
	"github.com/liwook/go-vue-selection/pkg/result"

	"github.com/gin-gonic/gin"
)

// CtxUserIDKey 是认证中间件写入 gin.Context 的当前用户ID键，
// handler 层通过 middlewares.CtxUserIDKey 读取，避免与业务层产生反向依赖。
const CtxUserIDKey = "userID"

// JWTAuthMiddleware 基于JWT的认证中间件
// 客户端通过标准 Authorization 头携带令牌：Authorization: Bearer <token>
func JWTAuthMiddleware(j *jwt.JWT) func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			result.Error(c, result.CodeNeedLogin)
			c.Abort()
			return
		}

		// 按空格分割，期望格式：Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			result.Error(c, result.CodeInvalidToken)
			c.Abort()
			return
		}

		// parts[1] 即为 token 字符串，使用项目内 jwt.ParseToken 解析
		mc, err := j.ParseToken(parts[1])
		if err != nil {
			result.Error(c, result.CodeInvalidToken)
			c.Abort()
			return
		}

		// 将当前请求的 userID 信息保存到上下文 c 上
		c.Set(CtxUserIDKey, mc.UserID)
		c.Next() // 后续处理函数可通过 c.Get(CtxUserIDKey) 获取当前用户信息
	}
}
