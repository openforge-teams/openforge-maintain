package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/pkg/utils"
)

const (
	contextUserIDKey   = "auth_user_id"
	contextUsernameKey = "auth_username"
)

// Auth JWT 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code:    401,
				Message: "authorization header is required",
				Data:    nil,
			})
			c.Abort()
			return
		}

		// 检查 Bearer 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code:    401,
				Message: "invalid authorization header format, expected Bearer <token>",
				Data:    nil,
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.Response{
				Code:    401,
				Message: "invalid or expired token",
				Data:    nil,
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set(contextUserIDKey, claims.UserID)
		c.Set(contextUsernameKey, claims.Username)

		c.Next()
	}
}

// GetAuthUserID 从上下文中获取已认证的用户ID
func GetAuthUserID(c *gin.Context) uint {
	userIDStr, exists := c.Get(contextUserIDKey)
	if !exists {
		return 0
	}
	userIDStrStr, ok := userIDStr.(string)
	if !ok {
		return 0
	}

	var userID uint
	for _, ch := range userIDStrStr {
		userID = userID*10 + uint(ch-'0')
	}
	return userID
}

// GetAuthUsername 从上下文中获取已认证的用户名
func GetAuthUsername(c *gin.Context) string {
	username, exists := c.Get(contextUsernameKey)
	if !exists {
		return ""
	}
	usernameStr, ok := username.(string)
	if !ok {
		return ""
	}
	return usernameStr
}
