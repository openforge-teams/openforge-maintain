package appstore

import (
	"fmt"
	"strings"
	"time"

	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"github.com/openforge-maintain/openforge-maintain/agent/repository"
)

// EnvVar represents an environment variable for an app.
type EnvVar struct {
	Name     string `json:"name" yaml:"name"`
	Default  string `json:"default" yaml:"default"`
	Required bool   `json:"required" yaml:"required"`
	Comment  string `json:"comment" yaml:"comment"`
}

// PortMapping represents a port mapping for an app.
type PortMapping struct {
	ContainerPort int    `json:"container_port" yaml:"container_port"`
	HostPort      int    `json:"host_port" yaml:"host_port"`
	Protocol      string `json:"protocol" yaml:"protocol"`
	Description   string `json:"description" yaml:"description"`
}

// VolumeMapping represents a volume mapping for an app.
type VolumeMapping struct {
	ContainerPath string `json:"container_path" yaml:"container_path"`
	HostPath       string `json:"host_path" yaml:"host_path"`
	Description    string `json:"description" yaml:"description"`
}

// AppMeta represents metadata for an application in the store.
type AppMeta struct {
	Key        string         `json:"key" yaml:"key"`
	Name       string         `json:"name" yaml:"name"`
	ModuleName string         `json:"module_name" yaml:"module_name"`
	Version    string         `json:"version" yaml:"version"`
	Category   string         `json:"category" yaml:"category"`
	Image      string         `json:"image" yaml:"image"`
	DockerImage string       `json:"docker_image" yaml:"docker_image"`
	Icon       string        `json:"icon" yaml:"icon"`
	Ports      []PortMapping `json:"ports" yaml:"ports"`
	Volumes    []VolumeMapping `json:"volumes" yaml:"volumes"`
	EnvVars    []EnvVar      `json:"env_vars" yaml:"env_vars"`
	RestartPolicy string     `json:"restart_policy" yaml:"restart_policy"`
}

// AppStoreService provides application store operations.
type AppStoreService struct {
	repo     *repository.AppRepository
	registry *AppRegistry
}

// NewAppStoreService creates a new AppStoreService.
func NewAppStoreService(repo *repository.AppRepository, registry *AppRegistry) *AppStoreService {
	return &AppStoreService{
		repo:     repo,
		registry: registry,
	}
}

// GetAppList returns the list of available applications from the registry.
func (s *AppStoreService) GetAppList() ([]AppMeta, error) {
	return s.registry.GetAll(), nil
}

// GetAppDetail returns detailed information about a specific application.
func (s *AppStoreService) GetAppDetail(appKey string) (*AppMeta, error) {
	app := s.registry.Get(appKey)
	if app == nil {
		return nil, fmt.Errorf("application not found: %s", appKey)
	}
	return app, nil
}

// Install installs an application with the given parameters.
func (s *AppStoreService) Install(appKey string, params map[string]string) error {
	app := s.registry.Get(appKey)
	if app == nil {
		return fmt.Errorf("application not found: %s", appKey)
	}

	// Validate required parameters
	for _, env := range app.EnvVars {
		if env.Required {
			if _, ok := params[env.Name]; !ok && env.Default == "" {
				return fmt.Errorf("required parameter missing: %s", env.Name)
			}
		}
	}

	// Generate compose file
	composeContent := s.generateComposeFile(app, params)

	// Generate env config
	envContent := s.generateEnvConfig(app, params)

	// Determine primary port
	var port int
	if len(app.Ports) > 0 {
		port = app.Ports[0].HostPort
	}

	install := &model.AppInstall{
		AppKey:      appKey,
		Name:        app.Name,
		Version:     app.Version,
		Status:      "installing",
		ComposeFile: composeContent,
		EnvConfig:   envContent,
		Port:        port,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(install); err != nil {
		return fmt.Errorf("failed to save install record: %w", err)
	}

	// In a real implementation, we would run docker compose up here
	install.Status = "running"
	install.UpdatedAt = time.Now()
	return s.repo.Update(install)
}

// Uninstall removes an installed application.
func (s *AppStoreService) Uninstall(id uint) error {
	install, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get install record: %w", err)
	}

	// In a real implementation, we would run docker compose down here
	install.Status = "uninstalled"
	install.UpdatedAt = time.Now()

	if err := s.repo.Update(install); err != nil {
		return fmt.Errorf("failed to update install record: %w", err)
	}

	return s.repo.Delete(id)
}

// Upgrade upgrades an installed application.
func (s *AppStoreService) Upgrade(id uint) error {
	install, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get install record: %w", err)
	}

	app := s.registry.Get(install.AppKey)
	if app == nil {
		return fmt.Errorf("application not found in registry: %s", install.AppKey)
	}

	install.Version = app.Version
	install.Status = "upgrading"
	install.UpdatedAt = time.Now()

	if err := s.repo.Update(install); err != nil {
		return fmt.Errorf("failed to update install record: %w", err)
	}

	// In a real implementation, we would pull new images and recreate containers
	install.Status = "running"
	install.UpdatedAt = time.Now()
	return s.repo.Update(install)
}

// GetInstalledApps returns all installed applications.
func (s *AppStoreService) GetInstalledApps() ([]model.AppInstall, error) {
	return s.repo.List()
}

// generateComposeFile generates a docker-compose.yml content for the app.
func (s *AppStoreService) generateComposeFile(app *AppMeta, params map[string]string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("version: '3.8'\nservices:\n  %s:\n", app.ModuleName))
	sb.WriteString(fmt.Sprintf("    image: %s\n", app.DockerImage))
	sb.WriteString(fmt.Sprintf("    restart: %s\n", app.RestartPolicy))

	// Ports
	if len(app.Ports) > 0 {
		sb.WriteString("    ports:\n")
		for _, p := range app.Ports {
			sb.WriteString(fmt.Sprintf("      - \"%d:%d\"\n", p.HostPort, p.ContainerPort))
		}
	}

	// Volumes
	if len(app.Volumes) > 0 {
		sb.WriteString("    volumes:\n")
		for _, v := range app.Volumes {
			hostPath := v.HostPath
			if val, ok := params[v.ContainerPath]; ok {
				hostPath = val
			}
			sb.WriteString(fmt.Sprintf("      - \"%s:%s\"\n", hostPath, v.ContainerPath))
		}
	}

	// Environment variables
	if len(app.EnvVars) > 0 {
		sb.WriteString("    environment:\n")
		for _, env := range app.EnvVars {
			value := env.Default
			if val, ok := params[env.Name]; ok {
				value = val
			}
			if value != "" {
				sb.WriteString(fmt.Sprintf("      - \"%s=%s\"\n", env.Name, value))
			}
		}
	}

	return sb.String()
}

// generateEnvConfig generates an environment configuration string.
func (s *AppStoreService) generateEnvConfig(app *AppMeta, params map[string]string) string {
	var sb strings.Builder
	for _, env := range app.EnvVars {
		value := env.Default
		if val, ok := params[env.Name]; ok {
			value = val
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", env.Name, value))
	}
	return sb.String()
}
