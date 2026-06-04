package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mechaniel-coder/netwatch/server/internal/api/handlers"
	"github.com/mechaniel-coder/netwatch/server/internal/config"
	"github.com/mechaniel-coder/netwatch/server/internal/metrics"
	"github.com/mechaniel-coder/netwatch/server/internal/ws"
)

// NewRouter builds and returns the Gin router with all routes registered.
func NewRouter(cfg *config.Config, pool *pgxpool.Pool, hub *ws.Hub) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	agentHandler := handlers.NewAgentHandler(pool, cfg)
	metricHandler := handlers.NewMetricHandler(pool, hub)
	alertHandler := handlers.NewAlertHandler(pool)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Prometheus-format metrics endpoint — scraped by prometheus service
	r.GET("/metrics", func(c *gin.Context) {
		metrics.Handler(c.Writer, c.Request)
	})

	// WebSocket — live metric stream (dashboard subscribes here)
	r.GET("/ws/metrics", hub.ServeWS)

	v1 := r.Group("/api/v1")
	{
		agents := v1.Group("/agents")
		{
			agents.POST("/register", agentHandler.Register)
			agents.POST("/:id/heartbeat", agentHandler.Heartbeat)
			agents.POST("/:id/metrics", metricHandler.Ingest)
			agents.GET("", agentHandler.List)
			agents.GET("/:id", agentHandler.Get)
			agents.GET("/:id/metrics", metricHandler.Query)
		}

		alerts := v1.Group("/alerts")
		{
			alerts.GET("", alertHandler.List)
			alerts.POST("", alertHandler.Create)
			alerts.PUT("/:id", alertHandler.Update)
			alerts.DELETE("/:id", alertHandler.Delete)
			alerts.GET("/events", alertHandler.Events)
		}
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
