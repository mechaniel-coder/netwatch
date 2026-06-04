package models

import (
	"time"

	"github.com/google/uuid"
)

type Metric struct {
	ID             int64     `json:"id" db:"id"`
	AgentID        uuid.UUID `json:"agent_id" db:"agent_id"`
	RecordedAt     time.Time `json:"recorded_at" db:"recorded_at"`
	CPUPercent     float64   `json:"cpu_percent" db:"cpu_percent"`
	MemUsedBytes   int64     `json:"mem_used_bytes" db:"mem_used_bytes"`
	MemTotalBytes  int64     `json:"mem_total_bytes" db:"mem_total_bytes"`
	DiskUsedBytes  int64     `json:"disk_used_bytes" db:"disk_used_bytes"`
	DiskTotalBytes int64     `json:"disk_total_bytes" db:"disk_total_bytes"`
	NetRxBytes     int64     `json:"net_rx_bytes" db:"net_rx_bytes"`
	NetTxBytes     int64     `json:"net_tx_bytes" db:"net_tx_bytes"`
}

type MetricSnapshot struct {
	CPUPercent     float64 `json:"cpu_percent" binding:"required,min=0,max=100"`
	MemUsedBytes   int64   `json:"mem_used_bytes" binding:"required,min=0"`
	MemTotalBytes  int64   `json:"mem_total_bytes" binding:"required,min=1"`
	DiskUsedBytes  int64   `json:"disk_used_bytes" binding:"required,min=0"`
	DiskTotalBytes int64   `json:"disk_total_bytes" binding:"required,min=1"`
	NetRxBytes     int64   `json:"net_rx_bytes" binding:"min=0"`
	NetTxBytes     int64   `json:"net_tx_bytes" binding:"min=0"`
}

type MetricQuery struct {
	AgentID string    `form:"agent_id" binding:"required,uuid"`
	From    time.Time `form:"from" time_format:"2006-01-02T15:04:05Z"`
	To      time.Time `form:"to" time_format:"2006-01-02T15:04:05Z"`
	Limit   int       `form:"limit,default=500"`
}
