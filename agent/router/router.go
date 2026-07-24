package router

import (
	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/agent/handler"
)

// AgentRouter sets up all routes for the Agent API.
type AgentRouter struct {
	dockerHandler    *handler.DockerHandler
	fileHandler      *handler.FileHandler
	websiteHandler   *handler.WebsiteHandler
	databaseHandler  *handler.DatabaseHandler
	cronHandler      *handler.CronHandler
	appStoreHandler  *handler.AppStoreHandler
	sslHandler       *handler.SSLHandler
	backupHandler    *handler.BackupHandler
	firewallHandler  *handler.FirewallHandler
	metricsHandler   *handler.MetricsHandler
	terminalHandler  *handler.TerminalHandler
	metricsWSHandler *handler.MetricsWSHandler
}

// NewAgentRouter creates a new AgentRouter with all dependencies.
func NewAgentRouter(
	dockerHandler *handler.DockerHandler,
	fileHandler *handler.FileHandler,
	websiteHandler *handler.WebsiteHandler,
	databaseHandler *handler.DatabaseHandler,
	cronHandler *handler.CronHandler,
	appStoreHandler *handler.AppStoreHandler,
	sslHandler *handler.SSLHandler,
	backupHandler *handler.BackupHandler,
	firewallHandler *handler.FirewallHandler,
	metricsHandler *handler.MetricsHandler,
	terminalHandler *handler.TerminalHandler,
	metricsWSHandler *handler.MetricsWSHandler,
) *AgentRouter {
	return &AgentRouter{
		dockerHandler:    dockerHandler,
		fileHandler:      fileHandler,
		websiteHandler:   websiteHandler,
		databaseHandler:  databaseHandler,
		cronHandler:      cronHandler,
		appStoreHandler:  appStoreHandler,
		sslHandler:       sslHandler,
		backupHandler:    backupHandler,
		firewallHandler:  firewallHandler,
		metricsHandler:   metricsHandler,
		terminalHandler:  terminalHandler,
		metricsWSHandler: metricsWSHandler,
	}
}

// Setup configures all routes and returns the Gin engine.
func (r *AgentRouter) Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	apiV2 := engine.Group("/api/v2")
	{
		// Container routes
		containers := apiV2.Group("/containers")
		{
			containers.GET("", r.dockerHandler.ListContainers)
			containers.POST("", r.dockerHandler.CreateContainer)
			containers.GET("/:id", r.dockerHandler.GetContainer)
			containers.POST("/:id/start", r.dockerHandler.StartContainer)
			containers.POST("/:id/stop", r.dockerHandler.StopContainer)
			containers.POST("/:id/restart", r.dockerHandler.RestartContainer)
			containers.DELETE("/:id", r.dockerHandler.RemoveContainer)
			containers.GET("/:id/logs", r.dockerHandler.GetContainerLogs)
			containers.GET("/:id/stats", r.dockerHandler.GetContainerStats)
			containers.POST("/:id/exec", r.dockerHandler.ExecContainer)
		}

		// Image routes
		images := apiV2.Group("/images")
		{
			images.GET("", r.dockerHandler.ListImages)
			images.POST("/pull", r.dockerHandler.PullImage)
			images.DELETE("", r.dockerHandler.RemoveImage)
		}

		// Volume routes
		volumes := apiV2.Group("/volumes")
		{
			volumes.GET("", r.dockerHandler.ListVolumes)
			volumes.POST("", r.dockerHandler.CreateVolume)
			volumes.DELETE("/:name", r.dockerHandler.RemoveVolume)
		}

		// Network routes
		networks := apiV2.Group("/networks")
		{
			networks.GET("", r.dockerHandler.ListNetworks)
			networks.POST("", r.dockerHandler.CreateNetwork)
			networks.DELETE("/:name", r.dockerHandler.RemoveNetwork)
		}

		// Compose routes
		compose := apiV2.Group("/compose")
		{
			compose.POST("/up", r.dockerHandler.ComposeUp)
			compose.POST("/down", r.dockerHandler.ComposeDown)
			compose.GET("/list", r.dockerHandler.ComposeList)
			compose.POST("/pull", r.dockerHandler.ComposePullImages)
		}

		// File management routes
		files := apiV2.Group("/files")
		{
			files.GET("/list", r.fileHandler.ListDir)
			files.GET("/content", r.fileHandler.GetFileContent)
			files.POST("/content", r.fileHandler.SaveFileContent)
			files.POST("/upload", r.fileHandler.Upload)
			files.GET("/download", r.fileHandler.Download)
			files.DELETE("", r.fileHandler.DeleteFile)
			files.POST("/rename", r.fileHandler.RenameFile)
			files.POST("/chmod", r.fileHandler.ChmodFile)
			files.POST("/chown", r.fileHandler.ChownFile)
			files.POST("/compress", r.fileHandler.CompressFiles)
			files.POST("/extract", r.fileHandler.ExtractFiles)
			files.POST("/mkdir", r.fileHandler.Mkdir)
			files.GET("/info", r.fileHandler.GetFileInfo)
		}

		// Website routes
		websites := apiV2.Group("/websites")
		{
			websites.GET("", r.websiteHandler.ListWebsites)
			websites.POST("", r.websiteHandler.CreateWebsite)
			websites.GET("/:id", r.websiteHandler.GetWebsite)
			websites.PUT("/:id", r.websiteHandler.UpdateWebsite)
			websites.DELETE("/:id", r.websiteHandler.DeleteWebsite)
			websites.GET("/:id/nginx", r.websiteHandler.GetNginxConfig)
			websites.POST("/ssl/enable", r.websiteHandler.EnableSSL)
			websites.POST("/ssl/disable", r.websiteHandler.DisableSSL)
		}

		// Database routes
		databases := apiV2.Group("/databases")
		{
			// MySQL
			mysql := databases.Group("/mysql")
			{
				mysql.GET("/databases", r.databaseHandler.ListMySQLDatabases)
				mysql.POST("/databases", r.databaseHandler.CreateMySQLDatabase)
				mysql.DELETE("/databases", r.databaseHandler.DeleteMySQLDatabase)
				mysql.POST("/users", r.databaseHandler.CreateMySQLUser)
				mysql.POST("/backup", r.databaseHandler.BackupMySQL)
				mysql.POST("/restore", r.databaseHandler.RestoreMySQL)
			}

			// PostgreSQL
			postgres := databases.Group("/postgres")
			{
				postgres.GET("/databases", r.databaseHandler.ListPostgresDatabases)
				postgres.POST("/databases", r.databaseHandler.CreatePostgresDatabase)
				postgres.DELETE("/databases", r.databaseHandler.DeletePostgresDatabase)
				postgres.POST("/users", r.databaseHandler.CreatePostgresUser)
				postgres.POST("/backup", r.databaseHandler.BackupPostgres)
				postgres.POST("/restore", r.databaseHandler.RestorePostgres)
			}

			// Redis
			redis := databases.Group("/redis")
			{
				redis.GET("/info", r.databaseHandler.GetRedisInfo)
				redis.POST("/config", r.databaseHandler.SetRedisConfig)
				redis.GET("/databases", r.databaseHandler.GetRedisDatabases)
			}
		}

		// Cron job routes
		cronjobs := apiV2.Group("/cronjobs")
		{
			cronjobs.GET("", r.cronHandler.ListCronJobs)
			cronjobs.POST("", r.cronHandler.CreateCronJob)
			cronjobs.PUT("/:id", r.cronHandler.UpdateCronJob)
			cronjobs.DELETE("/:id", r.cronHandler.DeleteCronJob)
			cronjobs.POST("/:id/start", r.cronHandler.StartCronJob)
			cronjobs.POST("/:id/stop", r.cronHandler.StopCronJob)
			cronjobs.POST("/:id/run", r.cronHandler.RunCronJobNow)
		}

		// App store routes
		appstore := apiV2.Group("/appstore")
		{
			apps := appstore.Group("/apps")
			{
				apps.GET("", r.appStoreHandler.GetAppList)
				apps.GET("/:appKey", r.appStoreHandler.GetAppDetail)
			}
			installs := appstore.Group("/installs")
			{
				installs.GET("", r.appStoreHandler.GetInstalledApps)
				installs.POST("/install", r.appStoreHandler.InstallApp)
				installs.DELETE("/:id", r.appStoreHandler.UninstallApp)
				installs.POST("/:id/upgrade", r.appStoreHandler.UpgradeApp)
			}
		}

		// SSL routes
		ssl := apiV2.Group("/ssl/certs")
		{
			ssl.GET("", r.sslHandler.ListCerts)
			ssl.POST("/request", r.sslHandler.RequestCert)
			ssl.GET("/:id", r.sslHandler.GetCertDetail)
			ssl.POST("/:id/renew", r.sslHandler.RenewCert)
			ssl.DELETE("/:id", r.sslHandler.DeleteCert)
		}

		// Backup routes
		backup := apiV2.Group("/backup")
		{
			backup.GET("/tasks", r.backupHandler.ListBackupTasks)
			backup.POST("/tasks", r.backupHandler.CreateBackupTask)
			backup.GET("/tasks/:id/backups", r.backupHandler.ListBackups)
			backup.DELETE("/tasks/:id/backups", r.backupHandler.DeleteBackup)
			backup.POST("/restore", r.backupHandler.RestoreBackup)
		}

		// Firewall routes
		firewall := apiV2.Group("/firewall")
		{
			firewall.GET("/rules", r.firewallHandler.ListRules)
			firewall.POST("/rules", r.firewallHandler.AddRule)
			firewall.DELETE("/rules/:id", r.firewallHandler.DeleteRule)
			firewall.POST("/rules/:id/enable", r.firewallHandler.EnableRule)
			firewall.POST("/rules/:id/disable", r.firewallHandler.DisableRule)
			firewall.GET("/status", r.firewallHandler.GetFirewallStatus)
			firewall.POST("/enable", r.firewallHandler.EnableFirewall)
			firewall.POST("/disable", r.firewallHandler.DisableFirewall)
		}

		// Metrics routes
		metrics := apiV2.Group("/metrics")
		{
			metrics.GET("/cpu", r.metricsHandler.GetCPU)
			metrics.GET("/memory", r.metricsHandler.GetMemory)
			metrics.GET("/disk", r.metricsHandler.GetDisk)
			metrics.GET("/network", r.metricsHandler.GetNetwork)
			metrics.GET("/overview", r.metricsHandler.GetOverview)
			metrics.GET("/processes", r.metricsHandler.GetProcesses)
		}
	}

	// WebSocket routes (not under /api/v2)
	ws := engine.Group("/ws")
	{
		ws.GET("/terminal", r.terminalHandler.HandleTerminal)
		ws.GET("/terminal/sessions", r.terminalHandler.GetSessions)
		ws.GET("/metrics", r.metricsWSHandler.HandleMetricsWS)
	}

	return engine
}
