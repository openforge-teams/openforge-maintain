package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// NetworkDTO represents a Docker network for API responses.
type NetworkDTO struct {
	Name       string            `json:"name"`
	ID         string            `json:"id"`
	Driver     string            `json:"driver"`
	Scope      string            `json:"scope"`
	Labels     map[string]string `json:"labels"`
	Containers map[string]struct {
		Name string `json:"name"`
		IPv4 string `json:"ipv4"`
	} `json:"containers"`
}

// NetworkService provides Docker network management operations.
type NetworkService struct {
	cli *client.Client
}

// NewNetworkService creates a new NetworkService.
func NewNetworkService(cli *client.Client) *NetworkService {
	return &NetworkService{cli: cli}
}

// List returns a list of Docker networks.
func (s *NetworkService) List(ctx context.Context) ([]NetworkDTO, error) {
	networks, err := s.cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	result := make([]NetworkDTO, 0, len(networks))
	for _, n := range networks {
		dto := NetworkDTO{
			Name:    n.Name,
			ID:      n.ID,
			Driver:  n.Driver,
			Scope:   n.Scope,
			Labels:  n.Labels,
		}
		if n.Containers != nil {
			dto.Containers = make(map[string]struct {
				Name string `json:"name"`
				IPv4 string `json:"ipv4"`
			})
			for id, c := range n.Containers {
				dto.Containers[id] = struct {
					Name string `json:"name"`
					IPv4 string `json:"ipv4"`
				}{
					Name: c.Name,
				}
				if len(c.IPv4Address) > 0 {
					dto.Containers[id].IPv4 = c.IPv4Address
				}
			}
		}
		result = append(result, dto)
	}
	return result, nil
}

// Create creates a new Docker network.
func (s *NetworkService) Create(ctx context.Context, name, driver string, labels map[string]string) error {
	_, err := s.cli.NetworkCreate(ctx, name, types.NetworkCreate{
		Driver: driver,
		Labels: labels,
	})
	if err != nil {
		return fmt.Errorf("failed to create network: %w", err)
	}
	return nil
}

// Remove removes a Docker network.
func (s *NetworkService) Remove(ctx context.Context, name string) error {
	err := s.cli.NetworkRemove(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to remove network: %w", err)
	}
	return nil
}

// Prune removes unused networks.
func (s *NetworkService) Prune(ctx context.Context) (types.NetworksPruneReport, error) {
	report, err := s.cli.NetworksPrune(ctx, filters.Args{})
	if err != nil {
		return types.NetworksPruneReport{}, fmt.Errorf("failed to prune networks: %w", err)
	}
	return report, nil
}
