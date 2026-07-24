package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/core/service"
	"github.com/openforge-maintain/openforge-maintain/pkg/utils"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	TOTPCode string `json:"totp_code"`
}

// RefreshRequest 刷新令牌请求
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名和密码登录系统
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body LoginRequest true "登录信息"
// @Success 200 {object} utils.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid request parameters: "+err.Error())
		return
	}

	ip := c.ClientIP()
	resp, err := h.authService.Login(req.Username, req.Password, req.TOTPCode, ip)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	utils.Success(c, resp)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 退出登录（客户端删除令牌）
// @Tags 认证
// @Produce json
// @Success 200 {object} utils.Response
// @Security BearerAuth
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT 是无状态的，登出主要由客户端处理
	// 可以在此处将 token 加入黑名单（如果实现了黑名单机制）
	utils.Success(c, nil)
}

// Refresh 刷新令牌
// @Summary 刷新令牌
// @Description 使用 refresh token 获取新的 access token
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body RefreshRequest true "刷新令牌"
// @Success 200 {object} utils.Response
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid request parameters: "+err.Error())
		return
	}

	tokenPair, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	utils.Success(c, tokenPair)
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body service.RegisterRequest true "注册信息"
// @Success 200 {object} utils.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid request parameters: "+err.Error())
		return
	}

	if err := h.authService.Register(req.Username, req.Password, req.Email); err != nil {
		utils.Error(c, http.StatusConflict, 409, err.Error())
		return
	}

	utils.Success(c, nil)
}
