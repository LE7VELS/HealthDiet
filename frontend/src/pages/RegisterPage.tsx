import { Link } from 'react-router-dom'
import { Surface } from '../components/ui'
import { RegisterForm } from '../features/auth/components/RegisterForm'

export function RegisterPage() {
  return (
    <main className="auth-page">
      <section className="auth-page__intro" aria-labelledby="register-intro-title">
        <Link className="auth-page__brand" to="/login" aria-label="智能膳食助手首页">
          <span className="auth-page__brand-mark" aria-hidden="true">食</span>
          <span>智能膳食助手</span>
        </Link>

        <div className="auth-page__intro-copy">
          <span className="auth-page__eyebrow">从今天开始，更了解每一餐</span>
          <h1 id="register-intro-title">让饮食记录成为一种轻松的日常</h1>
          <p>集中管理饮食、营养趋势和个人档案，为每一次选择提供清晰的膳食参考。</p>
        </div>

        <ul className="auth-page__benefits" aria-label="产品特点">
          <li><span aria-hidden="true">01</span>快速记录每日餐食</li>
          <li><span aria-hidden="true">02</span>查看直观营养趋势</li>
          <li><span aria-hidden="true">03</span>获得可执行的膳食建议</li>
        </ul>
      </section>

      <section className="auth-page__form-area" aria-labelledby="register-title">
        <Surface className="register-card">
          <div className="register-card__header">
            <span className="register-card__step">创建个人账号</span>
            <h2 id="register-title">欢迎加入</h2>
            <p>先完成基础注册，健康档案可以稍后继续完善。</p>
          </div>

          <RegisterForm />

          <div className="register-card__footer">
            已有账号？<Link to="/login">返回登录</Link>
          </div>
        </Surface>
      </section>
    </main>
  )
}
