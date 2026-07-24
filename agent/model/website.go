package model

import "time"

type Website struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	PrimaryDomain string    `json:"primary_domain" gorm:"size:255"`
	Type          string    `json:"type" gorm:"size:20"` // static/proxy/php
	RootDir       string    `json:"root_dir" gorm:"size:255"`
	ProxyTarget   string    `json:"proxy_target" gorm:"size:255"`
	SSLStatus     string    `json:"ssl_status" gorm:"size:20"`
	CertID        uint      `json:"cert_id"`
	Alias         string    `json:"alias" gorm:"size:100"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
