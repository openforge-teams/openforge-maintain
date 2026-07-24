package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openforge-maintain/openforge-maintain/core/middleware"
	"github.com/openforge-maintain/openforge-maintain/core/service"
	"github.com/openforge-maintain/openforge-maintain/pkg/utils"
)

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UserHandler 用户处理器
type UserHandler struct {
	userService    *service.UserService
	authService    *service.AuthService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService *service.UserService, authService *service.AuthService) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
	}
}

// GetProfile 获取当前用户信息
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetAuthUserID(c)
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		utils.Error(c, http.StatusNotFound, 404, "user not found")
		return
	}

	utils.Success(c, user)
}

// UpdateProfile 更新当前用户信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid request parameters: "+err.Error())
		return
	}

	userID := middleware.GetAuthUserID(c)
	if err := h.userService.UpdateProfile(userID, req); err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "failed to update profile: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid request parameters: "+err.Error())
		return
	}

	userID := middleware.GetAuthUserID(c)
	if err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	utils.Success(c, nil)
}

// ListUsers 获取用户列表
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	users, total, err := h.userService.ListUsers(page, size)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "failed to list users: "+err.Error())
		return
	}

	utils.PageSuccess(c, users, total, page, size)
}

// GetUser 获取指定用户信息
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid user id")
		return
	}

	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, 404, "user not found")
		return
	}

	utils.Success(c, user)
}

// UpdateUser 更新指定用户信息
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid user id")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid request parameters: "+err.Error())
		return
	}

	if err := h.userService.UpdateUser(uint(id), req); err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "failed to update user: "+err.Error())
		return
	}

	utils.Success(c, nil)
}

// DeleteUser 删除指定用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, 400, "invalid user id")
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		utils.Error(c, http.StatusInternalServerError, 500, "failed to delete user: "+err.Error())
		return
	}

	utils.Success(c, nil)
}
