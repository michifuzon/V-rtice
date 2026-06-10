import { Navigate } from 'react-router-dom'

/**
 * Redirige a /login si no hay token en localStorage.
 */
export default function PrivateRoute({ children }) {
  const token = localStorage.getItem('token')
  return token ? children : <Navigate to="/login" replace />
}
