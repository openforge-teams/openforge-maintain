package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	securityEntryContextKey = "security_entry_matched"
)

// SecurityEntry 安全入口中间件
// 从环境变量 MAINTAIN_SECURITY_ENTRY 读取安全入口路径
// 如果设置了安全入口，只有包含该路径的请求才会被放行
// 安全入口主要用于隐藏管理openforge-maintain的真实路径
func SecurityEntry() gin.HandlerFunc {
	securityEntry := os.Getenv("MAINTAIN_SECURITY_ENTRY")

	// 如果未设置安全入口，放行所有请求
	if securityEntry == "" {
		return func(c *gin.Context) {
			c.Set(securityEntryContextKey, true)
			c.Next()
		}
	}

	// 清理安全入口路径，确保不以 / 开头和结尾
	securityEntry = strings.Trim(securityEntry, "/")

	return func(c *gin.Context) {
		requestPath := strings.Trim(c.Request.URL.Path, "/")

		// 检查请求路径是否包含安全入口
		if !strings.Contains(requestPath, securityEntry) {
			// 不包含安全入口路径，返回 404（不暴露真实存在）
			c.AbortWithStatus(404)
			return
		}

		c.Set(securityEntryContextKey, true)
		c.Next()
	}
}
