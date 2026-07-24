package repository

import (
	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"gorm.io/gorm"
)

// AppRepository provides CRUD operations for AppInstall.
type AppRepository struct {
	db *gorm.DB
}

// NewAppRepository creates a new AppRepository.
func NewAppRepository(db *gorm.DB) *AppRepository {
	return &AppRepository{db: db}
}

// Create creates a new AppInstall record.
func (r *AppRepository) Create(app *model.AppInstall) error {
	return r.db.Create(app).Error
}

// GetByID returns an AppInstall by ID.
func (r *AppRepository) GetByID(id uint) (*model.AppInstall, error) {
	var app model.AppInstall
	if err := r.db.First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

// Update updates an AppInstall record.
func (r *AppRepository) Update(app *model.AppInstall) error {
	return r.db.Save(app).Error
}

// Delete deletes an AppInstall record.
func (r *AppRepository) Delete(id uint) error {
	return r.db.Delete(&model.AppInstall{}, id).Error
}

// List returns all AppInstall records.
func (r *AppRepository) List() ([]model.AppInstall, error) {
	var apps []model.AppInstall
	if err := r.db.Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

// GetByAppKey returns an AppInstall by AppKey.
func (r *AppRepository) GetByAppKey(appKey string) (*model.AppInstall, error) {
	var app model.AppInstall
	if err := r.db.Where("app_key = ?", appKey).First(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}
