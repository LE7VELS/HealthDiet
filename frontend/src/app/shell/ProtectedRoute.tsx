import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { isDemoSessionAvailable } from '../../lib/storage/session'

export function ProtectedRoute() {
  const location = useLocation()

  if (!isDemoSessionAvailable()) {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  return <Outlet />
}
