package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/glebarez/sqlite"
	"github.com/openforge-maintain/openforge-maintain/agent/handler"
	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"github.com/openforge-maintain/openforge-maintain/agent/repository"
	"github.com/openforge-maintain/openforge-maintain/agent/router"
	"github.com/openforge-maintain/openforge-maintain/agent/service/appstore"
	"github.com/openforge-maintain/openforge-maintain/agent/service/backup"
	"github.com/openforge-maintain/openforge-maintain/agent/service/cron"
	"github.com/openforge-maintain/openforge-maintain/agent/service/database"
	"github.com/openforge-maintain/openforge-maintain/agent/service/docker"
	"github.com/openforge-maintain/openforge-maintain/agent/service/file"
	"github.com/openforge-maintain/openforge-maintain/agent/service/firewall"
	"github.com/openforge-maintain/openforge-maintain/agent/service/ssl"
	"github.com/openforge-maintain/openforge-maintain/agent/service/system"
	"github.com/openforge-maintain/openforge-maintain/agent/service/website"
	"gorm.io/gorm"
)

func main() {
	// Read configuration from environment variables
	port := os.Getenv("AGENT_PORT")
	if port == "" {
		port = "10000"
	}

	dbPath := os.Getenv("AGENT_DB_PATH")
	if dbPath == "" {
		dbPath = "agent.db"
	}

	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
	}

	// Initialize SQLite database
	db, err := initDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Docker client
	dockerClient, err := client.NewClientWithOpts(client.WithHost(dockerHost), client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Failed to initialize Docker client: %v", err)
	}
	defer dockerClient.Close()

	// Initialize repositories
	appRepo := repository.NewAppRepository(db)
	websiteRepo := repository.NewWebsiteRepository(db)
	sslRepo := repository.NewSSLRepository(db)
	cronRepo := repository.NewCronRepository(db)
	backupRepo := repository.NewBackupRepository(db)
	firewallRepo := repository.NewFirewallRepository(db)

	// Initialize services
	containerService := docker.NewContainerService(dockerClient)
	imageService := docker.NewImageService(dockerClient)
	volumeService := docker.NewVolumeService(dockerClient)
	networkService := docker.NewNetworkService(dockerClient)
	composeService := docker.NewComposeService(dockerHost)

	fileService := file.NewFileManagerService("/")

	websiteService := website.NewWebsiteService(websiteRepo, sslRepo, db)

	mysqlService := database.NewMySQLService("mysql", "")
	postgresService := database.NewPostgresService("postgresql", "postgres", "")
	redisService := database.NewRedisService("redis")

	cronManager := cron.NewCronManager(cronRepo)

	appRegistry := appstore.NewAppRegistry()
	appStoreService := appstore.NewAppStoreService(appRepo, appRegistry)

	sslService := ssl.NewSSLService(sslRepo)

	backupService := backup.NewBackupService(backupRepo)

	firewallService := firewall.NewFirewallService(firewallRepo)

	metricsService := system.NewMetricsService()

	// Initialize handlers
	dockerHandler := handler.NewDockerHandler(containerService, imageService, volumeService, networkService, composeService)
	fileHandler := handler.NewFileHandler(fileService)
	websiteHandler := handler.NewWebsiteHandler(websiteService)
	databaseHandler := handler.NewDatabaseHandler(mysqlService, postgresService, redisService)
	cronHandler := handler.NewCronHandler(cronManager)
	appStoreHandler := handler.NewAppStoreHandler(appStoreService)
	sslHandler := handler.NewSSLHandler(sslService)
	backupHandler := handler.NewBackupHandler(backupService)
	firewallHandler := handler.NewFirewallHandler(firewallService)
	metricsHandler := handler.NewMetricsHandler(metricsService)
	terminalHandler := handler.NewTerminalHandler()
	metricsWSHandler := handler.NewMetricsWSHandler(nil)

	// Setup router
	agentRouter := router.NewAgentRouter(
		dockerHandler,
		fileHandler,
		websiteHandler,
		databaseHandler,
		cronHandler,
		appStoreHandler,
		sslHandler,
		backupHandler,
		firewallHandler,
		metricsHandler,
		terminalHandler,
		metricsWSHandler,
	)

	engine := agentRouter.Setup()

	// Start cron scheduler
	if err := cronManager.StartScheduler(); err != nil {
		log.Printf("Warning: failed to start cron scheduler: %v", err)
	}
	defer cronManager.StopScheduler()

	// Start system metrics collector (runs in background)
	go func() {
		// Metrics collector would periodically collect and cache metrics
		// In a production system, this would use a ticker
	}()

	// Start HTTP server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Agent server starting on %s", addr)
	log.Printf("Docker host: %s", dockerHost)
	log.Printf("Database path: %s", dbPath)

	if err := engine.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initDB initializes the SQLite database and runs auto-migration.
func initDB(dbPath string) (*gorm.DB, error) {
	// Ensure database directory exists
	dbDir := filepath.Dir(dbPath)
	if dbDir != "." && dbDir != "" {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&model.AppInstall{},
		&model.Website{},
		&model.SSLCert{},
		&model.CronJob{},
		&model.BackupTask{},
		&model.FirewallRule{},
	); err != nil {
		return nil, fmt.Errorf("failed to run auto-migration: %w", err)
	}

	log.Println("Database initialized and migrated successfully")
	return db, nil
}
