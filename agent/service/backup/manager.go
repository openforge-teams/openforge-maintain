package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"github.com/openforge-maintain/openforge-maintain/agent/repository"
)

const backupBaseDir = "/opt/backups"

// BackupRecord represents a backup file record.
type BackupRecord struct {
	ID       uint      `json:"id"`
	TaskID   uint      `json:"task_id"`
	FileName string    `json:"file_name"`
	FileSize int64     `json:"file_size"`
	Created  time.Time `json:"created"`
}

// BackupService provides backup and restore operations.
type BackupService struct {
	repo     *repository.BackupRepository
	baseDir  string
}

// NewBackupService creates a new BackupService.
func NewBackupService(repo *repository.BackupRepository) *BackupService {
	return &BackupService{
		repo:    repo,
		baseDir: backupBaseDir,
	}
}

// CreateBackup creates a backup based on the task configuration.
func (s *BackupService) CreateBackup(task *model.BackupTask) error {
	backupDir := filepath.Join(s.baseDir, fmt.Sprintf("task_%d", task.ID))
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("%s_%s.tar.gz", task.Name, timestamp)
	backupPath := filepath.Join(backupDir, fileName)

	// In a real implementation, this would perform the actual backup based on TargetType
	// For now, create a placeholder file
	f, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	f.Close()

	task.Status = "completed"
	task.UpdatedAt = time.Now()
	return s.repo.Update(task)
}

// Restore restores from a backup file.
func (s *BackupService) Restore(taskID uint, backupFile string) error {
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get backup task: %w", err)
	}

	backupDir := filepath.Join(s.baseDir, fmt.Sprintf("task_%d", taskID))
	backupPath := filepath.Join(backupDir, backupFile)

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupFile)
	}

	// In a real implementation, this would perform the actual restore based on TargetType
	task.Status = "restoring"
	task.UpdatedAt = time.Now()
	s.repo.Update(task)

	task.Status = "completed"
	task.UpdatedAt = time.Now()
	return s.repo.Update(task)
}

// ListBackups lists all backup files for a task.
func (s *BackupService) ListBackups(taskID uint) ([]BackupRecord, error) {
	backupDir := filepath.Join(s.baseDir, fmt.Sprintf("task_%d", taskID))

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []BackupRecord{}, nil
		}
		return nil, fmt.Errorf("failed to list backup directory: %w", err)
	}

	var records []BackupRecord
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}

		records = append(records, BackupRecord{
			TaskID:   taskID,
			FileName: entry.Name(),
			FileSize: info.Size(),
			Created:  info.ModTime(),
		})
	}

	// Sort by creation time descending
	sort.Slice(records, func(i, j int) bool {
		return records[i].Created.After(records[j].Created)
	})

	return records, nil
}

// DeleteBackup deletes a specific backup file.
func (s *BackupService) DeleteBackup(taskID uint, backupFile string) error {
	backupDir := filepath.Join(s.baseDir, fmt.Sprintf("task_%d", taskID))
	backupPath := filepath.Join(backupDir, backupFile)

	return os.Remove(backupPath)
}

// ScheduleBackup sets up a scheduled backup task.
func (s *BackupService) ScheduleBackup() error {
	tasks, err := s.repo.List()
	if err != nil {
		return fmt.Errorf("failed to list backup tasks: %w", err)
	}

	for _, task := range tasks {
		if task.Schedule == "" {
			continue
		}
		// In a real implementation, this would register the backup with the cron scheduler
		_ = task.Schedule
	}

	return nil
}

// CleanOldBackups removes old backup files based on retention policy.
func (s *BackupService) CleanOldBackups(taskID uint) error {
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		return fmt.Errorf("failed to get backup task: %w", err)
	}

	if task.Retention <= 0 {
		return nil // No retention limit
	}

	records, err := s.ListBackups(taskID)
	if err != nil {
		return err
	}

	// Keep only the most recent backups up to retention limit
	if len(records) > task.Retention {
		for _, record := range records[task.Retention:] {
			_ = s.DeleteBackup(taskID, record.FileName)
		}
	}

	return nil
}
