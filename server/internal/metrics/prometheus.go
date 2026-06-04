// Package metrics exposes a Prometheus-compatible text-format /metrics endpoint
// using only the Go standard library — no external dependencies required.
//
// Counters and gauges are updated via the exported Increment/Set helpers.
// All reads and writes use sync/atomic for lock-free correctness.
package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	startTime = time.Now()

	agentRegistrations int64
	metricsIngested    int64
	wsConnections      int64
	alertsFired        int64
)

// IncAgentRegistrations increments the agent registration counter.
func IncAgentRegistrations() { atomic.AddInt64(&agentRegistrations, 1) }

// IncMetricsIngested increments the metric-ingest counter.
func IncMetricsIngested() { atomic.AddInt64(&metricsIngested, 1) }

// IncWSConnections increments the active WebSocket connection gauge.
func IncWSConnections() { atomic.AddInt64(&wsConnections, 1) }

// DecWSConnections decrements the active WebSocket connection gauge.
func DecWSConnections() { atomic.AddInt64(&wsConnections, -1) }

// IncAlertsFired increments the alerts-fired counter.
func IncAlertsFired() { atomic.AddInt64(&alertsFired, 1) }

// Handler writes Prometheus text-exposition format to w.
// Register it at GET /metrics on your HTTP mux.
func Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	uptimeSeconds := time.Since(startTime).Seconds()

	fmt.Fprintf(w, "# HELP netwatch_uptime_seconds Seconds since the server started.\n")
	fmt.Fprintf(w, "# TYPE netwatch_uptime_seconds gauge\n")
	fmt.Fprintf(w, "netwatch_uptime_seconds %.2f\n\n", uptimeSeconds)

	fmt.Fprintf(w, "# HELP netwatch_agent_registrations_total Total agent registrations since startup.\n")
	fmt.Fprintf(w, "# TYPE netwatch_agent_registrations_total counter\n")
	fmt.Fprintf(w, "netwatch_agent_registrations_total %d\n\n", atomic.LoadInt64(&agentRegistrations))

	fmt.Fprintf(w, "# HELP netwatch_metrics_ingested_total Total metric snapshots ingested since startup.\n")
	fmt.Fprintf(w, "# TYPE netwatch_metrics_ingested_total counter\n")
	fmt.Fprintf(w, "netwatch_metrics_ingested_total %d\n\n", atomic.LoadInt64(&metricsIngested))

	fmt.Fprintf(w, "# HELP netwatch_websocket_connections Current number of active WebSocket dashboard connections.\n")
	fmt.Fprintf(w, "# TYPE netwatch_websocket_connections gauge\n")
	fmt.Fprintf(w, "netwatch_websocket_connections %d\n\n", atomic.LoadInt64(&wsConnections))

	fmt.Fprintf(w, "# HELP netwatch_alerts_fired_total Total alert events fired since startup.\n")
	fmt.Fprintf(w, "# TYPE netwatch_alerts_fired_total counter\n")
	fmt.Fprintf(w, "netwatch_alerts_fired_total %d\n", atomic.LoadInt64(&alertsFired))
}
