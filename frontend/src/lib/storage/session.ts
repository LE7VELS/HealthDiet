import type { AuthUser, LoginResponse } from '../../types'

const SESSION_KEY = 'healthdiet-auth-session'

export type AuthSession = Pick<LoginResponse, 'accessToken' | 'expiresIn' | 'user'> & {
  createdAt: number
}

// “保持登录”使用 localStorage；普通登录仅在当前浏览器会话内有效。
export function activateSession(result: LoginResponse, remember = false): void {
  clearSession()
  const storage = remember ? window.localStorage : window.sessionStorage
  storage.setItem(SESSION_KEY, JSON.stringify({
    accessToken: result.accessToken,
    expiresIn: result.expiresIn,
    user: result.user,
    createdAt: Date.now(),
  } satisfies AuthSession))
}

export function getSession(): AuthSession | null {
  const raw = window.sessionStorage.getItem(SESSION_KEY)
    ?? window.localStorage.getItem(SESSION_KEY)
  if (!raw) return null

  try {
    const session = JSON.parse(raw) as AuthSession
    const expiresAt = session.createdAt + session.expiresIn * 1000
    if (!session.accessToken || expiresAt <= Date.now()) {
      clearSession()
      return null
    }
    return session
  } catch {
    clearSession()
    return null
  }
}

export function getAccessToken(): string | null {
  return getSession()?.accessToken ?? null
}

export function isSessionAvailable(): boolean {
  return getSession() !== null
}

// updateSessionUser 用 /auth/me 的可信响应刷新用户信息，同时保留原 Token、存储位置和过期时间。
// 这里不能重新调用 activateSession，否则每次身份验证都会错误地延长本地会话寿命。
export function updateSessionUser(user: AuthUser): void {
  const storage = window.sessionStorage.getItem(SESSION_KEY)
    ? window.sessionStorage
    : window.localStorage.getItem(SESSION_KEY)
      ? window.localStorage
      : null
  const session = getSession()
  if (!storage || !session) return

  storage.setItem(SESSION_KEY, JSON.stringify({ ...session, user } satisfies AuthSession))
}

export function clearSession(): void {
  // 同时清理两种存储，避免用户切换“保持登录”后残留旧会话。
  window.sessionStorage.removeItem(SESSION_KEY)
  window.localStorage.removeItem(SESSION_KEY)
}
