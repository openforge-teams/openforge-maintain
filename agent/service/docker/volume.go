package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// VolumeDTO represents a Docker volume for API responses.
type VolumeDTO struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	Labels     map[string]string `json:"labels"`
	Size       int64             `json:"size"`
}

// VolumeService provides Docker volume management operations.
type VolumeService struct {
	cli *client.Client
}

// NewVolumeService creates a new VolumeService.
func NewVolumeService(cli *client.Client) *VolumeService {
	return &VolumeService{cli: cli}
}

// List returns a list of Docker volumes.
func (s *VolumeService) List(ctx context.Context) ([]VolumeDTO, error) {
	volumes, err := s.cli.VolumeList(ctx, types.VolumeListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	result := make([]VolumeDTO, 0, len(volumes.Volumes))
	for _, v := range volumes.Volumes {
		var size int64
		if v.UsageData != nil {
			size = v.UsageData.Size
		}
		result = append(result, VolumeDTO{
			Name:       v.Name,
			Driver:     v.Driver,
			Mountpoint: v.Mountpoint,
			Labels:     v.Labels,
			Size:       size,
		})
	}
	return result, nil
}

// Create creates a new Docker volume.
func (s *VolumeService) Create(ctx context.Context, name, driver string, labels map[string]string) error {
	_, err := s.cli.VolumeCreate(ctx, types.VolumeCreateBody{
		Name:   name,
		Driver: driver,
		Labels: labels,
	})
	if err != nil {
		return fmt.Errorf("failed to create volume: %w", err)
	}
	return nil
}

// Remove removes a Docker volume.
func (s *VolumeService) Remove(ctx context.Context, name string, force bool) error {
	err := s.cli.VolumeRemove(ctx, name, force)
	if err != nil {
		return fmt.Errorf("failed to remove volume: %w", err)
	}
	return nil
}

// Prune removes unused volumes.
func (s *VolumeService) Prune(ctx context.Context) (types.VolumesPruneReport, error) {
	report, err := s.cli.VolumesPrune(ctx, filters.Args{})
	if err != nil {
		return types.VolumesPruneReport{}, fmt.Errorf("failed to prune volumes: %w", err)
	}
	return report, nil
}
