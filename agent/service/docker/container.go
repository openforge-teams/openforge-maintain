package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// ContainerDTO is a simplified container representation for API responses.
type ContainerDTO struct {
	ID      string   `json:"id"`
	Names   []string `json:"names"`
	Image   string   `json:"image"`
	State   string   `json:"state"`
	Status  string   `json:"status"`
	Ports   []string `json:"ports"`
	Created int64    `json:"created"`
}

// ContainerDetail holds detailed container information.
type ContainerDetail struct {
	ContainerDTO
	Image    string                        `json:"image"`
	Command  string                        `json:"command"`
	CreatedAt int64                         `json:"created_at"`
	IP       string                        `json:"ip"`
	Env      []string                      `json:"env"`
	Labels   map[string]string             `json:"labels"`
	Mounts   []types.MountPoint            `json:"mounts"`
	Networks map[string]*types.EndpointSettings `json:"networks"`
}

// ContainerConfig holds the configuration for creating a container.
type ContainerConfig struct {
	Image         string            `json:"image"`
	Name          string            `json:"name"`
	Cmd           []string          `json:"cmd"`
	Env           []string          `json:"env"`
	Ports         map[string]string `json:"ports"`
	Volumes       map[string]string `json:"volumes"`
	RestartPolicy string            `json:"restart_policy"`
	Network       string            `json:"network"`
	Labels        map[string]string `json:"labels"`
}

// ContainerStats holds container resource usage statistics.
type ContainerStats struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	CPUPercent float64 `json:"cpu_percent"`
	Memory     struct {
		Usage   uint64  `json:"usage"`
		Limit   uint64  `json:"limit"`
		Percent float64 `json:"percent"`
	} `json:"memory"`
	Network struct {
		RxBytes uint64 `json:"rx_bytes"`
		TxBytes uint64 `json:"tx_bytes"`
	} `json:"network"`
	BlockIO struct {
		ReadBytes  uint64 `json:"read_bytes"`
		WriteBytes uint64 `json:"write_bytes"`
	} `json:"block_io"`
}

// ContainerService provides container management operations.
type ContainerService struct {
	cli *client.Client
}

// NewContainerService creates a new ContainerService.
func NewContainerService(cli *client.Client) *ContainerService {
	return &ContainerService{cli: cli}
}

// List returns a list of containers.
func (s *ContainerService) List(ctx context.Context, all bool) ([]ContainerDTO, error) {
	containers, err := s.cli.ContainerList(ctx, types.ContainerListOptions{All: all})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	result := make([]ContainerDTO, 0, len(containers))
	for _, c := range containers {
		ports := make([]string, 0, len(c.Ports))
		for _, p := range c.Ports {
			if p.PublicPort != 0 {
				ports = append(ports, fmt.Sprintf("%d->%d/%s", p.PublicPort, p.PrivatePort, p.Type))
			}
		}
		result = append(result, ContainerDTO{
			ID:      c.ID[:12],
			Names:   c.Names,
			Image:   c.Image,
			State:   c.State,
			Status:  c.Status,
			Ports:   ports,
			Created: c.Created,
		})
	}
	return result, nil
}

// Get returns detailed information about a container.
func (s *ContainerService) Get(ctx context.Context, id string) (*ContainerDetail, error) {
	containerJSON, err := s.cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	var ip string
	for _, n := range containerJSON.NetworkSettings.Networks {
		if ip == "" {
			ip = n.IPAddress
		}
	}

	return &ContainerDetail{
		ContainerDTO: ContainerDTO{
			ID:      containerJSON.ID[:12],
			Names:   []string{containerJSON.Name},
			Image:   containerJSON.Config.Image,
			State:   containerJSON.State.Status,
			Status:  containerJSON.State.String(),
			Created: containerJSON.Created,
		},
		Image:     containerJSON.Config.Image,
		Command:   strings.Join(containerJSON.Config.Cmd, " "),
		CreatedAt: containerJSON.Created,
		IP:        ip,
		Env:       containerJSON.Config.Env,
		Labels:    containerJSON.Config.Labels,
		Mounts:    containerJSON.Mounts,
		Networks:  containerJSON.NetworkSettings.Networks,
	}, nil
}

// Create creates and starts a new container.
func (s *ContainerService) Create(ctx context.Context, config ContainerConfig) error {
	// Parse port bindings
	portBindings := make(map[nat.Port][]nat.PortBinding)
	exposedPorts := make(nat.PortSet)
	for hostPort, containerPort := range config.Ports {
		p, err := nat.NewPort("tcp", containerPort)
		if err != nil {
			return fmt.Errorf("invalid port %s: %w", containerPort, err)
		}
		exposedPorts[p] = struct{}{}
		portBindings[p] = []nat.PortBinding{
			{HostIP: "0.0.0.0", HostPort: hostPort},
		}
	}

	// Parse restart policy
	var restartPolicy container.RestartPolicy
	switch config.RestartPolicy {
	case "always":
		restartPolicy = container.RestartPolicy{Name: "always"}
	case "unless-stopped":
		restartPolicy = container.RestartPolicy{Name: "unless-stopped"}
	case "on-failure":
		restartPolicy = container.RestartPolicy{Name: "on-failure"}
	default:
		restartPolicy = container.RestartPolicy{Name: "no"}
	}

	// Parse volume bindings
	var binds []string
	for hostPath, containerPath := range config.Volumes {
		binds = append(binds, fmt.Sprintf("%s:%s", hostPath, containerPath))
	}

	containerConfig := &container.Config{
		Image:        config.Image,
		Cmd:          config.Cmd,
		Env:          config.Env,
		Labels:       config.Labels,
		ExposedPorts: exposedPorts,
	}

	hostConfig := &container.HostConfig{
		PortBindings:  portBindings,
		Binds:         binds,
		RestartPolicy: restartPolicy,
	}

	networkingConfig := &network.NetworkingConfig{}
	if config.Network != "" {
		networkingConfig.EndpointsConfig = map[string]*network.EndpointSettings{
			config.Network: {},
		}
	}

	resp, err := s.cli.ContainerCreate(ctx, containerConfig, hostConfig, networkingConfig, nil, config.Name)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	if err := s.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

// Start starts a stopped container.
func (s *ContainerService) Start(ctx context.Context, id string) error {
	return s.cli.ContainerStart(ctx, id, types.ContainerStartOptions{})
}

// Stop stops a running container.
func (s *ContainerService) Stop(ctx context.Context, id string) error {
	return s.cli.ContainerStop(ctx, id, container.StopOptions{})
}

// Restart restarts a container.
func (s *ContainerService) Restart(ctx context.Context, id string) error {
	timeout := int(10)
	return s.cli.ContainerRestart(ctx, id, container.StopOptions{Timeout: &timeout})
}

// Remove removes a container.
func (s *ContainerService) Remove(ctx context.Context, id string) error {
	return s.cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{Force: true})
}

// Logs returns the logs of a container.
func (s *ContainerService) Logs(ctx context.Context, id string, tail string) (string, error) {
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
	}
	if tail != "" && tail != "all" {
		options.Tail = tail
	}

	reader, err := s.cli.ContainerLogs(ctx, id, options)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read container logs: %w", err)
	}

	return string(data), nil
}

// Stats returns resource usage statistics for a container.
func (s *ContainerService) Stats(ctx context.Context, id string) (*ContainerStats, error) {
	resp, err := s.cli.ContainerStats(ctx, id, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get container stats: %w", err)
	}
	defer resp.Body.Close()

	var stats types.StatsJSON
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode container stats: %w", err)
	}

	result := &ContainerStats{
		ID:   stats.ID[:12],
		Name: strings.TrimPrefix(stats.Name, "/"),
	}

	// CPU calculation
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemCPUUsage - stats.PreCPUStats.SystemCPUUsage)
	if systemDelta > 0 && cpuDelta > 0 {
		result.CPUPercent = (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	// Memory calculation
	result.Memory.Usage = stats.MemoryStats.Usage
	result.Memory.Limit = stats.MemoryStats.Limit
	if result.Memory.Limit > 0 {
		result.Memory.Percent = float64(result.Memory.Usage) / float64(result.Memory.Limit) * 100.0
	}

	// Network calculation
	for _, n := range stats.Networks {
		result.Network.RxBytes += n.RxBytes
		result.Network.TxBytes += n.TxBytes
	}

	// Block IO
	for _, entry := range stats.BlkioStats.IoServiceBytesRecursive {
		switch entry.Op {
		case "Read":
			result.BlockIO.ReadBytes += entry.Value
		case "Write":
			result.BlockIO.WriteBytes += entry.Value
		}
	}

	return result, nil
}

// Exec executes a command inside a running container.
func (s *ContainerService) Exec(ctx context.Context, id string, cmd []string) (string, error) {
	execConfig, err := s.cli.ContainerExecCreate(ctx, id, types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create exec: %w", err)
	}

	resp, err := s.cli.ContainerExecAttach(ctx, execConfig.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec: %w", err)
	}
	defer resp.Close()

	data, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to read exec output: %w", err)
	}

	return string(data), nil
}
