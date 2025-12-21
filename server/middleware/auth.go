package middleware

import (
	"net/http"
	"play-tools/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// 优先从Cookie中获取token
		cookieToken, err := c.Cookie("token")
		if err == nil && cookieToken != "" {
			token = cookieToken
		} else {
			// 如果Cookie中没有，从Header中获取
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 401,
					"msg":  "未登录",
				})
				c.Abort()
				return
			}

			// 解析Bearer Token
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 401,
					"msg":  "Token格式错误",
				})
				c.Abort()
				return
			}
			token = parts[1]
		}

		// 验证Token
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Token无效",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
