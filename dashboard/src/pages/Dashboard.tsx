import { useEffect } from 'react'
import { useMetricsStore, type MetricSnapshot } from '../store/useMetricsStore'
import { fetchAgents, openMetricsSocket } from '../api/client'
import AgentList from '../components/AgentList'
import AlertBanner from '../components/AlertBanner'

export default function Dashboard() {
  const { setAgents, applySnapshot, setConnected } = useMetricsStore()

  // Load agent list on mount
  useEffect(() => {
    fetchAgents().then(setAgents).catch(console.error)
  }, [setAgents])

  // Open WebSocket for live metric updates
  useEffect(() => {
    const ws = openMetricsSocket((data) => {
      applySnapshot(data as MetricSnapshot)
    })
    ws.onopen  = () => setConnected(true)
    ws.onclose = () => setConnected(false)
    return () => ws.close()
  }, [applySnapshot, setConnected])

  return (
    <div className="space-y-6">
      <AlertBanner />
      <h1 className="text-xl font-bold text-green-400 tracking-wide">// INFRASTRUCTURE OVERVIEW</h1>
      <AgentList />
    </div>
  )
}
