import { Link } from 'react-router-dom'
import { Surface } from '../components/ui/Surface'

type PlaceholderPageProps = {
  title: string
  description?: string
}

export function PlaceholderPage({
  title,
  description = '页面路由与布局已就绪，正式业务内容将在后续阶段实现。',
}: PlaceholderPageProps) {
  const isPublicPage = title === '登录' || title === '注册' || title === '页面不存在'

  const content = (
    <Surface className="placeholder">
      <span className="placeholder__badge">初始化占位页</span>
      <h2>{title}</h2>
      <p>{description}</p>
      {title === '登录' && <Link to="/register">前往注册占位页</Link>}
      {title === '注册' && <Link to="/login">返回登录占位页</Link>}
      {title === '页面不存在' && <Link to="/login">返回登录页</Link>}
    </Surface>
  )

  if (!isPublicPage) {
    return content
  }

  return (
    <main className="public-page">
      <div className="public-page__brand">智能膳食助手</div>
      {content}
    </main>
  )
}
