package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/core/model"
	"github.com/openforge-maintain/openforge-maintain/core/repository"
)

// Audit 审计日志中间件
func Audit(auditRepo repository.AuditRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 请求完成后记录审计日志
		duration := time.Since(startTime)
		userIDStr, _ := c.Get(contextUserIDKey)
		usernameStr, _ := c.Get(contextUsernameKey)

		var userID uint
		if userIDStr != nil {
			if idStr, ok := userIDStr.(string); ok {
				id, _ := strconv.ParseUint(idStr, 10, 32)
				userID = uint(id)
			}
		}

		username := ""
		if usernameStr != nil {
			if u, ok := usernameStr.(string); ok {
				username = u
			}
		}

		auditLog := &model.AuditLog{
			UserID:   userID,
			Action:   c.Request.Method + " " + c.Request.URL.Path,
			Resource: c.Request.URL.Path,
			Detail:   "user: " + username + " | status: " + strconv.Itoa(c.Writer.Status()) + " | duration: " + duration.String(),
			IP:       c.ClientIP(),
			CreatedAt: time.Now(),
		}

		// 异步写入审计日志，不影响请求响应
		go func() {
			_ = auditRepo.Create(auditLog)
		}()
	}
}
