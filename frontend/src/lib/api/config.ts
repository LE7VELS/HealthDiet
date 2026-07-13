const rawApiBaseUrl = import.meta.env.VITE_API_BASE_URL?.trim() ?? ''

export const apiConfig = {
  baseUrl: rawApiBaseUrl.replace(/\/$/, ''),
  useMocks: import.meta.env.VITE_USE_MOCKS === 'true',
} as const
