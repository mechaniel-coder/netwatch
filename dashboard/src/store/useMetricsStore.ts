import { create } from 'zustand'

export interface Agent {
  id: string
  hostname: string
  ip_address: string
  os: string
  version: string
  online: boolean
  last_seen: string | null
}

export interface MetricSnapshot {
  agent_id: string
  metrics: {
    cpu_percent: number
    mem_used_bytes: number
    mem_total_bytes: number
    disk_used_bytes: number
    disk_total_bytes: number
    net_rx_bytes: number
    net_tx_bytes: number
  }
}

interface MetricsState {
  agents: Agent[]
  latestMetrics: Record<string, MetricSnapshot['metrics']>
  connected: boolean
  setAgents: (agents: Agent[]) => void
  applySnapshot: (snap: MetricSnapshot) => void
  setConnected: (v: boolean) => void
}

export const useMetricsStore = create<MetricsState>((set) => ({
  agents: [],
  latestMetrics: {},
  connected: false,
  setAgents: (agents) => set({ agents }),
  applySnapshot: (snap) =>
    set((state) => ({
      latestMetrics: { ...state.latestMetrics, [snap.agent_id]: snap.metrics },
    })),
  setConnected: (connected) => set({ connected }),
}))
