package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/agent/service/appstore"
	"github.com/openforge-maintain/openforge-maintain/pkg/response"
)

// AppStoreHandler handles app store API requests.
type AppStoreHandler struct {
	appStoreService *appstore.AppStoreService
}

// NewAppStoreHandler creates a new AppStoreHandler.
func NewAppStoreHandler(appStoreService *appstore.AppStoreService) *AppStoreHandler {
	return &AppStoreHandler{appStoreService: appStoreService}
}

// GetAppList returns the list of available applications.
func (h *AppStoreHandler) GetAppList(c *gin.Context) {
	apps, err := h.appStoreService.GetAppList()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, apps)
}

// GetAppDetail returns detailed information about a specific application.
func (h *AppStoreHandler) GetAppDetail(c *gin.Context) {
	appKey := c.Param("appKey")
	app, err := h.appStoreService.GetAppDetail(appKey)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, app)
}

// InstallApp installs an application.
func (h *AppStoreHandler) InstallApp(c *gin.Context) {
	var req struct {
		AppKey string            `json:"app_key" binding:"required"`
		Params map[string]string `json:"params"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.appStoreService.Install(req.AppKey, req.Params); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// UninstallApp uninstalls an application.
func (h *AppStoreHandler) UninstallApp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.appStoreService.Uninstall(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// UpgradeApp upgrades an installed application.
func (h *AppStoreHandler) UpgradeApp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.appStoreService.Upgrade(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetInstalledApps returns all installed applications.
func (h *AppStoreHandler) GetInstalledApps(c *gin.Context) {
	apps, err := h.appStoreService.GetInstalledApps()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, apps)
}
