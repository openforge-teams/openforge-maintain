package appstore

// AppRegistry holds the built-in application definitions.
type AppRegistry struct {
	apps map[string]*AppMeta
}

// NewAppRegistry creates a new AppRegistry with built-in applications.
func NewAppRegistry() *AppRegistry {
	r := &AppRegistry{
		apps: make(map[string]*AppMeta),
	}
	r.registerBuiltinApps()
	return r
}

// Get returns an app by key, or nil if not found.
func (r *AppRegistry) Get(key string) *AppMeta {
	return r.apps[key]
}

// GetAll returns all registered apps.
func (r *AppRegistry) GetAll() []AppMeta {
	result := make([]AppMeta, 0, len(r.apps))
	for _, app := range r.apps {
		result = append(result, *app)
	}
	return result
}

// Register adds an app to the registry.
func (r *AppRegistry) Register(app AppMeta) {
	r.apps[app.Key] = &app
}

// registerBuiltinApps registers all built-in application definitions.
func (r *AppRegistry) registerBuiltinApps() {
	// WordPress
	r.Register(AppMeta{
		Key:           "wordpress",
		Name:          "WordPress",
		ModuleName:    "wordpress",
		Version:       "6.4.2",
		Category:      "CMS",
		DockerImage:   "wordpress:6.4.2",
		Icon:          "wordpress",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 80, HostPort: 8081, Protocol: "tcp", Description: "HTTP"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/var/www/html", HostPath: "/opt/wordpress/data", Description: "Website data"},
		},
		EnvVars: []EnvVar{
			{Name: "WORDPRESS_DB_HOST", Default: "mysql", Required: true, Comment: "MySQL host"},
			{Name: "WORDPRESS_DB_USER", Default: "wordpress", Required: true, Comment: "MySQL user"},
			{Name: "WORDPRESS_DB_PASSWORD", Default: "", Required: true, Comment: "MySQL password"},
			{Name: "WORDPRESS_DB_NAME", Default: "wordpress", Required: true, Comment: "MySQL database name"},
		},
	})

	// Halo
	r.Register(AppMeta{
		Key:           "halo",
		Name:          "Halo",
		ModuleName:    "halo",
		Version:       "2.17.0",
		Category:      "CMS",
		DockerImage:   "halohub/halo:2.17.0",
		Icon:          "halo",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 8090, HostPort: 8090, Protocol: "tcp", Description: "HTTP"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/root/.halo2", HostPath: "/opt/halo/data", Description: "Halo data"},
		},
		EnvVars: []EnvVar{},
	})

	// MySQL
	r.Register(AppMeta{
		Key:           "mysql",
		Name:          "MySQL",
		ModuleName:    "mysql",
		Version:       "8.0.36",
		Category:      "Database",
		DockerImage:   "mysql:8.0.36",
		Icon:          "mysql",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 3306, HostPort: 3306, Protocol: "tcp", Description: "MySQL"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/var/lib/mysql", HostPath: "/opt/mysql/data", Description: "MySQL data"},
		},
		EnvVars: []EnvVar{
			{Name: "MYSQL_ROOT_PASSWORD", Default: "", Required: true, Comment: "Root password"},
			{Name: "MYSQL_DATABASE", Default: "", Required: false, Comment: "Initial database"},
		},
	})

	// Redis
	r.Register(AppMeta{
		Key:           "redis",
		Name:          "Redis",
		ModuleName:    "redis",
		Version:       "7.2.4",
		Category:      "Database",
		DockerImage:   "redis:7.2.4-alpine",
		Icon:          "redis",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 6379, HostPort: 6379, Protocol: "tcp", Description: "Redis"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/data", HostPath: "/opt/redis/data", Description: "Redis data"},
		},
		EnvVars: []EnvVar{},
	})

	// PostgreSQL
	r.Register(AppMeta{
		Key:           "postgresql",
		Name:          "PostgreSQL",
		ModuleName:    "postgresql",
		Version:       "16.1",
		Category:      "Database",
		DockerImage:   "postgres:16.1-alpine",
		Icon:          "postgresql",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 5432, HostPort: 5432, Protocol: "tcp", Description: "PostgreSQL"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/var/lib/postgresql/data", HostPath: "/opt/postgresql/data", Description: "PostgreSQL data"},
		},
		EnvVars: []EnvVar{
			{Name: "POSTGRES_USER", Default: "postgres", Required: true, Comment: "Superuser name"},
			{Name: "POSTGRES_PASSWORD", Default: "", Required: true, Comment: "Superuser password"},
		},
	})

	// Nginx
	r.Register(AppMeta{
		Key:           "nginx",
		Name:          "Nginx",
		ModuleName:    "nginx",
		Version:       "1.25.4",
		Category:      "Web Server",
		DockerImage:   "nginx:1.25.4-alpine",
		Icon:          "nginx",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 80, HostPort: 80, Protocol: "tcp", Description: "HTTP"},
			{ContainerPort: 443, HostPort: 443, Protocol: "tcp", Description: "HTTPS"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/etc/nginx/conf.d", HostPath: "/opt/nginx/conf.d", Description: "Nginx configs"},
			{ContainerPath: "/usr/share/nginx/html", HostPath: "/opt/nginx/html", Description: "Website files"},
		},
		EnvVars: []EnvVar{},
	})

	// OpenResty
	r.Register(AppMeta{
		Key:           "openresty",
		Name:          "OpenResty",
		ModuleName:    "openresty",
		Version:       "1.25.3.1",
		Category:      "Web Server",
		DockerImage:   "openresty/openresty:1.25.3.1-alpine",
		Icon:          "openresty",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 80, HostPort: 80, Protocol: "tcp", Description: "HTTP"},
			{ContainerPort: 443, HostPort: 443, Protocol: "tcp", Description: "HTTPS"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/usr/local/openresty/nginx/conf", HostPath: "/opt/openresty/conf", Description: "Configs"},
			{ContainerPath: "/usr/local/openresty/nginx/html", HostPath: "/opt/openresty/html", Description: "Website files"},
		},
		EnvVars: []EnvVar{},
	})

	// MaxKB
	r.Register(AppMeta{
		Key:           "maxkb",
		Name:          "MaxKB",
		ModuleName:    "maxkb",
		Version:       "1.5.1",
		Category:      "AI",
		DockerImage:   "maxkb/maxkb:latest",
		Icon:          "maxkb",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 8080, HostPort: 8080, Protocol: "tcp", Description: "HTTP"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/opt/maxkb/data", HostPath: "/opt/maxkb/data", Description: "MaxKB data"},
		},
		EnvVars: []EnvVar{},
	})

	// n8n
	r.Register(AppMeta{
		Key:           "n8n",
		Name:          "n8n",
		ModuleName:    "n8n",
		Version:       "1.26.2",
		Category:      "Automation",
		DockerImage:   "n8nio/n8n:1.26.2",
		Icon:          "n8n",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 5678, HostPort: 5678, Protocol: "tcp", Description: "Web UI"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/home/node/.n8n", HostPath: "/opt/n8n/data", Description: "n8n data"},
		},
		EnvVars: []EnvVar{
			{Name: "N8N_BASIC_AUTH_ACTIVE", Default: "true", Required: false, Comment: "Enable basic auth"},
			{Name: "N8N_BASIC_AUTH_USER", Default: "admin", Required: false, Comment: "Auth username"},
			{Name: "N8N_BASIC_AUTH_PASSWORD", Default: "", Required: false, Comment: "Auth password"},
		},
	})

	// Ollama
	r.Register(AppMeta{
		Key:           "ollama",
		Name:          "Ollama",
		ModuleName:    "ollama",
		Version:       "0.1.30",
		Category:      "AI",
		DockerImage:   "ollama/ollama:latest",
		Icon:          "ollama",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 11434, HostPort: 11434, Protocol: "tcp", Description: "API"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/root/.ollama", HostPath: "/opt/ollama/data", Description: "Model data"},
		},
		EnvVars: []EnvVar{},
	})

	// MinIO
	r.Register(AppMeta{
		Key:           "minio",
		Name:          "MinIO",
		ModuleName:    "minio",
		Version:       "latest",
		Category:      "Storage",
		DockerImage:   "minio/minio:latest",
		Icon:          "minio",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 9000, HostPort: 9000, Protocol: "tcp", Description: "API"},
			{ContainerPort: 9001, HostPort: 9001, Protocol: "tcp", Description: "Console"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/data", HostPath: "/opt/minio/data", Description: "Storage data"},
		},
		EnvVars: []EnvVar{
			{Name: "MINIO_ROOT_USER", Default: "minioadmin", Required: true, Comment: "Root user"},
			{Name: "MINIO_ROOT_PASSWORD", Default: "minioadmin", Required: true, Comment: "Root password"},
		},
	})

	// Portainer
	r.Register(AppMeta{
		Key:           "portainer",
		Name:          "Portainer",
		ModuleName:    "portainer",
		Version:       "2.19.5",
		Category:      "DevOps",
		DockerImage:   "portainer/portainer-ce:2.19.5-alpine",
		Icon:          "portainer",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 9000, HostPort: 9443, Protocol: "tcp", Description: "Web UI (HTTPS)"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/data", HostPath: "/opt/portainer/data", Description: "Portainer data"},
			{ContainerPath: "/var/run/docker.sock", HostPath: "/var/run/docker.sock", Description: "Docker socket"},
		},
		EnvVars: []EnvVar{},
	})

	// Gitea
	r.Register(AppMeta{
		Key:           "gitea",
		Name:          "Gitea",
		ModuleName:    "gitea",
		Version:       "1.21.4",
		Category:      "DevOps",
		DockerImage:   "gitea/gitea:1.21.4",
		Icon:          "gitea",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 3000, HostPort: 3000, Protocol: "tcp", Description: "HTTP"},
			{ContainerPort: 22, HostPort: 2222, Protocol: "tcp", Description: "SSH"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/data", HostPath: "/opt/gitea/data", Description: "Gitea data"},
		},
		EnvVars: []EnvVar{},
	})

	// Nextcloud
	r.Register(AppMeta{
		Key:           "nextcloud",
		Name:          "Nextcloud",
		ModuleName:    "nextcloud",
		Version:       "28.0.1",
		Category:      "Cloud",
		DockerImage:   "nextcloud:28.0.1-apache",
		Icon:          "nextcloud",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 80, HostPort: 8082, Protocol: "tcp", Description: "HTTP"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/var/www/html", HostPath: "/opt/nextcloud/data", Description: "Nextcloud data"},
		},
		EnvVars: []EnvVar{
			{Name: "NEXTCLOUD_ADMIN_USER", Default: "admin", Required: true, Comment: "Admin username"},
			{Name: "NEXTCLOUD_ADMIN_PASSWORD", Default: "", Required: true, Comment: "Admin password"},
		},
	})

	// Jellyfin
	r.Register(AppMeta{
		Key:           "jellyfin",
		Name:          "Jellyfin",
		ModuleName:    "jellyfin",
		Version:       "10.8.13",
		Category:      "Media",
		DockerImage:   "jellyfin/jellyfin:10.8.13",
		Icon:          "jellyfin",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 8096, HostPort: 8096, Protocol: "tcp", Description: "Web UI"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/config", HostPath: "/opt/jellyfin/config", Description: "Config"},
			{ContainerPath: "/media", HostPath: "/opt/jellyfin/media", Description: "Media library"},
		},
		EnvVars: []EnvVar{},
	})

	// HomeAssistant
	r.Register(AppMeta{
		Key:           "homeassistant",
		Name:          "Home Assistant",
		ModuleName:    "homeassistant",
		Version:       "2024.1.5",
		Category:      "IoT",
		DockerImage:   "homeassistant/home-assistant:2024.1.5",
		Icon:          "homeassistant",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 8123, HostPort: 8123, Protocol: "tcp", Description: "Web UI"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/config", HostPath: "/opt/homeassistant/config", Description: "Config"},
		},
		EnvVars: []EnvVar{},
	})

	// Ghost
	r.Register(AppMeta{
		Key:           "ghost",
		Name:          "Ghost",
		ModuleName:    "ghost",
		Version:       "5.75.0",
		Category:      "CMS",
		DockerImage:   "ghost:5.75.0",
		Icon:          "ghost",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 2368, HostPort: 2368, Protocol: "tcp", Description: "HTTP"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/var/lib/ghost/content", HostPath: "/opt/ghost/content", Description: "Content"},
		},
		EnvVars: []EnvVar{
			{Name: "url", Default: "http://localhost:2368", Required: true, Comment: "Site URL"},
		},
	})

	// MariaDB
	r.Register(AppMeta{
		Key:           "mariadb",
		Name:          "MariaDB",
		ModuleName:    "mariadb",
		Version:       "11.2.3",
		Category:      "Database",
		DockerImage:   "mariadb:11.2.3",
		Icon:          "mariadb",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 3306, HostPort: 3307, Protocol: "tcp", Description: "MariaDB"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/var/lib/mysql", HostPath: "/opt/mariadb/data", Description: "MariaDB data"},
		},
		EnvVars: []EnvVar{
			{Name: "MYSQL_ROOT_PASSWORD", Default: "", Required: true, Comment: "Root password"},
		},
	})

	// MongoDB
	r.Register(AppMeta{
		Key:           "mongodb",
		Name:          "MongoDB",
		ModuleName:    "mongodb",
		Version:       "7.0.5",
		Category:      "Database",
		DockerImage:   "mongo:7.0.5",
		Icon:          "mongodb",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 27017, HostPort: 27017, Protocol: "tcp", Description: "MongoDB"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/data/db", HostPath: "/opt/mongodb/data", Description: "MongoDB data"},
		},
		EnvVars: []EnvVar{
			{Name: "MONGO_INITDB_ROOT_USERNAME", Default: "root", Required: false, Comment: "Root username"},
			{Name: "MONGO_INITDB_ROOT_PASSWORD", Default: "", Required: true, Comment: "Root password"},
		},
	})

	// RabbitMQ
	r.Register(AppMeta{
		Key:           "rabbitmq",
		Name:          "RabbitMQ",
		ModuleName:    "rabbitmq",
		Version:       "3.13.0",
		Category:      "Middleware",
		DockerImage:   "rabbitmq:3.13.0-management-alpine",
		Icon:          "rabbitmq",
		RestartPolicy: "unless-stopped",
		Ports: []PortMapping{
			{ContainerPort: 5672, HostPort: 5672, Protocol: "tcp", Description: "AMQP"},
			{ContainerPort: 15672, HostPort: 15672, Protocol: "tcp", Description: "Management UI"},
		},
		Volumes: []VolumeMapping{
			{ContainerPath: "/var/lib/rabbitmq", HostPath: "/opt/rabbitmq/data", Description: "RabbitMQ data"},
		},
		EnvVars: []EnvVar{
			{Name: "RABBITMQ_DEFAULT_USER", Default: "admin", Required: true, Comment: "Default user"},
			{Name: "RABBITMQ_DEFAULT_PASS", Default: "", Required: true, Comment: "Default password"},
		},
	})
}
