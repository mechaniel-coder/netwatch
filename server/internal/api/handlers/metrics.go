package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mechaniel-coder/netwatch/server/internal/metrics"
	"github.com/mechaniel-coder/netwatch/server/internal/models"
	"github.com/mechaniel-coder/netwatch/server/internal/ws"
)

type MetricHandler struct {
	pool *pgxpool.Pool
	hub  *ws.Hub
}

func NewMetricHandler(pool *pgxpool.Pool, hub *ws.Hub) *MetricHandler {
	return &MetricHandler{pool: pool, hub: hub}
}

// Ingest handles POST /api/v1/agents/:id/metrics
// Agents call this every reporting interval (default: 15s).
func (h *MetricHandler) Ingest(c *gin.Context) {
	agentID := c.Param("id")

	var snap models.MetricSnapshot
	if err := c.ShouldBindJSON(&snap); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const q = `
		INSERT INTO metrics
		  (agent_id, cpu_percent, mem_used_bytes, mem_total_bytes,
		   disk_used_bytes, disk_total_bytes, net_rx_bytes, net_tx_bytes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`
	if _, err := h.pool.Exec(c.Request.Context(), q,
		agentID,
		snap.CPUPercent, snap.MemUsedBytes, snap.MemTotalBytes,
		snap.DiskUsedBytes, snap.DiskTotalBytes,
		snap.NetRxBytes, snap.NetTxBytes,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store metrics"})
		return
	}

	// Broadcast the new snapshot to all connected dashboard clients
	payload, _ := json.Marshal(gin.H{"agent_id": agentID, "metrics": snap})
	h.hub.Broadcast(payload)

	metrics.IncMetricsIngested()
	c.JSON(http.StatusOK, gin.H{"status": "accepted"})
}

// Query handles GET /api/v1/agents/:id/metrics
func (h *MetricHandler) Query(c *gin.Context) {
	var q models.MetricQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	q.AgentID = c.Param("id")
	if q.Limit <= 0 || q.Limit > 5000 {
		q.Limit = 500
	}

	const sql = `
		SELECT id, agent_id, recorded_at,
		       cpu_percent, mem_used_bytes, mem_total_bytes,
		       disk_used_bytes, disk_total_bytes, net_rx_bytes, net_tx_bytes
		FROM metrics
		WHERE agent_id = $1
		  AND ($2::timestamptz IS NULL OR recorded_at >= $2)
		  AND ($3::timestamptz IS NULL OR recorded_at <= $3)
		ORDER BY recorded_at DESC
		LIMIT $4
	`
	var from, to interface{}
	if !q.From.IsZero() {
		from = q.From
	}
	if !q.To.IsZero() {
		to = q.To
	}

	rows, err := h.pool.Query(c.Request.Context(), sql, q.AgentID, from, to, q.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	defer rows.Close()

	var metrics []models.Metric
	for rows.Next() {
		var m models.Metric
		if err := rows.Scan(&m.ID, &m.AgentID, &m.RecordedAt,
			&m.CPUPercent, &m.MemUsedBytes, &m.MemTotalBytes,
			&m.DiskUsedBytes, &m.DiskTotalBytes, &m.NetRxBytes, &m.NetTxBytes,
		); err != nil {
			continue
		}
		metrics = append(metrics, m)
	}
	c.JSON(http.StatusOK, metrics)
}
