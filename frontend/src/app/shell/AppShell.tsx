import { useState } from 'react'
import { NavLink, Outlet, useLocation } from 'react-router-dom'
import { pageTitles, primaryNavigation } from '../route-config'

function NavigationLinks({ mobile = false }: { mobile?: boolean }) {
  const items = mobile
    ? primaryNavigation.filter((item) => item.path !== '/app/recommendations')
    : primaryNavigation

  return items.map((item) => (
    <NavLink
      className={({ isActive }) => `shell__nav-link${isActive ? ' shell__nav-link--active' : ''}`}
      key={item.path}
      to={item.path}
    >
      <span aria-hidden="true" className="shell__nav-dot" />
      <span>{mobile ? item.shortLabel : item.label}</span>
    </NavLink>
  ))
}

export function AppShell() {
  const [collapsed, setCollapsed] = useState(false)
  const location = useLocation()
  const title = location.pathname.startsWith('/app/foods/')
    ? '食品 / 菜谱详情'
    : (pageTitles[location.pathname] ?? '智能膳食助手')

  return (
    <div className={`shell${collapsed ? ' shell--collapsed' : ''}`}>
      <aside className="shell__sidebar">
        <div className="shell__brand">
          <span className="shell__brand-mark" aria-hidden="true">食</span>
          <span className="shell__brand-text">智能膳食助手</span>
        </div>
        <nav aria-label="主导航" className="shell__desktop-nav">
          <NavigationLinks />
        </nav>
        <button
          aria-label={collapsed ? '展开侧边栏' : '收起侧边栏'}
          className="shell__collapse-button"
          onClick={() => setCollapsed((value) => !value)}
          type="button"
        >
          {collapsed ? '›' : '‹ 收起'}
        </button>
      </aside>

      <div className="shell__body">
        <header className="shell__header">
          <div>
            <span className="shell__eyebrow">智能膳食助手</span>
            <h1>{title}</h1>
          </div>
          <button className="shell__user-button" type="button" aria-label="打开用户菜单">
            演示用户
          </button>
        </header>
        <main className="shell__main" id="main-content">
          <Outlet />
        </main>
      </div>

      <nav aria-label="移动端主导航" className="shell__mobile-nav">
        <NavigationLinks mobile />
      </nav>
    </div>
  )
}
