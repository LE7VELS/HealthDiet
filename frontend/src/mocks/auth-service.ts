import { ApiError } from '../lib/api/client'
import type { RegisterRequest, RegisterResponse } from '../types'

const registeredEmails = new Set(['demo@example.com', 'existing@example.com'])
const MOCK_DELAY_MS = 650

export async function registerMockUser(request: RegisterRequest): Promise<RegisterResponse> {
  await new Promise((resolve) => window.setTimeout(resolve, MOCK_DELAY_MS))

  const normalizedEmail = request.email.trim().toLowerCase()
  if (registeredEmails.has(normalizedEmail)) {
    throw new ApiError('该邮箱已注册，请直接登录或更换邮箱。', 409)
  }

  registeredEmails.add(normalizedEmail)

  return {
    sessionId: crypto.randomUUID(),
    user: {
      id: crypto.randomUUID(),
      username: request.username,
      email: normalizedEmail,
    },
  }
}
