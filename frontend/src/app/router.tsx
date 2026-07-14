import { Navigate, createBrowserRouter } from 'react-router-dom'
import { AppShell } from './shell/AppShell'
import { ProtectedRoute } from './shell/ProtectedRoute'
import { PublicOnlyRoute } from './shell/PublicOnlyRoute'
import { PlaceholderPage } from '../pages/PlaceholderPage'
import { LoginPage } from '../pages/LoginPage'
import { RegisterPage } from '../pages/RegisterPage'
import { DashboardPage } from '../pages/DashboardPage'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Navigate to="/login" replace />,
  },
  {
    element: <PublicOnlyRoute />,
    children: [
      { path: '/login', element: <LoginPage /> },
      { path: '/register', element: <RegisterPage /> },
    ],
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        path: '/app',
        element: <AppShell />,
        children: [
          { index: true, element: <Navigate to="dashboard" replace /> },
          { path: 'dashboard', element: <DashboardPage /> },
          { path: 'meals', element: <PlaceholderPage title="饮食记录" /> },
          { path: 'meals/new', element: <PlaceholderPage title="新建饮食记录" /> },
          { path: 'foods', element: <PlaceholderPage title="食品与菜谱库" /> },
          { path: 'foods/:id', element: <PlaceholderPage title="食品 / 菜谱详情" /> },
          { path: 'reports', element: <PlaceholderPage title="营养报告" /> },
          { path: 'recommendations', element: <PlaceholderPage title="膳食建议" /> },
          { path: 'profile', element: <PlaceholderPage title="个人档案" /> },
        ],
      },
    ],
  },
  {
    path: '*',
    element: <PlaceholderPage title="页面不存在" description="请通过应用导航访问已规划页面。" />,
  },
])
