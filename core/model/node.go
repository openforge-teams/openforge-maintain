package model

import "time"

// Node 节点模型
type Node struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	AgentAddr string    `json:"agent_addr" gorm:"size:255"`
	Token     string    `json:"token" gorm:"size:255"`
	Status    string    `json:"status" gorm:"size:20;default:offline"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
