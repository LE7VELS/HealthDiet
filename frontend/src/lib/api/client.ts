import { apiConfig } from './config'
import { clearSession, getAccessToken } from '../storage/session'

export type ApiFieldError = { field: string; message: string }

export class ApiError extends Error {
  constructor(
    message: string,
    readonly status: number,
    readonly code = 'REQUEST_FAILED',
    readonly fields: ApiFieldError[] = [],
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export async function apiRequest<T>(path: string, init?: RequestInit): Promise<T> {
  const accessToken = getAccessToken()
  // 所有真实 API 请求从这里统一拼接基础地址，页面和业务模块不直接拼 URL。
  const response = await fetch(`${apiConfig.baseUrl}${path}`, {
    ...init,
    headers: {
      Accept: 'application/json',
      ...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}),
      ...init?.headers,
    },
  })

  if (!response.ok) {
    const body = await response.json().catch(() => null) as {
      error?: { code?: string; message?: string; fields?: ApiFieldError[] }
    } | null
    const error = body?.error
    if (response.status === 401 && error?.code === 'UNAUTHENTICATED') {
      clearSession()
      if (!window.location.pathname.startsWith('/login')) {
        window.location.assign('/login')
      }
    }
    throw new ApiError(
      error?.message ?? '请求失败，请稍后重试。',
      response.status,
      error?.code,
      error?.fields,
    )
  }

  if (response.status === 204) {
    // 204 没有响应体，跳过 JSON 解析；调用方应将对应返回类型声明为 void。
    return undefined as T
  }

  return response.json() as Promise<T>
}
