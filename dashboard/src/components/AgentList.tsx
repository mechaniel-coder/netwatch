import { useMetricsStore } from '../store/useMetricsStore'
import MetricsCard from './MetricsCard'

export default function AgentList() {
  const { agents, latestMetrics } = useMetricsStore()

  if (agents.length === 0) {
    return <p className="text-gray-600 text-sm">No agents registered yet. Deploy an agent to get started.</p>
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      {agents.map((agent) => (
        <div key={agent.id} className="border border-gray-800 rounded-lg p-4 bg-gray-900">
          <div className="flex items-center justify-between mb-3">
            <div>
              <p className="text-white font-semibold">{agent.hostname}</p>
              <p className="text-gray-500 text-xs">{agent.ip_address} · {agent.os}</p>
            </div>
            <span className={`text-xs px-2 py-1 rounded-full font-medium ${
              agent.online ? 'bg-green-900 text-green-300' : 'bg-red-900 text-red-400'
            }`}>
              {agent.online ? 'ONLINE' : 'OFFLINE'}
            </span>
          </div>

          {latestMetrics[agent.id] ? (
            <div className="grid grid-cols-2 gap-2">
              <MetricsCard label="CPU" value={latestMetrics[agent.id].cpu_percent} unit="%" />
              <MetricsCard label="RAM" value={toGB(latestMetrics[agent.id].mem_used_bytes)} unit="GB" />
              <MetricsCard label="DISK" value={toGB(latestMetrics[agent.id].disk_used_bytes)} unit="GB" />
              <MetricsCard label="NET TX" value={toMB(latestMetrics[agent.id].net_tx_bytes)} unit="MB" />
            </div>
          ) : (
            <p className="text-gray-700 text-xs">Waiting for first metric snapshot...</p>
          )}
        </div>
      ))}
    </div>
  )
}

const toGB = (bytes: number) => (bytes / 1_073_741_824).toFixed(1)
const toMB = (bytes: number) => (bytes / 1_048_576).toFixed(1)
