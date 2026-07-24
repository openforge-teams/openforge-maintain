package repository

import (
	"github.com/openforge-maintain/openforge-maintain/core/model"
	"gorm.io/gorm"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	List(page, size int) ([]model.User, int64, error)
	ChangePassword(userID uint, newPassword string) error
}

// userRepository 用户仓库实现
type userRepository struct {
	DB *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

// Update 更新用户
func (r *userRepository) Update(user *model.User) error {
	return r.DB.Save(user).Error
}

// Delete 删除用户
func (r *userRepository) Delete(id uint) error {
	return r.DB.Delete(&model.User{}, id).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.DB.Preload("Role").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.Preload("Role").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List 分页获取用户列表
func (r *userRepository) List(page, size int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if err := r.DB.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.DB.Preload("Role").Offset(offset).Limit(size).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ChangePassword 修改用户密码
func (r *userRepository) ChangePassword(userID uint, newPassword string) error {
	return r.DB.Model(&model.User{}).Where("id = ?", userID).Update("password", newPassword).Error
}
