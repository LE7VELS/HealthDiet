const DEMO_SESSION_KEY = 'smart-diet-demo-session'

export function activateDemoSession(sessionId: string): void {
  window.sessionStorage.setItem(DEMO_SESSION_KEY, sessionId)
}

export function isDemoSessionAvailable(): boolean {
  if (import.meta.env.VITE_USE_MOCKS === 'true') {
    return true
  }

  return Boolean(window.sessionStorage.getItem(DEMO_SESSION_KEY))
}
