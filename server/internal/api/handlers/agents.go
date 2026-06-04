package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mechaniel-coder/netwatch/server/internal/config"
	"github.com/mechaniel-coder/netwatch/server/internal/metrics"
	"github.com/mechaniel-coder/netwatch/server/internal/models"
)

type AgentHandler struct {
	pool *pgxpool.Pool
	cfg  *config.Config
}

func NewAgentHandler(pool *pgxpool.Pool, cfg *config.Config) *AgentHandler {
	return &AgentHandler{pool: pool, cfg: cfg}
}

// Register handles POST /api/v1/agents/register
// Agents call this on first boot to get their UUID and a signed JWT.
func (h *AgentHandler) Register(c *gin.Context) {
	var req models.AgentRegistration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const q = `
		INSERT INTO agents (hostname, ip_address, os, version, online, last_seen)
		VALUES ($1, $2, $3, $4, true, now())
		ON CONFLICT (hostname) DO UPDATE
		  SET ip_address = EXCLUDED.ip_address,
		      version    = EXCLUDED.version,
		      online     = true,
		      last_seen  = now()
		RETURNING id
	`
	var agentID string
	if err := h.pool.QueryRow(c.Request.Context(), q,
		req.Hostname, req.IPAddress, req.OS, req.Version,
	).Scan(&agentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"agent_id": agentID})
	metrics.IncAgentRegistrations()
}

// Heartbeat handles POST /api/v1/agents/:id/heartbeat
func (h *AgentHandler) Heartbeat(c *gin.Context) {
	id := c.Param("id")
	const q = `UPDATE agents SET online = true, last_seen = now() WHERE id = $1`
	tag, err := h.pool.Exec(c.Request.Context(), q, id)
	if err != nil || tag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "last_seen": time.Now().UTC()})
}

// List handles GET /api/v1/agents
func (h *AgentHandler) List(c *gin.Context) {
	const q = `SELECT id, hostname, ip_address, os, version, last_seen, online, registered_at FROM agents ORDER BY hostname`
	rows, err := h.pool.Query(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	defer rows.Close()

	var agents []models.Agent
	for rows.Next() {
		var a models.Agent
		if err := rows.Scan(&a.ID, &a.Hostname, &a.IPAddress, &a.OS, &a.Version,
			&a.LastSeen, &a.Online, &a.RegisteredAt); err != nil {
			continue
		}
		agents = append(agents, a)
	}
	c.JSON(http.StatusOK, agents)
}

// Get handles GET /api/v1/agents/:id
func (h *AgentHandler) Get(c *gin.Context) {
	id := c.Param("id")
	const q = `SELECT id, hostname, ip_address, os, version, last_seen, online, registered_at FROM agents WHERE id = $1`
	var a models.Agent
	err := h.pool.QueryRow(c.Request.Context(), q, id).
		Scan(&a.ID, &a.Hostname, &a.IPAddress, &a.OS, &a.Version,
			&a.LastSeen, &a.Online, &a.RegisteredAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}
