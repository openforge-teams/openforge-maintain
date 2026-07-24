package model

import "time"

type AppInstall struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	AppKey      string    `json:"app_key" gorm:"size:100;index"`
	Name        string    `json:"name" gorm:"size:100"`
	Version     string    `json:"version" gorm:"size:50"`
	Status      string    `json:"status" gorm:"size:20"`
	ComposeFile string    `json:"compose_file" gorm:"type:text"`
	EnvConfig   string    `json:"env_config" gorm:"type:text"`
	Port        int       `json:"port"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
