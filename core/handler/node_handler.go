package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/core/service"
	"github.com/openforge-maintain/openforge-maintain/pkg/utils"
)

// NodeHandler 节点处理器
type NodeHandler struct {
	nodeService *service.NodeService
}

// NewNodeHandler 创建节点处理器实例
func NewNodeHandler(nodeService *service.NodeService) *NodeHandler {
	return &NodeHandler{
		nodeService: nodeService,
	}
}

// ListNodes 获取节点列表
func (h *NodeHandler) ListNodes(c *gin.Context) {
	nodes, err := h.nodeService.ListNodes()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "failed to list nodes: "+err.Error())
		return
	}

	utils.Success(c, nodes)
}

// CreateNode 创建节点
func (h *NodeHandler) CreateNode(c *gin.Context) {
	var req service.CreateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid request parameters: "+err.Error())
		return
	}

	node, err := h.nodeService.CreateNode(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "failed to create node: "+err.Error())
		return
	}

	utils.Success(c, node)
}

// GetNode 获取节点详情
func (h *NodeHandler) GetNode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid node id")
		return
	}

	node, err := h.nodeService.GetNode(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, 404, "node not found")
		return
	}

	utils.Success(c, node)
}

// UpdateNode 更新节点
func (h *NodeHandler) UpdateNode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid node id")
		return
	}

	var req service.UpdateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid request parameters: "+err.Error())
		return
	}

	if err := h.nodeService.UpdateNode(uint(id), req); err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "failed to update node: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

// DeleteNode 删除节点
func (h *NodeHandler) DeleteNode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid node id")
		return
	}

	if err := h.nodeService.DeleteNode(uint(id)); err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "failed to delete node: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

// CheckNodeStatus 检查节点状态
func (h *NodeHandler) CheckNodeStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid node id")
		return
	}

	node, err := h.nodeService.CheckNodeStatus(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, 404, "node not found")
		return
	}

	utils.Success(c, node)
}
