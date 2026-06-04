interface MetricsCardProps {
  label: string
  value: string | number
  unit: string
}

export default function MetricsCard({ label, value, unit }: MetricsCardProps) {
  return (
    <div className="bg-gray-950 border border-gray-800 rounded p-2">
      <p className="text-gray-600 text-xs tracking-widest">{label}</p>
      <p className="text-white font-mono text-lg font-bold leading-none mt-1">
        {value}<span className="text-gray-500 text-xs ml-1">{unit}</span>
      </p>
    </div>
  )
}
