package model

import "time"

type SSLCert struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Domain      string    `json:"domain" gorm:"size:255"`
	Provider    string    `json:"provider" gorm:"size:50"`
	CAType      string    `json:"ca_type" gorm:"size:50"`
	CertPath    string    `json:"cert_path" gorm:"size:255"`
	KeyPath     string    `json:"key_path" gorm:"size:255"`
	DNSProvider string    `json:"dns_provider" gorm:"size:50"`
	AutoRenew   bool      `json:"auto_renew"`
	ExpiredAt   time.Time `json:"expired_at"`
	CreatedAt   time.Time `json:"created_at"`
}
