import {
  getCurrentMockUser,
  loginMockUser,
  logoutMockUser,
  registerMockUser,
} from '../../mocks/auth-service'
import type {
  AuthUser,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
} from '../../types'
import { apiRequest } from './client'
import { apiConfig } from './config'

type DataResponse<T> = { data: T }

export function registerUser(request: RegisterRequest): Promise<RegisterResponse> {
  // Mock 与真实请求在 API 层切换，表单组件不感知具体数据来源。
  if (apiConfig.useMocks) {
    return registerMockUser(request)
  }

  return apiRequest<DataResponse<RegisterResponse>>('/auth/register', {
    method: 'POST',
    data: request,
  }).then((response) => response.data)
}

export function loginUser(request: LoginRequest): Promise<LoginResponse> {
  // 保持 Mock 和真实后端返回同一 DTO，便于后续关闭 Mock 时不改页面逻辑。
  if (apiConfig.useMocks) {
    return loginMockUser(request)
  }

  return apiRequest<DataResponse<LoginResponse>>('/auth/login', {
    method: 'POST',
    data: request,
  }).then((response) => response.data)
}

// getCurrentUser 使用当前请求客户端自动附带的 Bearer Token 验证会话并取得当前账号。
export function getCurrentUser(): Promise<AuthUser> {
  if (apiConfig.useMocks) {
    return getCurrentMockUser()
  }

  return apiRequest<DataResponse<AuthUser>>('/auth/me').then((response) => response.data)
}

export function logoutUser(): Promise<void> {
  if (apiConfig.useMocks) {
    return logoutMockUser()
  }

  // 当前合同使用无状态访问 Token，不存在服务端注销接口。
  return Promise.resolve()
}
