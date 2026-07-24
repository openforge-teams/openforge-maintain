package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/agent/service/ssl"
	"github.com/openforge-maintain/openforge-maintain/pkg/response"
)

// SSLHandler handles SSL certificate API requests.
type SSLHandler struct {
	sslService *ssl.SSLService
}

// NewSSLHandler creates a new SSLHandler.
func NewSSLHandler(sslService *ssl.SSLService) *SSLHandler {
	return &SSLHandler{sslService: sslService}
}

// RequestCert requests a new SSL certificate.
func (h *SSLHandler) RequestCert(c *gin.Context) {
	var req struct {
		Domain      string            `json:"domain" binding:"required"`
		Provider    string            `json:"provider"`
		DNSProvider string            `json:"dns_provider"`
		DNSConfig   map[string]string `json:"dns_config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.sslService.RequestCert(req.Domain, req.Provider, req.DNSProvider, req.DNSConfig); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// RenewCert renews an SSL certificate.
func (h *SSLHandler) RenewCert(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.sslService.RenewCert(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// ListCerts lists all SSL certificates.
func (h *SSLHandler) ListCerts(c *gin.Context) {
	certs, err := h.sslService.ListCerts()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, certs)
}

// GetCertDetail returns detailed information about an SSL certificate.
func (h *SSLHandler) GetCertDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	cert, err := h.sslService.GetCertDetail(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cert)
}

// DeleteCert deletes an SSL certificate.
func (h *SSLHandler) DeleteCert(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.sslService.DeleteCert(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}
