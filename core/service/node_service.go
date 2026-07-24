package service

import (
	"github.com/openforge-maintain/openforge-maintain/core/model"
	"github.com/openforge-maintain/openforge-maintain/core/repository"
)

// CreateNodeRequest 创建节点请求
type CreateNodeRequest struct {
	Name      string `json:"name" binding:"required,min=1,max=100"`
	AgentAddr string `json:"agent_addr" binding:"required"`
}

// UpdateNodeRequest 更新节点请求
type UpdateNodeRequest struct {
	Name      string `json:"name" binding:"omitempty,min=1,max=100"`
	AgentAddr string `json:"agent_addr" binding:"omitempty"`
	Status    string `json:"status" binding:"omitempty,oneof=online offline"`
}

// NodeService 节点服务
type NodeService struct {
	nodeRepo repository.NodeRepository
}

// NewNodeService 创建节点服务实例
func NewNodeService(nodeRepo repository.NodeRepository) *NodeService {
	return &NodeService{
		nodeRepo: nodeRepo,
	}
}

// ListNodes 获取所有节点
func (s *NodeService) ListNodes() ([]model.Node, error) {
	return s.nodeRepo.List()
}

// CreateNode 创建节点
func (s *NodeService) CreateNode(req CreateNodeRequest) (*model.Node, error) {
	node := &model.Node{
		Name:      req.Name,
		AgentAddr: req.AgentAddr,
		Status:    "offline",
	}

	if err := s.nodeRepo.Create(node); err != nil {
		return nil, err
	}

	return node, nil
}

// GetNode 获取节点详情
func (s *NodeService) GetNode(id uint) (*model.Node, error) {
	return s.nodeRepo.GetByID(id)
}

// UpdateNode 更新节点
func (s *NodeService) UpdateNode(id uint, req UpdateNodeRequest) error {
	node, err := s.nodeRepo.GetByID(id)
	if err != nil {
		return err
	}

	if req.Name != "" {
		node.Name = req.Name
	}
	if req.AgentAddr != "" {
		node.AgentAddr = req.AgentAddr
	}
	if req.Status != "" {
		node.Status = req.Status
	}

	return s.nodeRepo.Update(node)
}

// DeleteNode 删除节点
func (s *NodeService) DeleteNode(id uint) error {
	return s.nodeRepo.Delete(id)
}

// CheckNodeStatus 检查节点状态
func (s *NodeService) CheckNodeStatus(id uint) (*model.Node, error) {
	node, err := s.nodeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return node, nil
}
