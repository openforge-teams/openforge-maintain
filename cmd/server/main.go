package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ==================== Models ====================

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	Email     string    `json:"email" gorm:"size:120"`
	RoleID    uint      `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ==================== Response ====================

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceID string      `json:"trace_id"`
}

type PageResponse struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

func success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{Code: 0, Message: "success", Data: data, TraceID: uuid.New().String()})
}

func successPage(c *gin.Context, list interface{}, total int64, page, size int) {
	c.JSON(200, Response{Code: 0, Message: "success", Data: PageResponse{List: list, Total: total, Page: page, Size: size}, TraceID: uuid.New().String()})
}

func errResp(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{Code: code, Message: msg, Data: nil, TraceID: uuid.New().String()})
}

// ==================== JWT ====================

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("openforge-maintain-secret-key")

func generateToken(userID, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "openforge-maintain",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func parseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// ==================== Auth Middleware ====================

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) < 8 || token[:7] != "Bearer " {
			errResp(c, 401, "unauthorized")
			c.Abort()
			return
		}
		claims, err := parseToken(token[7:])
		if err != nil {
			errResp(c, 401, "invalid token")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// ==================== Main ====================

func main() {
	// Database
	dbPath := "/tmp/openforge-maintain.db"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	db.AutoMigrate(&User{})

	// Default admin
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("maintain@2024"), bcrypt.DefaultCost)
		db.Create(&User{Username: "admin", Password: string(hash), Email: "admin@openforge-maintain.local", RoleID: 1})
		log.Println("Created default admin user (admin / maintain@2024)")
	}

	// Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Trace-ID")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Serve frontend static files
	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		// Try common locations
		for _, d := range []string{"frontend/dist", "/opt/openforge-maintain/frontend"} {
			if _, err := os.Stat(d); err == nil {
				frontendDir = d
				break
			}
		}
	}
	if frontendDir != "" {
		r.Static("/assets", filepath.Join(frontendDir, "assets"))
		r.StaticFile("/favicon.ico", filepath.Join(frontendDir, "favicon.ico"))
		r.NoRoute(func(c *gin.Context) {
			c.File(filepath.Join(frontendDir, "index.html"))
		})
		log.Printf("Serving frontend from: %s", frontendDir)
	}

	// ==================== Core API Routes ====================
	core := r.Group("/api/v2/core")
	{
		// Auth
		core.POST("/auth/login", func(c *gin.Context) {
			var req struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			c.ShouldBindJSON(&req)
			var user User
			if db.Where("username = ?", req.Username).First(&user).Error != nil {
				errResp(c, 401, "invalid credentials")
				return
			}
			if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
				errResp(c, 401, "invalid credentials")
				return
			}
			token, _ := generateToken(strconv.Itoa(int(user.ID)), user.Username)
			success(c, gin.H{
				"access_token":  token,
				"refresh_token": token,
				"user": gin.H{
					"id":       user.ID,
					"username": user.Username,
					"email":    user.Email,
					"role_id":  user.RoleID,
				},
			})
		})

		core.POST("/auth/logout", authMiddleware(), func(c *gin.Context) {
			success(c, nil)
		})

		core.GET("/auth/profile", authMiddleware(), func(c *gin.Context) {
			userID := c.GetString("user_id")
			var user User
			if db.First(&user, userID).Error != nil {
				errResp(c, 404, "user not found")
				return
			}
			success(c, gin.H{"id": user.ID, "username": user.Username, "email": user.Email, "role_id": user.RoleID, "created_at": user.CreatedAt})
		})
	}

	// ==================== Agent API Routes ====================
	agent := r.Group("/api/v2")
	agent.Use(authMiddleware())
	{
		// Metrics
		agent.GET("/metrics/overview", func(c *gin.Context) {
			success(c, gin.H{
				"cpu_usage":         23.5,
				"memory_usage":      61.2,
				"memory_total":      15.6,
				"memory_used":       9.55,
				"disk_usage":        45.8,
				"disk_total":        256.0,
				"disk_used":         117.2,
				"container_count":   24,
				"container_running": 12,
				"network_rx":        128.5,
				"network_tx":        42.3,
				"uptime":            259200,
				"load_avg":          []float64{0.52, 0.38, 0.29},
				"hostname":          "openforge-server",
				"os":                "Ubuntu 22.04 LTS",
			})
		})

		agent.GET("/metrics/cpu", func(c *gin.Context) {
			now := time.Now().Unix()
			points := make([]gin.H, 0)
			for i := 0; i < 60; i++ {
				points = append(points, gin.H{
					"time":  now - int64(60-i),
					"value": 18 + float64(i%20)*0.5,
				})
			}
			success(c, gin.H{"points": points, "cores": 4, "usage": 23.5})
		})

		agent.GET("/metrics/memory", func(c *gin.Context) {
			now := time.Now().Unix()
			points := make([]gin.H, 0)
			for i := 0; i < 60; i++ {
				points = append(points, gin.H{
					"time":  now - int64(60-i),
					"value": 58 + float64(i%15)*0.4,
				})
			}
			success(c, gin.H{"points": points, "total": 15.6, "used": 9.55, "cached": 3.2, "free": 2.85})
		})

		agent.GET("/metrics/disk", func(c *gin.Context) {
			success(c, []gin.H{
				{"device": "/dev/sda1", "mount": "/", "total": 256.0, "used": 117.2, "free": 138.8, "usage": 45.8},
				{"device": "/dev/sda2", "mount": "/home", "total": 512.0, "used": 89.4, "free": 422.6, "usage": 17.5},
				{"device": "/dev/sdb1", "mount": "/data", "total": 1024.0, "used": 623.8, "free": 400.2, "usage": 60.9},
			})
		})

		agent.GET("/metrics/network", func(c *gin.Context) {
			now := time.Now().Unix()
			rx := make([]gin.H, 0)
			tx := make([]gin.H, 0)
			for i := 0; i < 60; i++ {
				rx = append(rx, gin.H{"time": now - int64(60-i), "value": float64(80+i%40) * 1.5})
				tx = append(tx, gin.H{"time": now - int64(60-i), "value": float64(30+i%20) * 1.2})
			}
			success(c, gin.H{"rx": rx, "tx": tx, "total_rx": 128.5, "total_tx": 42.3})
		})

		agent.GET("/metrics/processes", func(c *gin.Context) {
			procs := []gin.H{
				{"pid": 1, "name": "systemd", "user": "root", "cpu": 0.0, "memory": 1.2, "status": "running"},
				{"pid": 1023, "name": "nginx", "user": "www-data", "cpu": 0.5, "memory": 2.8, "status": "running"},
				{"pid": 2045, "name": "docker", "user": "root", "cpu": 1.2, "memory": 15.3, "status": "running"},
				{"pid": 3098, "name": "mysql", "user": "mysql", "cpu": 2.1, "memory": 28.7, "status": "running"},
				{"pid": 4521, "name": "redis-server", "user": "redis", "cpu": 0.3, "memory": 4.5, "status": "running"},
				{"pid": 5678, "name": "node", "user": "deploy", "cpu": 3.8, "memory": 12.1, "status": "running"},
				{"pid": 7201, "name": "sshd", "user": "root", "cpu": 0.1, "memory": 0.8, "status": "running"},
				{"pid": 8934, "name": "openforge-core", "user": "root", "cpu": 0.8, "memory": 5.2, "status": "running"},
				{"pid": 8935, "name": "openforge-agent", "user": "root", "cpu": 1.5, "memory": 8.9, "status": "running"},
			}
			successPage(c, procs, 45, 1, 20)
		})

		// Containers
		agent.GET("/containers", func(c *gin.Context) {
			containers := []gin.H{
				{"id": "a1b2c3d4", "name": "nginx-proxy", "image": "nginx:1.25", "state": "running", "status": "Up 3 days", "ports": []gin.H{{"ip": "0.0.0.0", "private_port": 80, "public_port": 80, "type": "tcp"}, {"ip": "0.0.0.0", "private_port": 443, "public_port": 443, "type": "tcp"}}, "created_at": "2024-07-21T10:30:00Z"},
				{"id": "e5f6g7h8", "name": "mysql-8", "image": "mysql:8.0", "state": "running", "status": "Up 3 days", "ports": []gin.H{{"ip": "0.0.0.0", "private_port": 3306, "public_port": 3306, "type": "tcp"}}, "created_at": "2024-07-21T10:30:00Z"},
				{"id": "i9j0k1l2", "name": "redis-7", "image": "redis:7-alpine", "state": "running", "status": "Up 3 days", "ports": []gin.H{{"ip": "0.0.0.0", "private_port": 6379, "public_port": 6379, "type": "tcp"}}, "created_at": "2024-07-21T10:30:00Z"},
				{"id": "m3n4o5p6", "name": "wordpress", "image": "wordpress:6.4", "state": "running", "status": "Up 2 days", "ports": []gin.H{{"ip": "0.0.0.0", "private_port": 80, "public_port": 8080, "type": "tcp"}}, "created_at": "2024-07-22T10:30:00Z"},
				{"id": "q7r8s9t0", "name": "halo", "image": "halohub/halo:2.15", "state": "running", "status": "Up 1 day", "ports": []gin.H{{"ip": "0.0.0.0", "private_port": 8090, "public_port": 8090, "type": "tcp"}}, "created_at": "2024-07-23T10:30:00Z"},
				{"id": "u1v2w3x4", "name": "minio", "image": "minio/minio:latest", "state": "running", "status": "Up 5 days", "ports": []gin.H{{"ip": "0.0.0.0", "private_port": 9000, "public_port": 9000, "type": "tcp"}, {"ip": "0.0.0.0", "private_port": 9001, "public_port": 9001, "type": "tcp"}}, "created_at": "2024-07-19T10:30:00Z"},
				{"id": "y5z6a7b8", "name": "gitea", "image": "gitea/gitea:1.21", "state": "running", "status": "Up 4 days", "ports": []gin.H{{"ip": "0.0.0.0", "private_port": 3000, "public_port": 3000, "type": "tcp"}, {"ip": "0.0.0.0", "private_port": 22, "public_port": 2222, "type": "tcp"}}, "created_at": "2024-07-20T10:30:00Z"},
				{"id": "c9d0e1f2", "name": "n8n", "image": "n8nio/n8n:latest", "state": "stopped", "status": "Exited (0) 2 hours ago", "ports": []gin.H{{"ip": "0.0.0.0", "private_port": 5678, "public_port": 5678, "type": "tcp"}}, "created_at": "2024-07-24T10:30:00Z"},
			}
			successPage(c, containers, 8, 1, 20)
		})

		agent.GET("/images", func(c *gin.Context) {
			images := []gin.H{
				{"id": "sha256:abc123", "repo_tags": []string{"nginx:1.25"}, "size": 187, "created": "2024-07-01"},
				{"id": "sha256:def456", "repo_tags": []string{"mysql:8.0"}, "size": 578, "created": "2024-07-01"},
				{"id": "sha256:ghi789", "repo_tags": []string{"redis:7-alpine"}, "size": 32, "created": "2024-07-01"},
				{"id": "sha256:jkl012", "repo_tags": []string{"wordpress:6.4"}, "size": 682, "created": "2024-07-02"},
				{"id": "sha256:mno345", "repo_tags": []string{"halohub/halo:2.15"}, "size": 412, "created": "2024-07-03"},
			}
			successPage(c, images, 5, 1, 20)
		})

		// Websites
		agent.GET("/websites", func(c *gin.Context) {
			sites := []gin.H{
				{"id": 1, "primary_domain": "example.com", "type": "static", "root_dir": "/var/www/example", "ssl_status": "enabled", "created_at": "2024-07-01"},
				{"id": 2, "primary_domain": "api.example.com", "type": "proxy", "proxy_target": "http://localhost:3000", "ssl_status": "enabled", "created_at": "2024-07-02"},
				{"id": 3, "primary_domain": "blog.example.com", "type": "php", "root_dir": "/var/www/blog", "ssl_status": "disabled", "created_at": "2024-07-05"},
			}
			successPage(c, sites, 3, 1, 20)
		})

		// Files
		agent.GET("/files/list", func(c *gin.Context) {
			dir := c.DefaultQuery("path", "/")
			files := []gin.H{
				{"name": "nginx", "path": dir + "/nginx", "is_dir": true, "size": 4096, "mod_time": "2024-07-20 10:30", "mode": "drwxr-xr-x", "owner": "root"},
				{"name": "docker", "path": dir + "/docker", "is_dir": true, "size": 4096, "mod_time": "2024-07-19 15:22", "mode": "drwxr-xr-x", "owner": "root"},
				{"name": "www", "path": dir + "/www", "is_dir": true, "size": 4096, "mod_time": "2024-07-18 08:15", "mode": "drwxr-xr-x", "owner": "www-data"},
				{"name": "backup", "path": dir + "/backup", "is_dir": true, "size": 4096, "mod_time": "2024-07-22 02:00", "mode": "drwx------", "owner": "root"},
				{"name": "config.yaml", "path": dir + "/config.yaml", "is_dir": false, "size": 2048, "mod_time": "2024-07-21 14:45", "mode": "-rw-r--r--", "owner": "root"},
				{"name": "access.log", "path": dir + "/access.log", "is_dir": false, "size": 524288, "mod_time": "2024-07-24 08:00", "mode": "-rw-r--r--", "owner": "www-data"},
			}
			success(c, files)
		})

		// Cron Jobs
		agent.GET("/cronjobs", func(c *gin.Context) {
			jobs := []gin.H{
				{"id": 1, "name": "每日备份", "spec": "0 2 * * *", "command": "/opt/backup.sh", "type": "script", "status": "enabled", "last_run": "2024-07-24 02:00:01", "next_run": "2024-07-25 02:00:00"},
				{"id": 2, "name": "SSL 证书续期", "spec": "0 3 * * 1", "command": "certbot renew", "type": "shell", "status": "enabled", "last_run": "2024-07-22 03:00:02", "next_run": "2024-07-29 03:00:00"},
				{"id": 3, "name": "日志清理", "spec": "0 4 * * 0", "command": "/opt/clean-logs.sh", "type": "script", "status": "disabled", "last_run": nil, "next_run": nil},
			}
			successPage(c, jobs, 3, 1, 20)
		})

		// App Store
		agent.GET("/appstore/apps", func(c *gin.Context) {
			apps := []gin.H{
				{"key": "wordpress", "name": "WordPress", "category": "CMS", "version": "6.4.2", "description": "Most popular CMS platform", "icon": "wordpress"},
				{"key": "halo", "name": "Halo", "category": "CMS", "version": "2.15", "description": "Modern blogging platform", "icon": "halo"},
				{"key": "mysql", "name": "MySQL", "category": "Database", "version": "8.0", "description": "Relational database", "icon": "mysql"},
				{"key": "redis", "name": "Redis", "category": "Database", "version": "7.2", "description": "In-memory data store", "icon": "redis"},
				{"key": "postgresql", "name": "PostgreSQL", "category": "Database", "version": "16", "description": "Advanced relational database", "icon": "postgresql"},
				{"key": "nginx", "name": "Nginx", "category": "Web Server", "version": "1.25", "description": "High-performance web server", "icon": "nginx"},
				{"key": "n8n", "name": "n8n", "category": "Automation", "version": "latest", "description": "Workflow automation", "icon": "n8n"},
				{"key": "minio", "name": "MinIO", "category": "Storage", "version": "latest", "description": "S3-compatible object storage", "icon": "minio"},
				{"key": "gitea", "name": "Gitea", "category": "DevOps", "version": "1.21", "description": "Self-hosted Git service", "icon": "gitea"},
				{"key": "nextcloud", "name": "Nextcloud", "category": "Storage", "version": "28", "description": "Private cloud platform", "icon": "nextcloud"},
				{"key": "ollama", "name": "Ollama", "category": "AI", "version": "latest", "description": "Local LLM runtime", "icon": "ollama"},
				{"key": "portainer", "name": "Portainer", "category": "DevOps", "version": "2.19", "description": "Container management UI", "icon": "portainer"},
			}
			success(c, apps)
		})

		// SSL Certs
		agent.GET("/ssl/certs", func(c *gin.Context) {
			certs := []gin.H{
				{"id": 1, "domain": "example.com", "provider": "letsencrypt", "ca_type": "RSA", "auto_renew": true, "expired_at": "2024-10-15"},
				{"id": 2, "domain": "api.example.com", "provider": "zerossl", "ca_type": "ECDSA", "auto_renew": true, "expired_at": "2024-11-20"},
			}
			successPage(c, certs, 2, 1, 20)
		})

		// Firewall Rules
		agent.GET("/firewall/rules", func(c *gin.Context) {
			rules := []gin.H{
				{"id": 1, "protocol": "tcp", "port": "80", "source": "0.0.0.0/0", "action": "allow", "comment": "HTTP"},
				{"id": 2, "protocol": "tcp", "port": "443", "source": "0.0.0.0/0", "action": "allow", "comment": "HTTPS"},
				{"id": 3, "protocol": "tcp", "port": "22", "source": "10.0.0.0/8", "action": "allow", "comment": "SSH - Internal"},
				{"id": 4, "protocol": "tcp", "port": "3306", "source": "127.0.0.1", "action": "allow", "comment": "MySQL - Local"},
				{"id": 5, "protocol": "tcp", "port": "9999", "source": "0.0.0.0/0", "action": "allow", "comment": "OpenForge Core"},
			}
			successPage(c, rules, 5, 1, 20)
		})

		// Backup Tasks
		agent.GET("/backup/tasks", func(c *gin.Context) {
			tasks := []gin.H{
				{"id": 1, "name": "全站备份", "target_type": "website", "dest_type": "local", "schedule": "0 2 * * *", "retention": 7, "status": "enabled"},
				{"id": 2, "name": "数据库备份", "target_type": "db", "dest_type": "s3", "schedule": "0 4 * * *", "retention": 30, "status": "enabled"},
			}
			successPage(c, tasks, 2, 1, 20)
		})

		// Database
		agent.GET("/databases/mysql", func(c *gin.Context) {
			dbs := []gin.H{
				{"name": "wordpress", "charset": "utf8mb4", "size": "256MB", "users": 2},
				{"name": "halo", "charset": "utf8mb4", "size": "128MB", "users": 1},
				{"name": "n8n", "charset": "utf8mb4", "size": "64MB", "users": 1},
			}
			success(c, dbs)
		})
	}

	// ==================== Dashboard API Routes ====================
	dashboard := r.Group("/api/v2/dashboard")
	dashboard.Use(authMiddleware())
	{
		dashboard.GET("/overview", func(c *gin.Context) {
			success(c, gin.H{
				"cpu_usage": 23.5, "memory_usage": 61.2, "memory_total": 15.6, "memory_used": 9.55,
				"disk_usage": 45.8, "disk_total": 256.0, "disk_used": 117.2,
				"container_count": 24, "container_running": 12, "container_stopped": 12,
				"network_in": 128.5, "network_out": 42.3,
				"uptime": 259200, "hostname": "openforge-server", "os": "Ubuntu 22.04 LTS",
			})
		})
		dashboard.GET("/cpu", func(c *gin.Context) {
			now := time.Now().Unix()
			history := make([]gin.H, 0)
			for i := 0; i < 60; i++ {
				history = append(history, gin.H{"time": fmt.Sprintf("%d", now-int64(60-i)), "value": 18 + float64(i%20)*0.5})
			}
			success(c, gin.H{"usage": 23.5, "cores": 4, "model": "Intel Xeon E5-2680 v4", "history": history})
		})
		dashboard.GET("/memory", func(c *gin.Context) {
			now := time.Now().Unix()
			history := make([]gin.H, 0)
			for i := 0; i < 60; i++ {
				history = append(history, gin.H{"time": fmt.Sprintf("%d", now-int64(60-i)), "value": 58 + float64(i%15)*0.4})
			}
			success(c, gin.H{"total": 15.6, "used": 9.55, "free": 2.85, "cached": 3.2, "buffers": 0.8, "history": history})
		})
		dashboard.GET("/disk", func(c *gin.Context) {
			success(c, gin.H{
				"total": 256.0, "used": 117.2, "free": 138.8,
				"partitions": []gin.H{
					{"device": "/dev/sda1", "mount": "/", "total": 256.0, "used": 117.2, "free": 138.8, "usage": 45.8},
					{"device": "/dev/sda2", "mount": "/home", "total": 512.0, "used": 89.4, "free": 422.6, "usage": 17.5},
					{"device": "/dev/sdb1", "mount": "/data", "total": 1024.0, "used": 623.8, "free": 400.2, "usage": 60.9},
				},
			})
		})
		dashboard.GET("/network", func(c *gin.Context) {
			now := time.Now().Unix()
			history := make([]gin.H, 0)
			for i := 0; i < 60; i++ {
				history = append(history, gin.H{
					"time": fmt.Sprintf("%d", now-int64(60-i)),
					"rx": float64(80+i%40) * 1.5, "tx": float64(30+i%20) * 1.2,
				})
			}
			success(c, gin.H{
				"interfaces": []gin.H{
					{"name": "eth0", "rx_bytes": 128.5 * 1024 * 1024, "tx_bytes": 42.3 * 1024 * 1024, "rx_speed": 15.2, "tx_speed": 5.8},
				},
				"history": history,
			})
		})
		dashboard.GET("/processes", func(c *gin.Context) {
			success(c, []gin.H{
				{"pid": 1, "name": "systemd", "user": "root", "cpu": 0.0, "memory": 1.2, "status": "running"},
				{"pid": 1023, "name": "nginx", "user": "www-data", "cpu": 0.5, "memory": 2.8, "status": "running"},
				{"pid": 2045, "name": "docker", "user": "root", "cpu": 1.2, "memory": 15.3, "status": "running"},
				{"pid": 3098, "name": "mysql", "user": "mysql", "cpu": 2.1, "memory": 28.7, "status": "running"},
				{"pid": 4521, "name": "redis-server", "user": "redis", "cpu": 0.3, "memory": 4.5, "status": "running"},
				{"pid": 5678, "name": "node", "user": "deploy", "cpu": 3.8, "memory": 12.1, "status": "running"},
				{"pid": 7201, "name": "sshd", "user": "root", "cpu": 0.1, "memory": 0.8, "status": "running"},
				{"pid": 8934, "name": "openforge-core", "user": "root", "cpu": 0.8, "memory": 5.2, "status": "running"},
			})
		})
	}

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		success(c, gin.H{"status": "healthy", "version": "1.0.0", "uptime": time.Now().Unix()})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	log.Printf("openforge-maintain server starting on :%s", port)
	log.Printf("API: http://localhost:%s/api/v2/core/auth/login", port)
	log.Printf("Health: http://localhost:%s/api/health", port)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}