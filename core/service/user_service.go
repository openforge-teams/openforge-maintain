package service

import (
	"github.com/openforge-maintain/openforge-maintain/core/model"
	"github.com/openforge-maintain/openforge-maintain/core/repository"
)

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	Email string `json:"email" binding:"omitempty,email"`
}

// UpdateUserRequest 更新用户请求（管理员）
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	RoleID   uint   `json:"role_id"`
}

// UserService 用户服务
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetProfile 获取用户信息
func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.GetByID(userID)
}

// UpdateProfile 更新用户信息
func (s *UserService) UpdateProfile(userID uint, req UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	return s.userRepo.Update(user)
}

// ListUsers 获取用户列表（分页）
func (s *UserService) ListUsers(page, size int) ([]model.User, int64, error) {
	return s.userRepo.List(page, size)
}

// GetUser 获取指定用户
func (s *UserService) GetUser(userID uint) (*model.User, error) {
	return s.userRepo.GetByID(userID)
}

// UpdateUser 更新用户（管理员）
func (s *UserService) UpdateUser(userID uint, req UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.RoleID > 0 {
		user.RoleID = req.RoleID
	}

	return s.userRepo.Update(user)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(userID uint) error {
	return s.userRepo.Delete(userID)
}
