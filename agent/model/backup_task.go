package model

import "time"

type BackupTask struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"size:100"`
	TargetType string    `json:"target_type" gorm:"size:20"` // app/website/db/system
	TargetID   uint      `json:"target_id"`
	DestType   string    `json:"dest_type" gorm:"size:20"`   // local/s3/oss/sftp
	DestConfig string    `json:"dest_config" gorm:"type:text"` // JSON
	Schedule   string    `json:"schedule" gorm:"size:100"`
	Retention  int       `json:"retention"`
	Status     string    `json:"status" gorm:"size:20"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
