const rawApiBaseUrl = import.meta.env.VITE_API_BASE_URL?.trim() ?? ''

export const apiConfig = {
  baseUrl: (rawApiBaseUrl || '/api/v1').replace(/\/$/, ''),
  // 真实 Go API 是默认主流程；Mock 仅在开发者显式开启时使用。
  useMocks: import.meta.env.VITE_USE_MOCKS === 'true',
} as const
