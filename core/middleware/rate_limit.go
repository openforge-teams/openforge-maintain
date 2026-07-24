package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 基于内存的速率限制器
type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*clientBucket
	rate     int           // 时间窗口内允许的最大请求数
	window   time.Duration // 时间窗口
}

// clientBucket 客户端请求桶
type clientBucket struct {
	count    int
	expireAt time.Time
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*clientBucket),
		rate:    rate,
		window:  window,
	}

	// 启动清理过期桶的协程
	go rl.cleanup()

	return rl
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	bucket, exists := rl.clients[key]
	if !exists || now.After(bucket.expireAt) {
		rl.clients[key] = &clientBucket{
			count:    1,
			expireAt: now.Add(rl.window),
		}
		return true
	}

	if bucket.count >= rl.rate {
		return false
	}

	bucket.count++
	return true
}

// cleanup 定期清理过期记录
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, bucket := range rl.clients {
			if now.After(bucket.expireAt) {
				delete(rl.clients, key)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit 速率限制中间件
func RateLimit(rate int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, window)

	return func(c *gin.Context) {
		// 使用客户端 IP 作为限制 key
		key := c.ClientIP()

		if !limiter.Allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "too many requests, please try again later",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
