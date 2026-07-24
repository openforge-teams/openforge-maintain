package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/agent/service/database"
	"github.com/openforge-maintain/openforge-maintain/pkg/response"
)

// DatabaseHandler handles database management API requests.
type DatabaseHandler struct {
	mysqlService    *database.MySQLService
	postgresService *database.PostgresService
	redisService    *database.RedisService
}

// NewDatabaseHandler creates a new DatabaseHandler.
func NewDatabaseHandler(
	mysqlService *database.MySQLService,
	postgresService *database.PostgresService,
	redisService *database.RedisService,
) *DatabaseHandler {
	return &DatabaseHandler{
		mysqlService:    mysqlService,
		postgresService: postgresService,
		redisService:    redisService,
	}
}

// ListMySQLDatabases lists all MySQL databases.
func (h *DatabaseHandler) ListMySQLDatabases(c *gin.Context) {
	dbs, err := h.mysqlService.List()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, dbs)
}

// CreateMySQLDatabase creates a new MySQL database.
func (h *DatabaseHandler) CreateMySQLDatabase(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Charset string `json:"charset"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.mysqlService.CreateDB(req.Name, req.Charset); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// DeleteMySQLDatabase drops a MySQL database.
func (h *DatabaseHandler) DeleteMySQLDatabase(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.mysqlService.DeleteDB(req.Name); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// CreateMySQLUser creates a MySQL user.
func (h *DatabaseHandler) CreateMySQLUser(c *gin.Context) {
	var req struct {
		Database string `json:"database" binding:"required"`
		User     string `json:"user" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.mysqlService.CreateUser(req.Database, req.User, req.Password); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// BackupMySQL creates a MySQL database backup.
func (h *DatabaseHandler) BackupMySQL(c *gin.Context) {
	var req struct {
		Database string `json:"database" binding:"required"`
		Dest     string `json:"dest" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.mysqlService.Backup(req.Database, req.Dest); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// RestoreMySQL restores a MySQL database from backup.
func (h *DatabaseHandler) RestoreMySQL(c *gin.Context) {
	var req struct {
		Database string `json:"database" binding:"required"`
		Source   string `json:"source" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.mysqlService.Restore(req.Database, req.Source); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListPostgresDatabases lists all PostgreSQL databases.
func (h *DatabaseHandler) ListPostgresDatabases(c *gin.Context) {
	dbs, err := h.postgresService.List()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, dbs)
}

// CreatePostgresDatabase creates a new PostgreSQL database.
func (h *DatabaseHandler) CreatePostgresDatabase(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Charset string `json:"charset"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.postgresService.CreateDB(req.Name, req.Charset); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// DeletePostgresDatabase drops a PostgreSQL database.
func (h *DatabaseHandler) DeletePostgresDatabase(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.postgresService.DeleteDB(req.Name); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// CreatePostgresUser creates a PostgreSQL user.
func (h *DatabaseHandler) CreatePostgresUser(c *gin.Context) {
	var req struct {
		Database string `json:"database" binding:"required"`
		User     string `json:"user" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.postgresService.CreateUser(req.Database, req.User, req.Password); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// BackupPostgres creates a PostgreSQL database backup.
func (h *DatabaseHandler) BackupPostgres(c *gin.Context) {
	var req struct {
		Database string `json:"database" binding:"required"`
		Dest     string `json:"dest" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.postgresService.Backup(req.Database, req.Dest); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// RestorePostgres restores a PostgreSQL database from backup.
func (h *DatabaseHandler) RestorePostgres(c *gin.Context) {
	var req struct {
		Database string `json:"database" binding:"required"`
		Source   string `json:"source" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.postgresService.Restore(req.Database, req.Source); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetRedisInfo returns Redis server information.
func (h *DatabaseHandler) GetRedisInfo(c *gin.Context) {
	info, err := h.redisService.GetInfo()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, info)
}

// SetRedisConfig sets a Redis configuration parameter.
func (h *DatabaseHandler) SetRedisConfig(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.redisService.SetConfig(req.Key, req.Value); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetRedisDatabases returns Redis database information.
func (h *DatabaseHandler) GetRedisDatabases(c *gin.Context) {
	dbs, err := h.redisService.GetDatabases()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, dbs)
}
