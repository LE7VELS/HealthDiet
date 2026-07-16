import { useQuery } from '@tanstack/react-query'
import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { Banner, Button, Surface, Text } from '../../components/ui'
import { ApiError, getCurrentUser } from '../../lib/api'
import {
  clearSession,
  isSessionAvailable,
  updateSessionUser,
} from '../../lib/storage/session'

export function ProtectedRoute() {
  const location = useLocation()
  const hasLocalSession = isSessionAvailable()
  const currentUserQuery = useQuery({
    queryKey: ['auth', 'me'],
    enabled: hasLocalSession,
    retry: false,
    staleTime: 60_000,
    // 在查询成功返回前同步可信用户信息，确保随后挂载的 AppShell 读取到最新会话。
    queryFn: async () => {
      const user = await getCurrentUser()
      updateSessionUser(user)
      return user
    },
  })

  if (!hasLocalSession) {
    // 保存原始地址，登录成功后返回用户原本要访问的业务页面。
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  if (currentUserQuery.isPending) {
    return (
      <main aria-live="polite" className="route-status-page">
        <Surface className="route-status-card">
          <Text as="p" color="secondary">正在验证登录状态…</Text>
        </Surface>
      </main>
    )
  }

  if (currentUserQuery.isError) {
    if (currentUserQuery.error instanceof ApiError && currentUserQuery.error.status === 401) {
      clearSession()
      return <Navigate to="/login" replace state={{ from: location }} />
    }

    // 临时网络或服务端错误不等同于退出登录，保留本地会话并允许用户主动重试。
    return (
      <main className="route-status-page">
        <Surface className="route-status-card">
          <Banner
            description={currentUserQuery.error instanceof Error
              ? currentUserQuery.error.message
              : '暂时无法验证登录状态，请稍后重试。'}
            status="error"
            title="无法进入应用"
          />
          <Button
            isDisabled={currentUserQuery.isFetching}
            isLoading={currentUserQuery.isFetching}
            label={currentUserQuery.isFetching ? '正在重试' : '重新验证'}
            onClick={() => void currentUserQuery.refetch()}
            variant="primary"
          />
        </Surface>
      </main>
    )
  }

  return <Outlet />
}
