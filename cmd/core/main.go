package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/openforge-maintain/openforge-maintain/core/model"
	"github.com/openforge-maintain/openforge-maintain/core/repository"
	"github.com/openforge-maintain/openforge-maintain/core/router"
	"github.com/openforge-maintain/openforge-maintain/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// ========== 读取配置 ==========

	// 服务端口
	port := os.Getenv("MAINTAIN_PORT")
	if port == "" {
		port = "9999"
	}

	// 数据库路径
	dbPath := os.Getenv("MAINTAIN_DB_PATH")
	if dbPath == "" {
		dbPath = "/opt/openforge-maintain/data/core.db"
	}

	// JWT 密钥（写入环境变量供 pkg/utils 使用）
	secret := os.Getenv("MAINTAIN_SECRET")
	if secret == "" {
		secret = "openforge-maintain-secret-key"
		os.Setenv("MAINTAIN_SECRET", secret)
	} else {
		os.Setenv("MAINTAIN_SECRET", secret)
	}

	// 安全入口路径
	securityEntry := os.Getenv("MAINTAIN_SECURITY_ENTRY")
	if securityEntry != "" {
		os.Setenv("MAINTAIN_SECURITY_ENTRY", securityEntry)
	}

	// ========== 初始化数据库 ==========

	// 确保数据库目录存在
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory %s: %v", dbDir, err)
	}

	// 连接 SQLite 数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Printf("Database connected: %s", dbPath)

	// ========== 自动迁移 ==========

	err = db.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.AuditLog{},
		&model.Node{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	log.Println("Database migration completed")

	// ========== 创建默认角色和用户 ==========

	if err := initDefaultData(db); err != nil {
		log.Fatalf("Failed to initialize default data: %v", err)
	}

	// ========== 初始化 Repository ==========

	userRepo := repository.NewUserRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	nodeRepo := repository.NewNodeRepository(db)

	// ========== 初始化路由 ==========

	coreRouter := router.NewCoreRouter(userRepo, auditRepo, nodeRepo)
	engine := coreRouter.Setup()

	// ========== 启动服务 ==========

	gin.SetMode(gin.ReleaseMode)

	addr := ":" + port
	log.Printf("openforge-maintain Core service starting on %s", addr)
	log.Printf("API base path: /api/v2/core")
	log.Printf("Security entry: %s", func() string {
		if securityEntry == "" {
			return "disabled"
		}
		return securityEntry
	}())

	if err := engine.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initDefaultData 初始化默认数据（角色和管理员账户）
func initDefaultData(db *gorm.DB) error {
	// 创建默认角色
	var adminRole model.Role
	result := db.Where("name = ?", "admin").First(&adminRole)
	if result.Error != nil {
		// 管理员角色不存在，创建
		adminRole = model.Role{
			Name:        "admin",
			Description: "系统管理员，拥有所有权限",
		}
		if err := db.Create(&adminRole).Error; err != nil {
			return fmt.Errorf("failed to create admin role: %w", err)
		}
		log.Println("Default admin role created")
	}

	var userRole model.Role
	result = db.Where("name = ?", "user").First(&userRole)
	if result.Error != nil {
		// 普通用户角色不存在，创建
		userRole = model.Role{
			Name:        "user",
			Description: "普通用户，拥有基本权限",
		}
		if err := db.Create(&userRole).Error; err != nil {
			return fmt.Errorf("failed to create user role: %w", err)
		}
		log.Println("Default user role created")
	}

	// 创建默认管理员账户
	var adminUser model.User
	result = db.Where("username = ?", "admin").First(&adminUser)
	if result.Error != nil {
		// 对默认密码进行哈希
		hashedPassword, err := utils.HashPassword("maintain@2024")
		if err != nil {
			return fmt.Errorf("failed to hash admin password: %w", err)
		}

		adminUser = model.User{
			Username: "admin",
			Password: hashedPassword,
			Email:    "admin@openforge-maintain.local",
			RoleID:   adminRole.ID,
		}

		if err := db.Create(&adminUser).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		log.Println("Default admin user created (username: admin, password: maintain@2024)")
		log.Println("WARNING: Please change the default password after first login!")
	}

	return nil
}
