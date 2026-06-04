# NetWatch — Infrastructure Monitoring Platform

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat&logo=go)
![React](https://img.shields.io/badge/React-19-61DAFB?style=flat&logo=react)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)
![WebSockets](https://img.shields.io/badge/WebSockets-Live_Metrics-brightgreen?style=flat)
![License](https://img.shields.io/badge/License-MIT-yellow?style=flat)

Real-time infrastructure monitoring for servers, VMs, and bare-metal machines. A lightweight Go agent runs on each target host, ships CPU, memory, disk, and network metrics to a central Go API server, and a React dashboard visualizes everything live over WebSockets with configurable threshold alerts.

Built for IT shops, MSPs, and anyone who manages more than one machine.

---

## Features

- **Lightweight Go agent** — single binary, ~6 MB, zero runtime dependencies
- **Real-time dashboard** — WebSocket-driven charts, no polling
- **Threshold alerting** — per-metric alert rules with severity levels (warning / critical)
- **Historical data** — PostgreSQL time-series storage with configurable retention
- **Multi-host support** — manage hundreds of agents from one dashboard
- **Agent auto-registration** — agents register on first contact, no manual provisioning
- **JWT-secured API** — all agent communications and dashboard API calls are authenticated
- **Docker Compose deployment** — entire stack up in one command
- **Graceful degradation** — dashboard shows last known state when agent goes offline

---

## Architecture

```mermaid
flowchart TD
    subgraph HOSTS [Target Hosts]
        A1[Go Agent] 
        A2[Go Agent]
        AN[Go Agent ...]
    end
    HOSTS -->|HTTPS POST /metrics\nJWT auth| API
    subgraph API [Go API Server — Gin]
        REST[REST endpoints]
        WS[WebSocket Hub]
        ALERT[Alert Engine]
        METRICS[/metrics\nPrometheus text format]
    end
    API --> PG[(PostgreSQL 16)]
    API --> REDIS[(Redis 7)]
    WS -->|live push| UI[React Dashboard\nRecharts · Zustand]
    PROM[Prometheus] -->|scrape /metrics| METRICS
    GRAFANA[Grafana] -->|query| PROM
```

---

## Tech Stack

| Layer         | Technology                                  |
|---------------|---------------------------------------------|
| Agent         | Go 1.22 · `gopsutil` · `net/http`           |
| API Server    | Go 1.22 · Gin · `gorilla/websocket`         |
| Auth          | JWT (RS256) · bcrypt                        |
| Database      | PostgreSQL 16 · `pgx/v5` · raw migrations   |
| Cache/State   | Redis 7                                      |
| Frontend      | React 19 · TypeScript · Vite · Zustand      |
| Charts        | Recharts                                    |
| Deployment    | Docker · Docker Compose                     |

---

## Quick Start

**Prerequisites:** Docker and Docker Compose installed.

```bash
git clone https://github.com/mechaniel-coder/netwatch.git
cd netwatch
cp .env.example .env        # review and edit as needed
docker compose up -d
```

Dashboard → [http://localhost:3000](http://localhost:3000)  
API → [http://localhost:8080](http://localhost:8080)

**Deploy an agent on a target host:**

```bash
# Download the agent binary (cross-compiled for Linux/Windows/macOS)
curl -O https://github.com/mechaniel-coder/netwatch/releases/latest/download/netwatch-agent-linux-amd64

chmod +x netwatch-agent-linux-amd64
./netwatch-agent-linux-amd64 \
  --server https://your-netwatch-server:8080 \
  --token YOUR_AGENT_TOKEN \
  --interval 15s
```

---

## Project Structure

```
netwatch/
├── server/                         # Central API server (Go)
│   ├── main.go                     # Entry point — router init, DB connect, WS hub start
│   ├── go.mod
│   ├── Dockerfile
│   └── internal/
│       ├── api/
│       │   ├── router.go           # Route registration
│       │   └── handlers/
│       │       ├── agents.go       # Agent registration + heartbeat
│       │       ├── metrics.go      # Metric ingestion + query
│       │       └── alerts.go       # Alert rule CRUD + trigger log
│       ├── db/
│       │   ├── postgres.go         # Connection pool, pgx config
│       │   └── migrations/
│       │       └── 001_init.sql    # Schema: agents, metrics, alert_rules, alert_events
│       ├── models/
│       │   ├── agent.go
│       │   ├── metric.go
│       │   └── alert.go
│       ├── ws/
│       │   └── hub.go              # WebSocket broadcast hub
│       └── config/
│           └── config.go           # Env-based config with boot-time validation
│
├── agent/                          # Lightweight host agent (Go)
│   ├── main.go                     # Entry point — collect → report loop
│   ├── go.mod
│   ├── Dockerfile
│   └── internal/
│       ├── collector/
│       │   ├── cpu.go              # CPU usage via gopsutil
│       │   ├── memory.go           # RAM + swap
│       │   ├── disk.go             # Per-mount disk usage + I/O
│       │   └── network.go          # Interface bytes/packets
│       └── reporter/
│           └── http.go             # HTTPS POST to server with retry + backoff
│
├── dashboard/                      # React frontend
│   ├── src/
│   │   ├── App.tsx
│   │   ├── components/
│   │   │   ├── MetricsCard.tsx     # Single-metric stat card
│   │   │   ├── AgentList.tsx       # All agents + online/offline status
│   │   │   ├── AlertBanner.tsx     # Active alert notification bar
│   │   │   └── LiveChart.tsx       # Recharts real-time line chart
│   │   ├── pages/
│   │   │   ├── Dashboard.tsx       # Overview — all agents, top metrics
│   │   │   └── Agents.tsx          # Single-agent drilldown
│   │   ├── store/
│   │   │   └── useMetricsStore.ts  # Zustand — WS state, agent list
│   │   └── api/
│   │       └── client.ts           # Typed fetch wrapper
│   └── Dockerfile
│
├── docker-compose.yml
├── .env.example
└── .gitignore
```

---

## API Reference

| Method | Endpoint                         | Description                        |
|--------|----------------------------------|------------------------------------|
| POST   | `/api/v1/agents/register`        | Agent self-registration            |
| POST   | `/api/v1/agents/:id/heartbeat`   | Liveness update                    |
| POST   | `/api/v1/agents/:id/metrics`     | Ingest a metric snapshot           |
| GET    | `/api/v1/agents`                 | List all agents + latest status    |
| GET    | `/api/v1/agents/:id/metrics`     | Query historical metrics           |
| GET    | `/api/v1/alerts`                 | List alert rules                   |
| POST   | `/api/v1/alerts`                 | Create alert rule                  |
| GET    | `/api/v1/alerts/events`          | Alert trigger history              |
| GET    | `/ws/metrics`                    | WebSocket — live metric stream     |

---

## Configuration

Copy `.env.example` to `.env`:

```env
# Server
SERVER_PORT=8080
JWT_SECRET=change-me-in-production

# PostgreSQL
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DB=netwatch
POSTGRES_USER=netwatch
POSTGRES_PASSWORD=change-me

# Redis
REDIS_URL=redis://redis:6379

# Agent token (shared secret for agent auth; rotate per-environment)
AGENT_TOKEN_SECRET=change-me-in-production

# Metric retention (days)
RETENTION_DAYS=30
```

---

## Roadmap

- [ ] Per-agent alert rule overrides
- [ ] SMTP + webhook notification channels
- [ ] Agent binary auto-update via server
- [x] Prometheus metrics endpoint (`/metrics` — no external library)
- [x] Grafana dashboard (auto-provisioned via docker-compose)
- [ ] Kubernetes deployment manifest
- [ ] Role-based access (admin / read-only viewer)
- [ ] Dark / light theme toggle

---

## Observability Stack

NetWatch ships with a pre-wired Prometheus + Grafana stack. No configuration required — just run `docker compose up` and Grafana auto-provisions the datasource and dashboard.

| Service    | URL                        | Default Credentials |
|------------|----------------------------|---------------------|
| Dashboard  | http://localhost:3000      | —                   |
| API        | http://localhost:8080      | —                   |
| Prometheus | http://localhost:9090      | —                   |
| Grafana    | http://localhost:3001      | admin / admin       |

The `/metrics` endpoint is implemented using only the Go standard library (`sync/atomic`, `fmt`, `net/http`) — no external Prometheus client dependency. Prometheus scrapes it every 15 seconds.

**Metrics exposed:**

| Metric | Type | Description |
|--------|------|-------------|
| `netwatch_uptime_seconds` | Gauge | Seconds since server started |
| `netwatch_agent_registrations_total` | Counter | Total agent registrations |
| `netwatch_metrics_ingested_total` | Counter | Total metric snapshots ingested |
| `netwatch_websocket_connections` | Gauge | Active dashboard WebSocket connections |
| `netwatch_alerts_fired_total` | Counter | Total alert events fired |

---

## Technical Highlights

**Single-binary agent (~6 MB)**
The Go agent is a single statically-linked binary with zero runtime dependencies. It uses `gopsutil` for cross-platform CPU/memory/disk/network readings and `net/http` for shipping data to the server. Deploying to a new host is `scp agent binary` + run.

**WebSocket hub architecture**
The server maintains a central `Hub` struct that manages all active WebSocket connections. When any agent posts a metric snapshot, the API handler calls `hub.Broadcast(payload)` — every connected dashboard client receives the update instantly without polling. The hub uses a `register`/`unregister` channel pattern (not a mutex) to avoid race conditions at high connection counts.

**Prometheus `/metrics` without client_golang**
Rather than adding `prometheus/client_golang` as a dependency (which pulls in ~10 transitive deps), the `/metrics` endpoint manually writes Prometheus text exposition format using `fmt.Fprintf` and `sync/atomic` load operations. The format is simple enough that the standard library is sufficient and the result is a fully valid scrape target.

**Graceful shutdown**
The server listens for `SIGINT`/`SIGTERM` and calls `http.Server.Shutdown(ctx)` with a 10-second context. This drains in-flight requests before exit — important for a metric ingest server where losing a batch means a gap in your time-series.

---

## What I Learned Building This

- The WebSocket hub register/unregister channel pattern (instead of mutexes) is a pattern from the Gorilla WebSocket examples but understanding *why* it avoids data races required tracing through what happens when `Broadcast` is called concurrently with `register`.
- Implementing the Prometheus text exposition format by hand forced me to read the actual spec — I learned the difference between `counter` and `gauge` at the protocol level, not just the conceptual level.
- Go's `sync/atomic` package is significantly simpler than I expected for building lock-free metrics — `AddInt64` and `LoadInt64` cover the vast majority of real instrumentation needs.

---

## License

MIT © Daniel Pall
