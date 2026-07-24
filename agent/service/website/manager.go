package website

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"github.com/openforge-maintain/openforge-maintain/agent/repository"
	"gorm.io/gorm"
)

const nginxSitesAvailable = "/etc/nginx/sites-available"
const nginxSitesEnabled = "/etc/nginx/sites-enabled"

// WebsiteService provides website management operations.
type WebsiteService struct {
	repo      *repository.WebsiteRepository
	certRepo  *repository.SSLRepository
	db        *gorm.DB
}

// NewWebsiteService creates a new WebsiteService.
func NewWebsiteService(repo *repository.WebsiteRepository, certRepo *repository.SSLRepository, db *gorm.DB) *WebsiteService {
	return &WebsiteService{
		repo:     repo,
		certRepo: certRepo,
		db:       db,
	}
}

// Create creates a new website and generates Nginx configuration.
func (s *WebsiteService) Create(website *model.Website) error {
	website.CreatedAt = time.Now()
	website.UpdatedAt = time.Now()

	if err := s.repo.Create(website); err != nil {
		return fmt.Errorf("failed to create website: %w", err)
	}

	config, err := s.GenerateNginxConfig(website)
	if err != nil {
		return fmt.Errorf("failed to generate nginx config: %w", err)
	}

	if err := s.writeNginxConfig(website.PrimaryDomain, config); err != nil {
		return fmt.Errorf("failed to write nginx config: %w", err)
	}

	if err := s.reloadNginx(); err != nil {
		return fmt.Errorf("failed to reload nginx: %w", err)
	}

	return nil
}

// Update updates a website and regenerates Nginx configuration.
func (s *WebsiteService) Update(website *model.Website) error {
	website.UpdatedAt = time.Now()

	if err := s.repo.Update(website); err != nil {
		return fmt.Errorf("failed to update website: %w", err)
	}

	config, err := s.GenerateNginxConfig(website)
	if err != nil {
		return fmt.Errorf("failed to generate nginx config: %w", err)
	}

	if err := s.writeNginxConfig(website.PrimaryDomain, config); err != nil {
		return fmt.Errorf("failed to write nginx config: %w", err)
	}

	if err := s.reloadNginx(); err != nil {
		return fmt.Errorf("failed to reload nginx: %w", err)
	}

	return nil
}

// Delete removes a website and its Nginx configuration.
func (s *WebsiteService) Delete(id uint) error {
	website, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get website: %w", err)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete website: %w", err)
	}

	configPath := filepath.Join(nginxSitesAvailable, website.PrimaryDomain)
	os.Remove(configPath)

	enabledPath := filepath.Join(nginxSitesEnabled, website.PrimaryDomain)
	os.Remove(enabledPath)

	if err := s.reloadNginx(); err != nil {
		return fmt.Errorf("failed to reload nginx: %w", err)
	}

	return nil
}

// Get returns a website by ID.
func (s *WebsiteService) Get(id uint) (*model.Website, error) {
	return s.repo.GetByID(id)
}

// List returns a paginated list of websites.
func (s *WebsiteService) List(page, size int) ([]model.Website, int64, error) {
	return s.repo.List(page, size)
}

// EnableSSL enables SSL for a website.
func (s *WebsiteService) EnableSSL(websiteID uint, certID uint) error {
	website, err := s.repo.GetByID(websiteID)
	if err != nil {
		return fmt.Errorf("failed to get website: %w", err)
	}

	cert, err := s.certRepo.GetByID(certID)
	if err != nil {
		return fmt.Errorf("failed to get ssl cert: %w", err)
	}

	website.SSLStatus = "enabled"
	website.CertID = certID
	website.UpdatedAt = time.Now()

	if err := s.repo.Update(website); err != nil {
		return fmt.Errorf("failed to update website: %w", err)
	}

	website.CertID = cert.ID // ensure cert is accessible for config generation
	config, err := s.GenerateNginxConfig(website)
	if err != nil {
		return fmt.Errorf("failed to generate nginx config: %w", err)
	}

	// Add cert paths to the config
	sslConfig := fmt.Sprintf(`
    ssl_certificate %s;
    ssl_certificate_key %s;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
`, cert.CertPath, cert.KeyPath)

	config = config + sslConfig

	if err := s.writeNginxConfig(website.PrimaryDomain, config); err != nil {
		return fmt.Errorf("failed to write nginx config: %w", err)
	}

	if err := s.reloadNginx(); err != nil {
		return fmt.Errorf("failed to reload nginx: %w", err)
	}

	return nil
}

// DisableSSL disables SSL for a website.
func (s *WebsiteService) DisableSSL(websiteID uint) error {
	website, err := s.repo.GetByID(websiteID)
	if err != nil {
		return fmt.Errorf("failed to get website: %w", err)
	}

	website.SSLStatus = "disabled"
	website.CertID = 0
	website.UpdatedAt = time.Now()

	if err := s.repo.Update(website); err != nil {
		return fmt.Errorf("failed to update website: %w", err)
	}

	config, err := s.GenerateNginxConfig(website)
	if err != nil {
		return fmt.Errorf("failed to generate nginx config: %w", err)
	}

	if err := s.writeNginxConfig(website.PrimaryDomain, config); err != nil {
		return fmt.Errorf("failed to write nginx config: %w", err)
	}

	if err := s.reloadNginx(); err != nil {
		return fmt.Errorf("failed to reload nginx: %w", err)
	}

	return nil
}

// GenerateNginxConfig generates an Nginx configuration template for a website.
func (s *WebsiteService) GenerateNginxConfig(website *model.Website) (string, error) {
	const tpl = `server {
    listen 80;
    server_name {{.Domain}};
{{- if eq .SSLStatus "enabled"}}
    listen 443 ssl http2;
{{- end}}

    access_log /var/log/nginx/{{.Domain}}.access.log;
    error_log  /var/log/nginx/{{.Domain}}.error.log;

{{- if eq .Type "static"}}
    root {{.RootDir}};
    index index.html index.htm;

    location / {
        try_files $uri $uri/ =404;
    }
{{- end}}

{{- if eq .Type "proxy"}}
    location / {
        proxy_pass {{.ProxyTarget}};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
{{- end}}

{{- if eq .Type "php"}}
    root {{.RootDir}};
    index index.php index.html;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        fastcgi_pass unix:/var/run/php/php-fpm.sock;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
{{- end}}
}`

	t, err := template.New("nginx").Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf strings.Builder
	if err := t.Execute(&buf, website); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// writeNginxConfig writes the Nginx configuration file.
func (s *WebsiteService) writeNginxConfig(domain, config string) error {
	if err := os.MkdirAll(nginxSitesAvailable, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(nginxSitesEnabled, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(nginxSitesAvailable, domain)
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return err
	}

	enabledPath := filepath.Join(nginxSitesEnabled, domain)
	os.Remove(enabledPath)
	return os.Symlink(configPath, enabledPath)
}

// reloadNginx reloads the Nginx service.
func (s *WebsiteService) reloadNginx() error {
	cmd := exec.Command("nginx", "-t")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("nginx config test failed: %s: %w", string(output), err)
	}

	cmd = exec.Command("nginx", "-s", "reload")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("nginx reload failed: %s: %w", string(output), err)
	}
	return nil
}
