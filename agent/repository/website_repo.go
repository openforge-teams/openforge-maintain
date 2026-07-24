package repository

import (
	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"gorm.io/gorm"
)

// WebsiteRepository provides CRUD operations for Website.
type WebsiteRepository struct {
	db *gorm.DB
}

// NewWebsiteRepository creates a new WebsiteRepository.
func NewWebsiteRepository(db *gorm.DB) *WebsiteRepository {
	return &WebsiteRepository{db: db}
}

// Create creates a new Website record.
func (r *WebsiteRepository) Create(website *model.Website) error {
	return r.db.Create(website).Error
}

// GetByID returns a Website by ID.
func (r *WebsiteRepository) GetByID(id uint) (*model.Website, error) {
	var website model.Website
	if err := r.db.First(&website, id).Error; err != nil {
		return nil, err
	}
	return &website, nil
}

// Update updates a Website record.
func (r *WebsiteRepository) Update(website *model.Website) error {
	return r.db.Save(website).Error
}

// Delete deletes a Website record.
func (r *WebsiteRepository) Delete(id uint) error {
	return r.db.Delete(&model.Website{}, id).Error
}

// List returns a paginated list of Websites.
func (r *WebsiteRepository) List(page, size int) ([]model.Website, int64, error) {
	var websites []model.Website
	var total int64

	r.db.Model(&model.Website{}).Count(&total)

	offset := (page - 1) * size
	if err := r.db.Offset(offset).Limit(size).Order("id desc").Find(&websites).Error; err != nil {
		return nil, 0, err
	}
	return websites, total, nil
}

// GetByDomain returns a Website by primary domain.
func (r *WebsiteRepository) GetByDomain(domain string) (*model.Website, error) {
	var website model.Website
	if err := r.db.Where("primary_domain = ?", domain).First(&website).Error; err != nil {
		return nil, err
	}
	return &website, nil
}
