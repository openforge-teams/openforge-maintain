package model

import "time"

type FirewallRule struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Protocol  string    `json:"protocol" gorm:"size:10"` // tcp/udp/icmp
	Port      string    `json:"port" gorm:"size:50"`    // 80, 443, 8000-9000
	Source    string    `json:"source" gorm:"size:50"`  // 0.0.0.0/0 or specific IP
	Action    string    `json:"action" gorm:"size:20"`  // allow/deny
	Comment   string    `json:"comment" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at"`
}
