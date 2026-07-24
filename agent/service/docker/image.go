package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// ImageDTO represents a Docker image for API responses.
type ImageDTO struct {
	ID       string   `json:"id"`
	RepoTags []string `json:"repo_tags"`
	Size     int64    `json:"size"`
	Created  int64    `json:"created"`
	Labels   map[string]string `json:"labels"`
}

// ImageService provides Docker image management operations.
type ImageService struct {
	cli *client.Client
}

// NewImageService creates a new ImageService.
func NewImageService(cli *client.Client) *ImageService {
	return &ImageService{cli: cli}
}

// List returns a list of Docker images.
func (s *ImageService) List(ctx context.Context) ([]ImageDTO, error) {
	images, err := s.cli.ImageList(ctx, types.ImageListOptions{All: false})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	result := make([]ImageDTO, 0, len(images))
	for _, img := range images {
		result = append(result, ImageDTO{
			ID:       img.ID,
			RepoTags: img.RepoTags,
			Size:     img.Size,
			Created:  img.Created,
			Labels:   img.Labels,
		})
	}
	return result, nil
}

// Pull pulls a Docker image from a registry.
func (s *ImageService) Pull(ctx context.Context, image string) error {
	resp, err := s.cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer resp.Close()

	_, err = io.ReadAll(resp)
	if err != nil {
		return fmt.Errorf("failed to read pull response: %w", err)
	}
	return nil
}

// Remove removes a Docker image.
func (s *ImageService) Remove(ctx context.Context, image string, force bool) error {
	_, err := s.cli.ImageRemove(ctx, image, types.ImageRemoveOptions{Force: force})
	if err != nil {
		return fmt.Errorf("failed to remove image: %w", err)
	}
	return nil
}

// Tag tags a Docker image.
func (s *ImageService) Tag(ctx context.Context, source, target string) error {
	err := s.cli.ImageTag(ctx, source, target)
	if err != nil {
		return fmt.Errorf("failed to tag image: %w", err)
	}
	return nil
}

// Search searches Docker Hub for images.
func (s *ImageService) Search(ctx context.Context, term string) ([]types.ImageSearchResult, error) {
	results, err := s.cli.ImageSearch(ctx, term, types.ImageSearchOptions{Limit: 25})
	if err != nil {
		return nil, fmt.Errorf("failed to search images: %w", err)
	}
	return results, nil
}

// Prune removes unused images.
func (s *ImageService) Prune(ctx context.Context) (types.ImagesPruneReport, error) {
	report, err := s.cli.ImagesPrune(ctx, filters.Args{})
	if err != nil {
		return types.ImagesPruneReport{}, fmt.Errorf("failed to prune images: %w", err)
	}
	return report, nil
}
