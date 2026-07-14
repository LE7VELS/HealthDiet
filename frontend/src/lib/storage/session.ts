const DEMO_SESSION_KEY = 'smart-diet-demo-session'

// “保持登录”使用 localStorage；普通登录仅在当前浏览器会话内有效。
export function activateDemoSession(sessionId: string, remember = false): void {
  clearDemoSession()
  const storage = remember ? window.localStorage : window.sessionStorage
  storage.setItem(DEMO_SESSION_KEY, sessionId)
}

export function isDemoSessionAvailable(): boolean {
  return Boolean(
    window.sessionStorage.getItem(DEMO_SESSION_KEY)
      ?? window.localStorage.getItem(DEMO_SESSION_KEY),
  )
}

export function clearDemoSession(): void {
  // 同时清理两种存储，避免用户切换“保持登录”后残留旧会话。
  window.sessionStorage.removeItem(DEMO_SESSION_KEY)
  window.localStorage.removeItem(DEMO_SESSION_KEY)
}
