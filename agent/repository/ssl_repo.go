package repository

import (
	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"gorm.io/gorm"
)

// SSLRepository provides CRUD operations for SSLCert.
type SSLRepository struct {
	db *gorm.DB
}

// NewSSLRepository creates a new SSLRepository.
func NewSSLRepository(db *gorm.DB) *SSLRepository {
	return &SSLRepository{db: db}
}

// Create creates a new SSLCert record.
func (r *SSLRepository) Create(cert *model.SSLCert) error {
	return r.db.Create(cert).Error
}

// GetByID returns an SSLCert by ID.
func (r *SSLRepository) GetByID(id uint) (*model.SSLCert, error) {
	var cert model.SSLCert
	if err := r.db.First(&cert, id).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

// Update updates an SSLCert record.
func (r *SSLRepository) Update(cert *model.SSLCert) error {
	return r.db.Save(cert).Error
}

// Delete deletes an SSLCert record.
func (r *SSLRepository) Delete(id uint) error {
	return r.db.Delete(&model.SSLCert{}, id).Error
}

// List returns all SSLCert records.
func (r *SSLRepository) List() ([]model.SSLCert, error) {
	var certs []model.SSLCert
	if err := r.db.Find(&certs).Error; err != nil {
		return nil, err
	}
	return certs, nil
}

// GetByDomain returns an SSLCert by domain.
func (r *SSLRepository) GetByDomain(domain string) (*model.SSLCert, error) {
	var cert model.SSLCert
	if err := r.db.Where("domain = ?", domain).First(&cert).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}
