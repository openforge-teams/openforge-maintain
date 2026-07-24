package model

import "time"

type CronJob struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"size:100"`
	Spec      string     `json:"spec" gorm:"size:100"` // cron expression
	Command   string     `json:"command" gorm:"type:text"`
	Type      string     `json:"type" gorm:"size:20"`  // shell/http/script/backup
	Status    string     `json:"status" gorm:"size:20"`
	LastRun   *time.Time `json:"last_run"`
	NextRun   *time.Time `json:"next_run"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
