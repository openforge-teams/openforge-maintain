package database

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// DBInfo represents a database instance.
type DBInfo struct {
	Name     string `json:"name"`
	Charset  string `json:"charset"`
	Size     string `json:"size"`
	Tables   int    `json:"tables"`
}

// MySQLService provides MySQL database management operations via Docker.
type MySQLService struct {
	containerName string
	rootPassword  string
}

// NewMySQLService creates a new MySQLService.
func NewMySQLService(containerName, rootPassword string) *MySQLService {
	return &MySQLService{
		containerName: containerName,
		rootPassword:  rootPassword,
	}
}

// execMySQL executes a MySQL command inside the container.
func (s *MySQLService) execMySQL(args ...string) (string, error) {
	fullArgs := []string{"exec", s.containerName, "mysql", "-uroot", "-p" + s.rootPassword}
	fullArgs = append(fullArgs, args...)
	cmd := exec.Command("docker", fullArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("mysql exec failed: %s: %w", string(output), err)
	}
	return string(output), nil
}

// List returns a list of databases.
func (s *MySQLService) List() ([]DBInfo, error) {
	output, err := s.execMySQL("-e", "SHOW DATABASES;")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var dbs []DBInfo
	for _, line := range lines[1:] { // skip header
		line = strings.TrimSpace(line)
		if line == "" || line == "information_schema" || line == "mysql" || line == "performance_schema" || line == "sys" {
			continue
		}
		dbs = append(dbs, DBInfo{
			Name: line,
		})
	}
	return dbs, nil
}

// CreateDB creates a new database.
func (s *MySQLService) CreateDB(name, charset string) error {
	if charset == "" {
		charset = "utf8mb4"
	}
	_, err := s.execMySQL("-e", fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET %s COLLATE %s_general_ci;", name, charset, charset))
	return err
}

// DeleteDB drops a database.
func (s *MySQLService) DeleteDB(name string) error {
	_, err := s.execMySQL("-e", fmt.Sprintf("DROP DATABASE IF EXISTS `%s`;", name))
	return err
}

// CreateUser creates a database user.
func (s *MySQLService) CreateUser(db, user, password string) error {
	_, err := s.execMySQL("-e", fmt.Sprintf(
		"CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY '%s'; GRANT ALL PRIVILEGES ON `%s`.* TO '%s'@'%%'; FLUSH PRIVILEGES;",
		user, password, db, user,
	))
	return err
}

// Backup creates a database backup using mysqldump.
func (s *MySQLService) Backup(db, dest string) error {
	args := []string{"exec", s.containerName, "mysqldump", "-uroot", "-p" + s.rootPassword, db}
	cmd := exec.Command("docker", args...)

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("mysqldump failed: %w", err)
	}

	if err := exec.Command("cp", "/dev/stdin", dest).Run(); err != nil {
		// Write directly
		if err := writeFile(dest, output); err != nil {
			return fmt.Errorf("failed to write backup file: %w", err)
		}
	}
	return nil
}

// Restore restores a database from a backup file.
func (s *MySQLService) Restore(db, src string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	args := []string{"exec", "-i", s.containerName, "mysql", "-uroot", "-p" + s.rootPassword, db}
	cmd := exec.Command("docker", args...)
	cmd.Stdin = strings.NewReader(string(data))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mysql restore failed: %s: %w", string(output), err)
	}
	return nil
}

// writeFile writes data to a file path.
func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
