const BASE = import.meta.env.VITE_API_URL ?? ''
const WS_BASE = import.meta.env.VITE_WS_URL ?? ''

export async function fetchAgents() {
  const res = await fetch(`${BASE}/api/v1/agents`)
  if (!res.ok) throw new Error(`fetchAgents: ${res.status}`)
  return res.json()
}

export async function fetchMetrics(agentId: string, from?: string, to?: string) {
  const params = new URLSearchParams({ agent_id: agentId, limit: '200' })
  if (from) params.set('from', from)
  if (to)   params.set('to', to)
  const res = await fetch(`${BASE}/api/v1/agents/${agentId}/metrics?${params}`)
  if (!res.ok) throw new Error(`fetchMetrics: ${res.status}`)
  return res.json()
}

export function openMetricsSocket(onMessage: (data: unknown) => void): WebSocket {
  const ws = new WebSocket(`${WS_BASE}/ws/metrics`)
  ws.onmessage = (e) => {
    try { onMessage(JSON.parse(e.data)) } catch {}
  }
  return ws
}
