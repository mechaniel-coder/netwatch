package models

import (
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Hostname     string     `json:"hostname" db:"hostname"`
	IPAddress    string     `json:"ip_address" db:"ip_address"`
	OS           string     `json:"os" db:"os"`
	Version      string     `json:"version" db:"version"`
	LastSeen     *time.Time `json:"last_seen" db:"last_seen"`
	Online       bool       `json:"online" db:"online"`
	RegisteredAt time.Time  `json:"registered_at" db:"registered_at"`
}

type AgentRegistration struct {
	Hostname  string `json:"hostname" binding:"required"`
	IPAddress string `json:"ip_address" binding:"required"`
	OS        string `json:"os" binding:"required"`
	Version   string `json:"version" binding:"required"`
}
