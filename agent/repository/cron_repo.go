package repository

import (
	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"gorm.io/gorm"
)

// CronRepository provides CRUD operations for CronJob.
type CronRepository struct {
	db *gorm.DB
}

// NewCronRepository creates a new CronRepository.
func NewCronRepository(db *gorm.DB) *CronRepository {
	return &CronRepository{db: db}
}

// Create creates a new CronJob record.
func (r *CronRepository) Create(job *model.CronJob) error {
	return r.db.Create(job).Error
}

// GetByID returns a CronJob by ID.
func (r *CronRepository) GetByID(id uint) (*model.CronJob, error) {
	var job model.CronJob
	if err := r.db.First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// Update updates a CronJob record.
func (r *CronRepository) Update(job *model.CronJob) error {
	return r.db.Save(job).Error
}

// Delete deletes a CronJob record.
func (r *CronRepository) Delete(id uint) error {
	return r.db.Delete(&model.CronJob{}, id).Error
}

// List returns all CronJob records.
func (r *CronRepository) List() ([]model.CronJob, error) {
	var jobs []model.CronJob
	if err := r.db.Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}
