import { Routes, Route, NavLink } from 'react-router-dom'
import Dashboard from './pages/Dashboard'
import Agents from './pages/Agents'

export default function App() {
  return (
    <div className="min-h-screen bg-gray-950 text-gray-100 font-mono">
      <nav className="border-b border-gray-800 px-6 py-3 flex items-center gap-8">
        <span className="text-green-400 font-bold tracking-widest text-sm">NETWATCH</span>
        <NavLink to="/"       className={({ isActive }) => isActive ? 'text-white' : 'text-gray-500 hover:text-gray-300'}>Dashboard</NavLink>
        <NavLink to="/agents" className={({ isActive }) => isActive ? 'text-white' : 'text-gray-500 hover:text-gray-300'}>Agents</NavLink>
      </nav>

      <main className="p-6">
        <Routes>
          <Route path="/"        element={<Dashboard />} />
          <Route path="/agents"  element={<Agents />} />
          <Route path="/agents/:id" element={<Agents />} />
        </Routes>
      </main>
    </div>
  )
}
