import { useMetricsStore } from '../store/useMetricsStore'

export default function AlertBanner() {
  const connected = useMetricsStore((s) => s.connected)

  return (
    <div className={`flex items-center gap-2 text-xs px-3 py-2 rounded border ${
      connected
        ? 'border-green-900 bg-green-950 text-green-400'
        : 'border-yellow-900 bg-yellow-950 text-yellow-400'
    }`}>
      <span className={`w-2 h-2 rounded-full ${connected ? 'bg-green-400 animate-pulse' : 'bg-yellow-400'}`} />
      {connected ? 'Live feed connected — metrics updating in real time' : 'Connecting to metric stream...'}
    </div>
  )
}
