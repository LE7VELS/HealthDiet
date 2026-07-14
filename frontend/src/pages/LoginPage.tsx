import { Link } from 'react-router-dom'
import { Badge, Heading, Surface, Text } from '../components/ui'
import { LoginForm } from '../features/auth/components/LoginForm'

export function LoginPage() {
  return (
    <main className="auth-page">
      <section className="auth-page__intro" aria-labelledby="login-intro-title">
        <Link className="auth-page__brand" to="/login" aria-label="智能膳食助手登录页">
          <span className="auth-page__brand-mark" aria-hidden="true">食</span>
          <span>智能膳食助手</span>
        </Link>

        <div className="auth-page__intro-copy">
          <Text className="auth-page__eyebrow" color="secondary" type="supporting">
            记录每一餐，看见每日变化
          </Text>
          <Heading id="login-intro-title" level={1} type="display-1">
            <span className="auth-page__title-line">从清晰记录开始</span>
            <span className="auth-page__title-line">管理日常饮食</span>
          </Heading>
          <Text as="p">登录后继续维护饮食记录、查看营养趋势，并逐步完善你的个人健康档案。</Text>
        </div>

        <ul className="auth-page__benefits" aria-label="产品特点">
          <li><span aria-hidden="true">01</span>快速记录每日餐食</li>
          <li><span aria-hidden="true">02</span>查看直观营养趋势</li>
          <li><span aria-hidden="true">03</span>获得可执行的膳食建议</li>
        </ul>
      </section>

      <section className="auth-page__form-area" aria-labelledby="login-title">
        <Surface className="register-card">
          <div className="register-card__header">
            <Badge className="register-card__step" label="登录个人账号" variant="green" />
            <Heading id="login-title" level={2}>欢迎回来</Heading>
            <Text as="p" color="secondary">使用你的用户名或邮箱继续访问智能膳食助手。</Text>
          </div>

          <LoginForm />

          <Text as="p" className="login-card__demo">
            演示账号：<strong>demo@example.com</strong><br />
            演示密码：<strong>Demo1234</strong>
          </Text>

          <Text as="div" className="register-card__footer" color="secondary">
            还没有账号？<Link to="/register">立即注册</Link>
          </Text>
        </Surface>
      </section>
    </main>
  )
}
