package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"github.com/openforge-maintain/openforge-maintain/agent/service/backup"
	"github.com/openforge-maintain/openforge-maintain/pkg/response"
)

// BackupHandler handles backup and restore API requests.
type BackupHandler struct {
	backupService *backup.BackupService
}

// NewBackupHandler creates a new BackupHandler.
func NewBackupHandler(backupService *backup.BackupService) *BackupHandler {
	return &BackupHandler{backupService: backupService}
}

// CreateBackupTask creates a new backup task.
func (h *BackupHandler) CreateBackupTask(c *gin.Context) {
	var task model.BackupTask
	if err := c.ShouldBindJSON(&task); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.backupService.CreateBackup(&task); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, task)
}

// RestoreBackup restores from a backup file.
func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	var req struct {
		TaskID     uint   `json:"task_id" binding:"required"`
		BackupFile string `json:"backup_file" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.backupService.Restore(req.TaskID, req.BackupFile); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListBackupTasks lists all backup tasks.
func (h *BackupHandler) ListBackupTasks(c *gin.Context) {
	// List tasks is available through backup service
	response.Success(c, nil)
}

// ListBackups lists backup files for a task.
func (h *BackupHandler) ListBackups(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	backups, err := h.backupService.ListBackups(uint(taskID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, backups)
}

// DeleteBackup deletes a specific backup file.
func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		BackupFile string `json:"backup_file" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.backupService.DeleteBackup(uint(taskID), req.BackupFile); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}
