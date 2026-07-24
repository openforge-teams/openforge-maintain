package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"github.com/openforge-maintain/openforge-maintain/agent/service/website"
	"github.com/openforge-maintain/openforge-maintain/pkg/response"
)

// WebsiteHandler handles website management API requests.
type WebsiteHandler struct {
	websiteService *website.WebsiteService
}

// NewWebsiteHandler creates a new WebsiteHandler.
func NewWebsiteHandler(websiteService *website.WebsiteService) *WebsiteHandler {
	return &WebsiteHandler{websiteService: websiteService}
}

// CreateWebsite creates a new website.
func (h *WebsiteHandler) CreateWebsite(c *gin.Context) {
	var ws model.Website
	if err := c.ShouldBindJSON(&ws); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.websiteService.Create(&ws); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, ws)
}

// UpdateWebsite updates a website.
func (h *WebsiteHandler) UpdateWebsite(c *gin.Context) {
	var ws model.Website
	if err := c.ShouldBindJSON(&ws); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.websiteService.Update(&ws); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, ws)
}

// DeleteWebsite deletes a website.
func (h *WebsiteHandler) DeleteWebsite(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.websiteService.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetWebsite returns a website by ID.
func (h *WebsiteHandler) GetWebsite(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	ws, err := h.websiteService.Get(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, ws)
}

// ListWebsites returns a paginated list of websites.
func (h *WebsiteHandler) ListWebsites(c *gin.Context) {
	page := parseInt(c.DefaultQuery("page", "1"), 1)
	size := parseInt(c.DefaultQuery("size", "10"), 10)

	websites, total, err := h.websiteService.List(page, size)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"list":  websites,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// EnableSSL enables SSL for a website.
func (h *WebsiteHandler) EnableSSL(c *gin.Context) {
	var req struct {
		WebsiteID uint `json:"website_id" binding:"required"`
		CertID    uint `json:"cert_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.websiteService.EnableSSL(req.WebsiteID, req.CertID); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// DisableSSL disables SSL for a website.
func (h *WebsiteHandler) DisableSSL(c *gin.Context) {
	var req struct {
		WebsiteID uint `json:"website_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.websiteService.DisableSSL(req.WebsiteID); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetNginxConfig generates and returns the Nginx config for a website.
func (h *WebsiteHandler) GetNginxConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	ws, err := h.websiteService.Get(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	config, err := h.websiteService.GenerateNginxConfig(ws)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"config": config})
}
