package database

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// RedisInfo represents Redis server information.
type RedisInfo struct {
	Version      string            `json:"version"`
	Mode         string            `json:"mode"`
	OS           string            `json:"os"`
	TCP          string            `json:"tcp_port"`
	Uptime       string            `json:"uptime"`
	Clients      int               `json:"clients"`
	UsedMemory   string            `json:"used_memory"`
	MaxMemory    string            `json:"max_memory"`
	TotalCommands int64             `json:"total_commands"`
	KeySpace     map[string]string `json:"keyspace"`
}

// RedisDB represents a Redis database.
type RedisDB struct {
	Index    int   `json:"index"`
	Keys     int64 `json:"keys"`
	Expires  int64 `json:"expires"`
	AvgTTL   int64 `json:"avg_ttl"`
}

// RedisService provides Redis management operations via Docker.
type RedisService struct {
	containerName string
}

// NewRedisService creates a new RedisService.
func NewRedisService(containerName string) *RedisService {
	return &RedisService{containerName: containerName}
}

// execRedis executes a redis-cli command inside the container.
func (s *RedisService) execRedis(args ...string) (string, error) {
	fullArgs := []string{"exec", s.containerName, "redis-cli"}
	fullArgs = append(fullArgs, args...)
	cmd := exec.Command("docker", fullArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("redis-cli exec failed: %s: %w", string(output), err)
	}
	return string(output), nil
}

// GetInfo returns Redis server information.
func (s *RedisService) GetInfo() (*RedisInfo, error) {
	output, err := s.execRedis("INFO", "all")
	if err != nil {
		return nil, err
	}

	info := &RedisInfo{KeySpace: make(map[string]string)}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "redis_version":
			info.Version = value
		case "redis_mode":
			info.Mode = value
		case "os":
			info.OS = value
		case "tcp_port":
			info.TCP = value
		case "uptime_in_days":
			info.Uptime = value + " days"
		case "connected_clients":
			info.Clients, _ = strconv.Atoi(value)
		case "used_memory_human":
			info.UsedMemory = value
		case "maxmemory_human":
			info.MaxMemory = value
		case "total_commands_processed":
			info.TotalCommands, _ = strconv.ParseInt(value, 10, 64)
		}

		if strings.HasPrefix(key, "db") {
			info.KeySpace[key] = value
		}
	}

	return info, nil
}

// SetConfig sets a Redis configuration parameter.
func (s *RedisService) SetConfig(key, value string) error {
	_, err := s.execRedis("CONFIG", "SET", key, value)
	return err
}

// GetDatabases returns a list of Redis databases with key counts.
func (s *RedisService) GetDatabases() ([]RedisDB, error) {
	output, err := s.execRedis("INFO", "keyspace")
	if err != nil {
		return nil, err
	}

	var dbs []RedisDB
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "db") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		var index int
		if _, err := fmt.Sscanf(parts[0], "db%d", &index); err != nil {
			continue
		}

		db := RedisDB{Index: index}
		fields := strings.Split(parts[1], ",")
		for _, field := range fields {
			field = strings.TrimSpace(field)
			kv := strings.SplitN(field, "=", 2)
			if len(kv) != 2 {
				continue
			}
			switch kv[0] {
			case "keys":
				db.Keys, _ = strconv.ParseInt(kv[1], 10, 64)
			case "expires":
				db.Expires, _ = strconv.ParseInt(kv[1], 10, 64)
			case "avg_ttl":
				db.AvgTTL, _ = strconv.ParseInt(kv[1], 10, 64)
			}
		}
		dbs = append(dbs, db)
	}
	return dbs, nil
}

// FlushDB flushes all keys in a specific database.
func (s *RedisService) FlushDB(dbIndex int) error {
	_, err := s.execRedis("-n", strconv.Itoa(dbIndex), "FLUSHDB")
	return err
}

// GetConfig returns the value of a Redis configuration parameter.
func (s *RedisService) GetConfig(key string) (string, error) {
	return s.execRedis("CONFIG", "GET", key)
}
