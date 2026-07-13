const DEMO_SESSION_KEY = 'smart-diet-demo-session'

export function isDemoSessionAvailable(): boolean {
  if (import.meta.env.VITE_USE_MOCKS === 'true') {
    return true
  }

  return window.sessionStorage.getItem(DEMO_SESSION_KEY) === 'active'
}
