package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"github.com/openforge-maintain/openforge-maintain/agent/service/cron"
	"github.com/openforge-maintain/openforge-maintain/pkg/response"
)

// CronHandler handles cron job API requests.
type CronHandler struct {
	cronManager *cron.CronManager
}

// NewCronHandler creates a new CronHandler.
func NewCronHandler(cronManager *cron.CronManager) *CronHandler {
	return &CronHandler{cronManager: cronManager}
}

// CreateCronJob creates a new cron job.
func (h *CronHandler) CreateCronJob(c *gin.Context) {
	var job model.CronJob
	if err := c.ShouldBindJSON(&job); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.cronManager.Create(&job); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, job)
}

// UpdateCronJob updates a cron job.
func (h *CronHandler) UpdateCronJob(c *gin.Context) {
	var job model.CronJob
	if err := c.ShouldBindJSON(&job); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.cronManager.Update(&job); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, job)
}

// DeleteCronJob deletes a cron job.
func (h *CronHandler) DeleteCronJob(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.cronManager.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListCronJobs lists all cron jobs.
func (h *CronHandler) ListCronJobs(c *gin.Context) {
	jobs, err := h.cronManager.List()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, jobs)
}

// StartCronJob starts a cron job.
func (h *CronHandler) StartCronJob(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.cronManager.Start(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// StopCronJob stops a cron job.
func (h *CronHandler) StopCronJob(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.cronManager.Stop(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// RunCronJobNow immediately executes a cron job.
func (h *CronHandler) RunCronJobNow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.cronManager.RunNow(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}
