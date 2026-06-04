-- NetWatch initial schema
-- Idempotent: safe to run on an already-initialized database.

CREATE TABLE IF NOT EXISTS agents (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hostname    TEXT NOT NULL,
    ip_address  TEXT NOT NULL,
    os          TEXT NOT NULL,
    version     TEXT NOT NULL,
    last_seen   TIMESTAMPTZ,
    online      BOOLEAN NOT NULL DEFAULT false,
    registered_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_agents_online ON agents (online);

-- ─────────────────────────────────────────────────────────────────────────────
-- Metrics: one row per reporting interval per agent.
-- Partitioning by month is recommended for high-volume deployments.
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS metrics (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    recorded_at     TIMESTAMPTZ NOT NULL DEFAULT now(),

    cpu_percent     NUMERIC(5,2),
    mem_used_bytes  BIGINT,
    mem_total_bytes BIGINT,
    disk_used_bytes BIGINT,
    disk_total_bytes BIGINT,
    net_rx_bytes    BIGINT,
    net_tx_bytes    BIGINT
);

CREATE INDEX IF NOT EXISTS idx_metrics_agent_time
    ON metrics (agent_id, recorded_at DESC);

-- ─────────────────────────────────────────────────────────────────────────────
-- Alert rules: per-agent or global threshold definitions.
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS alert_rules (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id    UUID REFERENCES agents(id) ON DELETE CASCADE,  -- NULL = applies to all agents
    metric      TEXT NOT NULL,          -- e.g. 'cpu_percent', 'mem_used_bytes'
    operator    TEXT NOT NULL,          -- 'gt' | 'lt' | 'gte' | 'lte'
    threshold   NUMERIC NOT NULL,
    severity    TEXT NOT NULL DEFAULT 'warning',  -- 'warning' | 'critical'
    enabled     BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ─────────────────────────────────────────────────────────────────────────────
-- Alert events: log of every alert that fired.
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS alert_events (
    id          BIGSERIAL PRIMARY KEY,
    rule_id     UUID NOT NULL REFERENCES alert_rules(id) ON DELETE CASCADE,
    agent_id    UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    value       NUMERIC NOT NULL,
    severity    TEXT NOT NULL,
    fired_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_alert_events_agent
    ON alert_events (agent_id, fired_at DESC);
