import { Navigate, Outlet } from 'react-router-dom'
import { isSessionAvailable } from '../../lib/storage/session'

export function PublicOnlyRoute() {
  // 已登录用户无需再次进入登录或注册流程。
  if (isSessionAvailable()) {
    return <Navigate to="/app/dashboard" replace />
  }

  return <Outlet />
}
