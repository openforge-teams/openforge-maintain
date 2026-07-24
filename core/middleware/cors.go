package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS 跨域资源共享中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 允许所有来源（生产环境应限制为具体域名）
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Trace-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Trace-ID")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // 24 小时预检缓存

		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// IsCORSPreflight 判断是否为 CORS 预检请求
func IsCORSPreflight(c *gin.Context) bool {
	return c.Request.Method == http.MethodOptions &&
		strings.Contains(c.GetHeader("Access-Control-Request-Method"), "GET") ||
		c.GetHeader("Access-Control-Request-Method") != ""
}
