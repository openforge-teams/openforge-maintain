package repository

import (
	"github.com/openforge-maintain/openforge-maintain/core/model"
	"gorm.io/gorm"
)

// AuditFilter 审计日志过滤器
type AuditFilter struct {
	UserID   uint
	Action   string
	Resource string
	IP       string
}

// AuditRepository 审计日志仓库接口
type AuditRepository interface {
	Create(log *model.AuditLog) error
	List(filter AuditFilter, page, size int) ([]model.AuditLog, int64, error)
}

// auditRepository 审计日志仓库实现
type auditRepository struct {
	DB *gorm.DB
}

// NewAuditRepository 创建审计日志仓库实例
func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{DB: db}
}

// Create 创建审计日志
func (r *auditRepository) Create(log *model.AuditLog) error {
	return r.DB.Create(log).Error
}

// List 分页查询审计日志
func (r *auditRepository) List(filter AuditFilter, page, size int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	query := r.DB.Model(&model.AuditLog{})

	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.Resource != "" {
		query = query.Where("resource = ?", filter.Resource)
	}
	if filter.IP != "" {
		query = query.Where("ip = ?", filter.IP)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := query.Order("created_at DESC").Offset(offset).Limit(size).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
