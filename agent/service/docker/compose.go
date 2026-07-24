package docker

import (
	"fmt"
	"os"
	"os/exec"
)

// ComposeService provides Docker Compose operations.
type ComposeService struct {
	dockerHost string
}

// NewComposeService creates a new ComposeService.
func NewComposeService(dockerHost string) *ComposeService {
	return &ComposeService{dockerHost: dockerHost}
}

// Up runs docker compose up -d.
func (s *ComposeService) Up(composePath string) error {
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return fmt.Errorf("compose file not found: %s", composePath)
	}

	args := []string{"compose", "-f", composePath, "up", "-d"}
	cmd := exec.Command("docker", args...)
	cmd.Env = append(os.Environ(), "DOCKER_HOST="+s.dockerHost)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker compose up failed: %s: %w", string(output), err)
	}
	return nil
}

// Down runs docker compose down.
func (s *ComposeService) Down(composePath string) error {
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return fmt.Errorf("compose file not found: %s", composePath)
	}

	args := []string{"compose", "-f", composePath, "down"}
	cmd := exec.Command("docker", args...)
	cmd.Env = append(os.Environ(), "DOCKER_HOST="+s.dockerHost)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker compose down failed: %s: %w", string(output), err)
	}
	return nil
}

// List runs docker compose ps.
func (s *ComposeService) List(composePath string) (string, error) {
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return "", fmt.Errorf("compose file not found: %s", composePath)
	}

	args := []string{"compose", "-f", composePath, "ps", "--format", "json"}
	cmd := exec.Command("docker", args...)
	cmd.Env = append(os.Environ(), "DOCKER_HOST="+s.dockerHost)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker compose ps failed: %s: %w", string(output), err)
	}
	return string(output), nil
}

// PullImages runs docker compose pull.
func (s *ComposeService) PullImages(composePath string) error {
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return fmt.Errorf("compose file not found: %s", composePath)
	}

	args := []string{"compose", "-f", composePath, "pull"}
	cmd := exec.Command("docker", args...)
	cmd.Env = append(os.Environ(), "DOCKER_HOST="+s.dockerHost)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker compose pull failed: %s: %w", string(output), err)
	}
	return nil
}

// Restart runs docker compose restart.
func (s *ComposeService) Restart(composePath string) error {
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return fmt.Errorf("compose file not found: %s", composePath)
	}

	args := []string{"compose", "-f", composePath, "restart"}
	cmd := exec.Command("docker", args...)
	cmd.Env = append(os.Environ(), "DOCKER_HOST="+s.dockerHost)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker compose restart failed: %s: %w", string(output), err)
	}
	return nil
}

// Config validates and returns the parsed compose configuration.
func (s *ComposeService) Config(composePath string) (string, error) {
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return "", fmt.Errorf("compose file not found: %s", composePath)
	}

	args := []string{"compose", "-f", composePath, "config"}
	cmd := exec.Command("docker", args...)
	cmd.Env = append(os.Environ(), "DOCKER_HOST="+s.dockerHost)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker compose config failed: %s: %w", string(output), err)
	}
	return string(output), nil
}

// Logs runs docker compose logs.
func (s *ComposeService) Logs(composePath string, tail string) (string, error) {
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return "", fmt.Errorf("compose file not found: %s", composePath)
	}

	args := []string{"compose", "-f", composePath, "logs"}
	if tail != "" && tail != "all" {
		args = append(args, "--tail", tail)
	}
	cmd := exec.Command("docker", args...)
	cmd.Env = append(os.Environ(), "DOCKER_HOST="+s.dockerHost)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker compose logs failed: %s: %w", string(output), err)
	}
	return string(output), nil
}
