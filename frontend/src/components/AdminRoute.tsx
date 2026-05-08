import { Navigate } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

export const AdminRoute = ({ children }: { children: React.ReactNode }) => {
  const { user, loading } = useAuth()

  if (loading) {
    return <div className="page-loading">Checking permissions...</div>
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  if (user.role !== 'admin') {
    return <Navigate to="/dashboard" replace />
  }

  return <>{children}</>
}
