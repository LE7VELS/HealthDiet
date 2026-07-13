import { registerMockUser } from '../../mocks/auth-service'
import type { RegisterRequest, RegisterResponse } from '../../types'
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
