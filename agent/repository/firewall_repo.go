package repository

import (
	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"gorm.io/gorm"
)

// FirewallRepository provides CRUD operations for FirewallRule.
type FirewallRepository struct {
	db *gorm.DB
}

// NewFirewallRepository creates a new FirewallRepository.
func NewFirewallRepository(db *gorm.DB) *FirewallRepository {
	return &FirewallRepository{db: db}
}

// Create creates a new FirewallRule record.
func (r *FirewallRepository) Create(rule *model.FirewallRule) error {
	return r.db.Create(rule).Error
}

// GetByID returns a FirewallRule by ID.
func (r *FirewallRepository) GetByID(id uint) (*model.FirewallRule, error) {
	var rule model.FirewallRule
	if err := r.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// Update updates a FirewallRule record.
func (r *FirewallRepository) Update(rule *model.FirewallRule) error {
	return r.db.Save(rule).Error
}

// Delete deletes a FirewallRule record.
func (r *FirewallRepository) Delete(id uint) error {
	return r.db.Delete(&model.FirewallRule{}, id).Error
}

// List returns all FirewallRule records.
func (r *FirewallRepository) List() ([]model.FirewallRule, error) {
	var rules []model.FirewallRule
	if err := r.db.Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}
