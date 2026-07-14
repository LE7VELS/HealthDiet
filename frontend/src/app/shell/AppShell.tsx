import { useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { NavLink, Outlet, useLocation, useNavigate } from 'react-router-dom'
import {
  Badge,
  ApplicationShell,
  DropdownMenu,
  DropdownMenuItem,
  Heading,
  HStack,
  Icon,
  SideNav,
  SideNavItem,
  Text,
  VStack,
  useAppToast,
} from '../../components/ui'
import { logoutUser } from '../../lib/api'
import { clearDemoSession } from '../../lib/storage/session'
import { pageTitles, primaryNavigation } from '../route-config'

const navigationIcons = [
  'info',
  'check',
  'search',
  'viewColumns',
  'success',
  'wrench',
] as const

function MobileNavigationLinks() {
  return primaryNavigation
    .filter((item) => item.path !== '/app/recommendations')
    .map((item) => (
      <NavLink
        className={({ isActive }) => `shell__nav-link${isActive ? ' shell__nav-link--active' : ''}`}
        key={item.path}
        to={item.path}
      >
        <span aria-hidden="true" className="shell__nav-dot" />
        <span>{item.shortLabel}</span>
      </NavLink>
    ))
}

export function AppShell() {
  const [collapsed, setCollapsed] = useState(false)
  const [isLoggingOut, setIsLoggingOut] = useState(false)
  const location = useLocation()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const showToast = useAppToast()
  const title = location.pathname.startsWith('/app/foods/')
    ? '食品 / 菜谱详情'
    : (pageTitles[location.pathname] ?? '智能膳食助手')

  async function handleLogout(): Promise<void> {
    setIsLoggingOut(true)
    try {
      await logoutUser()
    } catch {
      // 即使远端会话注销失败，也要清除本地会话，避免用户被困在登录状态。
    } finally {
      clearDemoSession()
      queryClient.clear()
      showToast({
        body: '你已安全退出登录。',
        autoHideDuration: 3500,
        uniqueID: 'logout-success',
      })
      navigate('/login', { replace: true })
    }
  }

  return (
    <ApplicationShell
      className={`shell${collapsed ? ' shell--collapsed' : ''}`}
      contentPadding={0}
      height="auto"
      mobileNav={false}
      sideNav={(
        <SideNav
          className="shell__side-nav"
          collapsible={{
            buttonLabel: collapsed ? '展开侧边栏' : '收起侧边栏',
            isCollapsed: collapsed,
            onCollapsedChange: setCollapsed,
          }}
          header={(
            <HStack align="center" className="shell__brand" gap={2}>
              <Badge className="shell__brand-mark" label="食" variant="green" />
              {!collapsed && <Text weight="bold">智能膳食助手</Text>}
            </HStack>
          )}
        >
          {primaryNavigation.map((item, index) => (
            <SideNavItem
              icon={<Icon icon={navigationIcons[index] ?? 'info'} />}
              isSelected={location.pathname === item.path}
              key={item.path}
              label={item.label}
              onClick={() => navigate(item.path)}
            />
          ))}
        </SideNav>
      )}
      topNav={(
        <header className="shell__header">
          <VStack gap={0.5}>
            <Text className="shell__eyebrow" color="secondary" type="supporting">
              智能膳食助手
            </Text>
            <Heading className="shell__page-title" level={1}>{title}</Heading>
          </VStack>

          <DropdownMenu
            button={{
              className: 'shell__user-button',
              icon: <Badge label="演" variant="green" />,
              isIconOnly: true,
              label: '打开演示用户菜单',
              size: 'md',
              tooltip: '账户菜单',
              variant: 'ghost',
            }}
            className="shell__user-menu-panel"
            hasChevron={false}
            menuWidth={252}
            placement="below"
          >
            <HStack align="center" className="shell__user-menu-account" gap={2}>
              <Badge label="演" variant="green" />
              <VStack gap={0.5}>
                <Text type="label" weight="bold">演示用户</Text>
                <Text color="secondary" type="supporting">demo@example.com</Text>
              </VStack>
            </HStack>
            <Badge className="shell__session-badge" label="演示会话" variant="green" />
            <DropdownMenuItem
              description="查看和完善账号资料"
              icon={<Icon icon="wrench" />}
              label="个人档案"
              onClick={() => navigate('/app/profile')}
            />
            <DropdownMenuItem
              description="清除本地演示会话"
              icon={<Icon icon="close" />}
              isDisabled={isLoggingOut}
              label={isLoggingOut ? '正在退出…' : '退出登录'}
              onClick={() => void handleLogout()}
            />
          </DropdownMenu>
        </header>
      )}
      variant="section"
    >
      <div className="shell__main">
        <Outlet />
      </div>

      {/* Astryx MobileNav 是抽屉式导航；产品要求固定底部五项，因此保留路由底栏。 */}
      <nav aria-label="移动端主导航" className="shell__mobile-nav">
        <MobileNavigationLinks />
      </nav>
    </ApplicationShell>
  )
}
