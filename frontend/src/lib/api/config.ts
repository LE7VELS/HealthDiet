const rawApiBaseUrl = import.meta.env.VITE_API_BASE_URL?.trim() ?? ''

export const apiConfig = {
  baseUrl: rawApiBaseUrl.replace(/\/$/, ''),
  // 当前阶段默认启用 Mock；只有显式配置 false 时才请求真实后端。
  useMocks: import.meta.env.VITE_USE_MOCKS !== 'false',
} as const
