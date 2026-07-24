package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/agent/service/docker"
	"github.com/openforge-maintain/openforge-maintain/pkg/response"
)

// DockerHandler handles container-related API requests.
type DockerHandler struct {
	containerService *docker.ContainerService
	imageService     *docker.ImageService
	volumeService    *docker.VolumeService
	networkService   *docker.NetworkService
	composeService   *docker.ComposeService
}

// NewDockerHandler creates a new DockerHandler.
func NewDockerHandler(
	containerService *docker.ContainerService,
	imageService *docker.ImageService,
	volumeService *docker.VolumeService,
	networkService *docker.NetworkService,
	composeService *docker.ComposeService,
) *DockerHandler {
	return &DockerHandler{
		containerService: containerService,
		imageService:     imageService,
		volumeService:    volumeService,
		networkService:   networkService,
		composeService:   composeService,
	}
}

// ListContainers lists all containers.
func (h *DockerHandler) ListContainers(c *gin.Context) {
	all := c.DefaultQuery("all", "false") == "true"
	containers, err := h.containerService.List(c.Request.Context(), all)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, containers)
}

// GetContainer returns details of a specific container.
func (h *DockerHandler) GetContainer(c *gin.Context) {
	id := c.Param("id")
	container, err := h.containerService.Get(c.Request.Context(), id)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, container)
}

// CreateContainer creates a new container.
func (h *DockerHandler) CreateContainer(c *gin.Context) {
	var config docker.ContainerConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.containerService.Create(c.Request.Context(), config); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// StartContainer starts a container.
func (h *DockerHandler) StartContainer(c *gin.Context) {
	id := c.Param("id")
	if err := h.containerService.Start(c.Request.Context(), id); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// StopContainer stops a container.
func (h *DockerHandler) StopContainer(c *gin.Context) {
	id := c.Param("id")
	if err := h.containerService.Stop(c.Request.Context(), id); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// RestartContainer restarts a container.
func (h *DockerHandler) RestartContainer(c *gin.Context) {
	id := c.Param("id")
	if err := h.containerService.Restart(c.Request.Context(), id); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// RemoveContainer removes a container.
func (h *DockerHandler) RemoveContainer(c *gin.Context) {
	id := c.Param("id")
	if err := h.containerService.Remove(c.Request.Context(), id); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetContainerLogs returns container logs.
func (h *DockerHandler) GetContainerLogs(c *gin.Context) {
	id := c.Param("id")
	tail := c.DefaultQuery("tail", "100")

	logs, err := h.containerService.Logs(c.Request.Context(), id, tail)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	c.String(http.StatusOK, logs)
}

// GetContainerStats returns container resource usage.
func (h *DockerHandler) GetContainerStats(c *gin.Context) {
	id := c.Param("id")
	stats, err := h.containerService.Stats(c.Request.Context(), id)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// ExecContainer executes a command inside a container.
func (h *DockerHandler) ExecContainer(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Command []string `json:"command"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	output, err := h.containerService.Exec(c.Request.Context(), id, req.Command)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"output": output})
}

// ListImages lists all Docker images.
func (h *DockerHandler) ListImages(c *gin.Context) {
	images, err := h.imageService.List(c.Request.Context())
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, images)
}

// PullImage pulls a Docker image.
func (h *DockerHandler) PullImage(c *gin.Context) {
	var req struct {
		Image string `json:"image" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.imageService.Pull(c.Request.Context(), req.Image); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// RemoveImage removes a Docker image.
func (h *DockerHandler) RemoveImage(c *gin.Context) {
	var req struct {
		Image string `json:"image" binding:"required"`
		Force bool   `json:"force"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.imageService.Remove(c.Request.Context(), req.Image, req.Force); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListVolumes lists all Docker volumes.
func (h *DockerHandler) ListVolumes(c *gin.Context) {
	volumes, err := h.volumeService.List(c.Request.Context())
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, volumes)
}

// CreateVolume creates a new Docker volume.
func (h *DockerHandler) CreateVolume(c *gin.Context) {
	var req struct {
		Name   string            `json:"name" binding:"required"`
		Driver string            `json:"driver"`
		Labels map[string]string `json:"labels"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.volumeService.Create(c.Request.Context(), req.Name, req.Driver, req.Labels); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// RemoveVolume removes a Docker volume.
func (h *DockerHandler) RemoveVolume(c *gin.Context) {
	name := c.Param("name")
	force := c.DefaultQuery("force", "false") == "true"

	if err := h.volumeService.Remove(c.Request.Context(), name, force); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListNetworks lists all Docker networks.
func (h *DockerHandler) ListNetworks(c *gin.Context) {
	networks, err := h.networkService.List(c.Request.Context())
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, networks)
}

// CreateNetwork creates a new Docker network.
func (h *DockerHandler) CreateNetwork(c *gin.Context) {
	var req struct {
		Name   string            `json:"name" binding:"required"`
		Driver string            `json:"driver"`
		Labels map[string]string `json:"labels"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.networkService.Create(c.Request.Context(), req.Name, req.Driver, req.Labels); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// RemoveNetwork removes a Docker network.
func (h *DockerHandler) RemoveNetwork(c *gin.Context) {
	name := c.Param("name")
	if err := h.networkService.Remove(c.Request.Context(), name); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ComposeUp runs docker compose up.
func (h *DockerHandler) ComposeUp(c *gin.Context) {
	var req struct {
		ComposePath string `json:"compose_path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.composeService.Up(req.ComposePath); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ComposeDown runs docker compose down.
func (h *DockerHandler) ComposeDown(c *gin.Context) {
	var req struct {
		ComposePath string `json:"compose_path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.composeService.Down(req.ComposePath); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ComposeList runs docker compose ps.
func (h *DockerHandler) ComposeList(c *gin.Context) {
	composePath := c.Query("compose_path")
	if composePath == "" {
		response.BadRequest(c, "compose_path is required")
		return
	}

	output, err := h.composeService.List(composePath)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	c.String(http.StatusOK, output)
}

// ComposePullImages runs docker compose pull.
func (h *DockerHandler) ComposePullImages(c *gin.Context) {
	var req struct {
		ComposePath string `json:"compose_path" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.composeService.PullImages(req.ComposePath); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// parseInt parses a string to int with a default value.
func parseInt(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}
