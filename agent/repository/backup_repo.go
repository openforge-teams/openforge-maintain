package repository

import (
	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"gorm.io/gorm"
)

// BackupRepository provides CRUD operations for BackupTask.
type BackupRepository struct {
	db *gorm.DB
}

// NewBackupRepository creates a new BackupRepository.
func NewBackupRepository(db *gorm.DB) *BackupRepository {
	return &BackupRepository{db: db}
}

// Create creates a new BackupTask record.
func (r *BackupRepository) Create(task *model.BackupTask) error {
	return r.db.Create(task).Error
}

// GetByID returns a BackupTask by ID.
func (r *BackupRepository) GetByID(id uint) (*model.BackupTask, error) {
	var task model.BackupTask
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// Update updates a BackupTask record.
func (r *BackupRepository) Update(task *model.BackupTask) error {
	return r.db.Save(task).Error
}

// Delete deletes a BackupTask record.
func (r *BackupRepository) Delete(id uint) error {
	return r.db.Delete(&model.BackupTask{}, id).Error
}

// List returns all BackupTask records.
func (r *BackupRepository) List() ([]model.BackupTask, error) {
	var tasks []model.BackupTask
	if err := r.db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// ListByTargetType returns BackupTasks filtered by target type.
func (r *BackupRepository) ListByTargetType(targetType string) ([]model.BackupTask, error) {
	var tasks []model.BackupTask
	if err := r.db.Where("target_type = ?", targetType).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
