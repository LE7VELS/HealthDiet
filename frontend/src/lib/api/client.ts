import axios, { type AxiosRequestConfig } from 'axios'
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

type ApiErrorResponse = {
  error?: { code?: string; message?: string; fields?: ApiFieldError[] }
}

// apiRequest 是页面访问后端的统一 Axios 边界：集中拼接基础地址、附加 JWT，并把网络或 API 错误
// 转换为稳定的 ApiError。业务模块只提供相对路径和 Axios 配置，不直接处理会话失效跳转。
export async function apiRequest<T>(path: string, config?: AxiosRequestConfig): Promise<T> {
  const accessToken = getAccessToken()
  try {
    // Axios 自动序列化 data 中的普通对象并解析 JSON 响应，调用方无需手动 JSON.stringify。
    const response = await axios.request<T>({
      ...config,
      baseURL: apiConfig.baseUrl,
      url: path,
      headers: {
        Accept: 'application/json',
        ...config?.headers,
        // Authorization 必须由统一会话层最后写入，业务调用方不能用自定义 Header 覆盖当前用户身份。
        ...(accessToken ? { Authorization: `Bearer ${accessToken}` } : {}),
      },
    })

    // Axios 对 204 返回空字符串；统一客户端合同仍以 void 表示无响应体。
    if (response.status === 204) {
      return undefined as T
    }
    return response.data
  } catch (cause) {
    if (axios.isAxiosError<ApiErrorResponse>(cause)) {
      const status = cause.response?.status ?? 0
      const error = cause.response?.data?.error
      if (status === 401 && error?.code === 'UNAUTHENTICATED') {
        clearSession()
        if (!window.location.pathname.startsWith('/login')) {
          window.location.assign('/login')
        }
      }
      throw new ApiError(
        error?.message ?? (cause.response ? '请求失败，请稍后重试。' : '无法连接服务器，请检查网络后重试。'),
        status,
        error?.code,
        error?.fields,
      )
    }
    throw cause
  }
}
