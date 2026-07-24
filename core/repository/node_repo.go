package repository

import (
	"github.com/openforge-maintain/openforge-maintain/core/model"
	"gorm.io/gorm"
)

// NodeRepository 节点仓库接口
type NodeRepository interface {
	Create(node *model.Node) error
	Update(node *model.Node) error
	Delete(id uint) error
	GetByID(id uint) (*model.Node, error)
	List() ([]model.Node, error)
}

// nodeRepository 节点仓库实现
type nodeRepository struct {
	DB *gorm.DB
}

// NewNodeRepository 创建节点仓库实例
func NewNodeRepository(db *gorm.DB) NodeRepository {
	return &nodeRepository{DB: db}
}

// Create 创建节点
func (r *nodeRepository) Create(node *model.Node) error {
	return r.DB.Create(node).Error
}

// Update 更新节点
func (r *nodeRepository) Update(node *model.Node) error {
	return r.DB.Save(node).Error
}

// Delete 删除节点
func (r *nodeRepository) Delete(id uint) error {
	return r.DB.Delete(&model.Node{}, id).Error
}

// GetByID 根据ID获取节点
func (r *nodeRepository) GetByID(id uint) (*model.Node, error) {
	var node model.Node
	err := r.DB.First(&node, id).Error
	if err != nil {
		return nil, err
	}
	return &node, nil
}

// List 获取所有节点
func (r *nodeRepository) List() ([]model.Node, error) {
	var nodes []model.Node
	err := r.DB.Find(&nodes).Error
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
