package database

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PostgresService provides PostgreSQL database management operations via Docker.
type PostgresService struct {
	containerName string
	user          string
	password      string
}

// NewPostgresService creates a new PostgresService.
func NewPostgresService(containerName, user, password string) *PostgresService {
	return &PostgresService{
		containerName: containerName,
		user:          user,
		password:      password,
	}
}

// execPostgres executes a psql command inside the container.
func (s *PostgresService) execPostgres(args ...string) (string, error) {
	fullArgs := []string{"exec", s.containerName, "psql", "-U", s.user}
	fullArgs = append(fullArgs, args...)
	cmd := exec.Command("docker", fullArgs...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", s.password))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("postgres exec failed: %s: %w", string(output), err)
	}
	return string(output), nil
}

// List returns a list of databases.
func (s *PostgresService) List() ([]DBInfo, error) {
	output, err := s.execPostgres("-c", "\\l")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var dbs []DBInfo
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "---") || strings.HasPrefix(line, "List") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			name := fields[0]
			if name == "template0" || name == "template1" || name == "postgres" {
				continue
			}
			dbs = append(dbs, DBInfo{
				Name:    name,
				Charset: fields[2],
			})
		}
	}
	return dbs, nil
}

// CreateDB creates a new database.
func (s *PostgresService) CreateDB(name, charset string) error {
	_, err := s.execPostgres("-c", fmt.Sprintf("CREATE DATABASE \"%s\" ENCODING '%s';", name, charset))
	return err
}

// DeleteDB drops a database.
func (s *PostgresService) DeleteDB(name string) error {
	_, err := s.execPostgres("-c", fmt.Sprintf("DROP DATABASE IF EXISTS \"%s\";", name))
	return err
}

// CreateUser creates a database user.
func (s *PostgresService) CreateUser(db, user, password string) error {
	_, err := s.execPostgres("-c", fmt.Sprintf(
		"CREATE USER \"%s\" WITH PASSWORD '%s'; GRANT ALL PRIVILEGES ON DATABASE \"%s\" TO \"%s\";",
		user, password, db, user,
	))
	return err
}

// Backup creates a database backup using pg_dump.
func (s *PostgresService) Backup(db, dest string) error {
	args := []string{"exec", s.containerName, "pg_dump", "-U", s.user, db}
	cmd := exec.Command("docker", args...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", s.password))

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("pg_dump failed: %w", err)
	}

	if err := os.WriteFile(dest, output, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}
	return nil
}

// Restore restores a database from a backup file.
func (s *PostgresService) Restore(db, src string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	args := []string{"exec", "-i", s.containerName, "psql", "-U", s.user, db}
	cmd := exec.Command("docker", args...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", s.password))
	cmd.Stdin = strings.NewReader(string(data))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("postgres restore failed: %s: %w", string(output), err)
	}
	return nil
}
