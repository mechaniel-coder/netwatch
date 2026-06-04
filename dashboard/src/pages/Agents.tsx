import { useParams } from 'react-router-dom'
import { useMetricsStore } from '../store/useMetricsStore'

export default function Agents() {
  const { id } = useParams()
  const { agents, latestMetrics } = useMetricsStore()

  if (id) {
    const agent = agents.find((a) => a.id === id)
    const metrics = latestMetrics[id]
    // TODO: render per-agent historical chart using fetchMetrics(id)
    return (
      <div>
        <h1 className="text-xl font-bold text-green-400 mb-4">// AGENT: {agent?.hostname ?? id}</h1>
        <pre className="text-xs text-gray-400 bg-gray-900 p-4 rounded border border-gray-800">
          {JSON.stringify({ agent, metrics }, null, 2)}
        </pre>
      </div>
    )
  }

  return (
    <div>
      <h1 className="text-xl font-bold text-green-400 mb-4">// ALL AGENTS</h1>
      <table className="w-full text-sm border border-gray-800 rounded overflow-hidden">
        <thead className="bg-gray-900 text-gray-500 text-xs tracking-widest">
          <tr>
            {['HOSTNAME', 'IP', 'OS', 'VERSION', 'STATUS', 'LAST SEEN'].map((h) => (
              <th key={h} className="text-left px-4 py-2">{h}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {agents.map((a) => (
            <tr key={a.id} className="border-t border-gray-800 hover:bg-gray-900 transition-colors">
              <td className="px-4 py-2 text-white">{a.hostname}</td>
              <td className="px-4 py-2 text-gray-400">{a.ip_address}</td>
              <td className="px-4 py-2 text-gray-400">{a.os}</td>
              <td className="px-4 py-2 text-gray-400">{a.version}</td>
              <td className="px-4 py-2">
                <span className={`text-xs px-2 py-0.5 rounded-full ${a.online ? 'bg-green-900 text-green-300' : 'bg-red-900 text-red-400'}`}>
                  {a.online ? 'ONLINE' : 'OFFLINE'}
                </span>
              </td>
              <td className="px-4 py-2 text-gray-600 text-xs">{a.last_seen ?? '—'}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
