package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"github.com/openforge-maintain/openforge-maintain/agent/service/firewall"
	"github.com/openforge-maintain/openforge-maintain/pkg/response"
)

// FirewallHandler handles firewall API requests.
type FirewallHandler struct {
	firewallService *firewall.FirewallService
}

// NewFirewallHandler creates a new FirewallHandler.
func NewFirewallHandler(firewallService *firewall.FirewallService) *FirewallHandler {
	return &FirewallHandler{firewallService: firewallService}
}

// ListRules lists all firewall rules.
func (h *FirewallHandler) ListRules(c *gin.Context) {
	rules, err := h.firewallService.ListRules()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, rules)
}

// AddRule adds a new firewall rule.
func (h *FirewallHandler) AddRule(c *gin.Context) {
	var rule model.FirewallRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.firewallService.AddRule(&rule); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, rule)
}

// DeleteRule deletes a firewall rule.
func (h *FirewallHandler) DeleteRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.firewallService.DeleteRule(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// EnableRule enables a firewall rule.
func (h *FirewallHandler) EnableRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.firewallService.EnableRule(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// DisableRule disables a firewall rule.
func (h *FirewallHandler) DisableRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.firewallService.DisableRule(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetFirewallStatus returns the current firewall status.
func (h *FirewallHandler) GetFirewallStatus(c *gin.Context) {
	status, err := h.firewallService.GetStatus()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"status": status})
}

// EnableFirewall enables the firewall.
func (h *FirewallHandler) EnableFirewall(c *gin.Context) {
	if err := h.firewallService.Enable(); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// DisableFirewall disables the firewall.
func (h *FirewallHandler) DisableFirewall(c *gin.Context) {
	if err := h.firewallService.Disable(); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}
