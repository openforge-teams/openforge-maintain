package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/agent/service/system"
	"github.com/openforge-maintain/openforge-maintain/pkg/response"
)

// MetricsHandler handles system monitoring API requests.
type MetricsHandler struct {
	metricsService *system.MetricsService
}

// NewMetricsHandler creates a new MetricsHandler.
func NewMetricsHandler(metricsService *system.MetricsService) *MetricsHandler {
	return &MetricsHandler{metricsService: metricsService}
}

// GetCPU returns CPU usage information.
func (h *MetricsHandler) GetCPU(c *gin.Context) {
	cpu, err := h.metricsService.GetCPU()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cpu)
}

// GetMemory returns memory usage information.
func (h *MetricsHandler) GetMemory(c *gin.Context) {
	mem, err := h.metricsService.GetMemory()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, mem)
}

// GetDisk returns disk usage information.
func (h *MetricsHandler) GetDisk(c *gin.Context) {
	disks, err := h.metricsService.GetDisk()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, disks)
}

// GetNetwork returns network interface statistics.
func (h *MetricsHandler) GetNetwork(c *gin.Context) {
	networks, err := h.metricsService.GetNetwork()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, networks)
}

// GetOverview returns a comprehensive system overview.
func (h *MetricsHandler) GetOverview(c *gin.Context) {
	overview, err := h.metricsService.GetOverview()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, overview)
}

// GetProcesses returns a paginated list of processes.
func (h *MetricsHandler) GetProcesses(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"), 1)
	size := parseInt(c.DefaultQuery("size", "20"), 20)

	processes, total, err := h.metricsService.GetProcesses(page, size)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"list":  processes,
		"total": total,
		"page":  page,
		"size":  size,
	})
}
