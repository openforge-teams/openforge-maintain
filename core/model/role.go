package model

// Role 角色模型
type Role struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"size:50;not null"`
	Description string `json:"description"`
}
