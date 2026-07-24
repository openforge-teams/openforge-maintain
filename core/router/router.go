package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/core/handler"
	"github.com/openforge-maintain/openforge-maintain/core/middleware"
	"github.com/openforge-maintain/openforge-maintain/core/repository"
	"github.com/openforge-maintain/openforge-maintain/core/service"
)

// CoreRouter 核心路由器
type CoreRouter struct {
	engine    *gin.Engine
	authRepo  repository.UserRepository
	auditRepo repository.AuditRepository
	nodeRepo  repository.NodeRepository
}

// NewCoreRouter 创建核心路由器
func NewCoreRouter(
	authRepo repository.UserRepository,
	auditRepo repository.AuditRepository,
	nodeRepo repository.NodeRepository,
) *CoreRouter {
	return &CoreRouter{
		engine:    gin.New(),
		authRepo:  authRepo,
		auditRepo: auditRepo,
		nodeRepo:  nodeRepo,
	}
}

// Setup 设置路由
func (r *CoreRouter) Setup() *gin.Engine {
	// 全局中间件
	r.engine.Use(gin.Recovery())

	// 安全入口中间件
	r.engine.Use(middleware.SecurityEntry())

	// CORS 中间件
	r.engine.Use(middleware.CORS())

	// 速率限制（每分钟 100 次请求）
	r.engine.Use(middleware.RateLimit(100, time.Minute))

	// 审计日志中间件
	r.engine.Use(middleware.Audit(r.auditRepo))

	// 初始化 services
	authService := service.NewAuthService(r.authRepo)
	userService := service.NewUserService(r.authRepo)
	nodeService := service.NewNodeService(r.nodeRepo)

	// 初始化 handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, authService)
	nodeHandler := handler.NewNodeHandler(nodeService)

	// API v2 路由组
	v2 := r.engine.Group("/api/v2/core")

	// 公开路由（无需认证）
	public := v2.Group("/auth")
	{
		public.POST("/login", authHandler.Login)
		public.POST("/register", authHandler.Register)
		public.POST("/refresh", authHandler.Refresh)
	}

	// 认证路由（需要 JWT 认证）
	protected := v2.Group("")
	protected.Use(middleware.Auth())
	{
		// 认证相关
		protected.POST("/auth/logout", authHandler.Logout)

		// 用户相关
		users := protected.Group("/users")
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.PUT("/password", userHandler.ChangePassword)
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// 节点相关
		nodes := protected.Group("/nodes")
		{
			nodes.GET("", nodeHandler.ListNodes)
			nodes.POST("", nodeHandler.CreateNode)
			nodes.GET("/:id", nodeHandler.GetNode)
			nodes.PUT("/:id", nodeHandler.UpdateNode)
			nodes.DELETE("/:id", nodeHandler.DeleteNode)
			nodes.GET("/:id/status", nodeHandler.CheckNodeStatus)
		}
	}

	// 健康检查
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    0,
			"message": "healthy",
			"data":    nil,
		})
	})

	return r.engine
}

// GetEngine 获取 Gin 引擎
func (r *CoreRouter) GetEngine() *gin.Engine {
	return r.engine
}
