package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertHandler struct {
	pool *pgxpool.Pool
}

func NewAlertHandler(pool *pgxpool.Pool) *AlertHandler {
	return &AlertHandler{pool: pool}
}

type alertRuleInput struct {
	AgentID   *string `json:"agent_id"`
	Metric    string  `json:"metric" binding:"required"`
	Operator  string  `json:"operator" binding:"required,oneof=gt lt gte lte"`
	Threshold float64 `json:"threshold" binding:"required"`
	Severity  string  `json:"severity" binding:"required,oneof=warning critical"`
}

// List handles GET /api/v1/alerts
func (h *AlertHandler) List(c *gin.Context) {
	const q = `SELECT id, agent_id, metric, operator, threshold, severity, enabled, created_at FROM alert_rules ORDER BY created_at DESC`
	rows, err := h.pool.Query(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	defer rows.Close()

	var rules []gin.H
	for rows.Next() {
		var r struct {
			ID        string  `json:"id"`
			AgentID   *string `json:"agent_id"`
			Metric    string  `json:"metric"`
			Operator  string  `json:"operator"`
			Threshold float64 `json:"threshold"`
			Severity  string  `json:"severity"`
			Enabled   bool    `json:"enabled"`
			CreatedAt string  `json:"created_at"`
		}
		rows.Scan(&r.ID, &r.AgentID, &r.Metric, &r.Operator, &r.Threshold, &r.Severity, &r.Enabled, &r.CreatedAt)
		rules = append(rules, gin.H{
			"id": r.ID, "agent_id": r.AgentID, "metric": r.Metric,
			"operator": r.Operator, "threshold": r.Threshold,
			"severity": r.Severity, "enabled": r.Enabled, "created_at": r.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, rules)
}

// Create handles POST /api/v1/alerts
func (h *AlertHandler) Create(c *gin.Context) {
	var input alertRuleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const q = `
		INSERT INTO alert_rules (agent_id, metric, operator, threshold, severity)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id string
	if err := h.pool.QueryRow(c.Request.Context(), q,
		input.AgentID, input.Metric, input.Operator, input.Threshold, input.Severity,
	).Scan(&id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create rule"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// Update handles PUT /api/v1/alerts/:id
func (h *AlertHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Enabled   *bool    `json:"enabled"`
		Threshold *float64 `json:"threshold"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const q = `
		UPDATE alert_rules
		SET enabled   = COALESCE($2, enabled),
		    threshold = COALESCE($3, threshold)
		WHERE id = $1
	`
	h.pool.Exec(c.Request.Context(), q, id, input.Enabled, input.Threshold)
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// Delete handles DELETE /api/v1/alerts/:id
func (h *AlertHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	h.pool.Exec(c.Request.Context(), `DELETE FROM alert_rules WHERE id = $1`, id)
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// Events handles GET /api/v1/alerts/events
func (h *AlertHandler) Events(c *gin.Context) {
	const q = `
		SELECT ae.id, ae.rule_id, ae.agent_id, ae.value, ae.severity, ae.fired_at, ae.resolved_at
		FROM alert_events ae
		ORDER BY ae.fired_at DESC
		LIMIT 200
	`
	rows, err := h.pool.Query(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	defer rows.Close()

	var events []gin.H
	for rows.Next() {
		var (
			id, ruleID, agentID, severity, firedAt string
			value                                  float64
			resolvedAt                             *string
		)
		rows.Scan(&id, &ruleID, &agentID, &value, &severity, &firedAt, &resolvedAt)
		events = append(events, gin.H{
			"id": id, "rule_id": ruleID, "agent_id": agentID,
			"value": value, "severity": severity,
			"fired_at": firedAt, "resolved_at": resolvedAt,
		})
	}
	c.JSON(http.StatusOK, events)
}
