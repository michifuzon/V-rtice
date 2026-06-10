import { Navigate } from 'react-router-dom'

export default function AdminRoute({ children }) {
  const token = localStorage.getItem('token')
  const usuario = JSON.parse(localStorage.getItem('usuario') || 'null')
  if (!token || usuario?.rol !== 'admin') return <Navigate to="/" replace />
  return children
}
