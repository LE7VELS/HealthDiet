import { loginMockUser, logoutMockUser, registerMockUser } from '../../mocks/auth-service'
import type { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse } from '../../types'
import { apiRequest } from './client'
import { apiConfig } from './config'

export function registerUser(request: RegisterRequest): Promise<RegisterResponse> {
  if (apiConfig.useMocks) {
    return registerMockUser(request)
  }

  return apiRequest<RegisterResponse>('/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  })
}

export function loginUser(request: LoginRequest): Promise<LoginResponse> {
  if (apiConfig.useMocks) {
    return loginMockUser(request)
  }

  return apiRequest<LoginResponse>('/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  })
}

export function logoutUser(): Promise<void> {
  if (apiConfig.useMocks) {
    return logoutMockUser()
  }

  return apiRequest<void>('/auth/logout', { method: 'POST' })
}
