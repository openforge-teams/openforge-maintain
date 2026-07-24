package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/openforge-maintain/openforge-maintain/core/model"
	"github.com/openforge-maintain/openforge-maintain/core/repository"
	"github.com/openforge-maintain/openforge-maintain/pkg/utils"
	"github.com/pquerna/otp/totp"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidTOTPCode    = errors.New("invalid TOTP code")
	ErrUserAlreadyExists  = errors.New("username already exists")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrPasswordIncorrect  = errors.New("old password is incorrect")
)

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	User         *model.User  `json:"user"`
}

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// AuthService 认证服务
type AuthService struct {
	userRepo repository.UserRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Login 用户登录
func (s *AuthService) Login(username, password, totpCode, ip string) (*LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 校验密码
	ok, err := utils.CheckPassword(password, user.Password)
	if err != nil || !ok {
		return nil, ErrInvalidCredentials
	}

	// 如果用户启用了 MFA，验证 TOTP 码
	if user.MFASecret != "" {
		if totpCode == "" {
			return nil, errors.New("TOTP code required")
		}
		valid := totp.Validate(totpCode, user.MFASecret)
		if !valid {
			return nil, ErrInvalidTOTPCode
		}
	}

	// 生成令牌
	userID := strconv.FormatUint(uint64(user.ID), 10)
	accessToken, refreshToken, err := utils.GenerateToken(userID, user.Username)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(utils.AccessTokenDuration.Seconds()),
		User:         user,
	}, nil
}

// Register 用户注册
func (s *AuthService) Register(username, password, email string) error {
	// 检查用户名是否已存在
	_, err := s.userRepo.GetByUsername(username)
	if err == nil {
		return ErrUserAlreadyExists
	}

	// 对密码进行哈希
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &model.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
		RoleID:   2, // 默认普通用户角色
	}

	return s.userRepo.Create(user)
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := utils.ParseToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// 验证用户是否存在
	userID, err := strconv.ParseUint(claims.UserID, 10, 32)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	_, err = s.userRepo.GetByID(uint(userID))
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// 生成新的令牌对
	accessToken, newRefreshToken, err := utils.GenerateToken(claims.UserID, claims.Username)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(utils.AccessTokenDuration.Seconds()),
	}, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID uint, oldPass, newPass string) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// 校验旧密码
	ok, err := utils.CheckPassword(oldPass, user.Password)
	if err != nil || !ok {
		return ErrPasswordIncorrect
	}

	// 哈希新密码
	hashedPassword, err := utils.HashPassword(newPass)
	if err != nil {
		return err
	}

	return s.userRepo.ChangePassword(userID, hashedPassword)
}

// GetTokenExpiry 获取令牌过期时间
func GetTokenExpiry() time.Duration {
	return utils.AccessTokenDuration
}
