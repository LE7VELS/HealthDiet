import { ApiError } from '../lib/api/client'
import { getSession } from '../lib/storage/session'
import type { AuthUser, LoginRequest, LoginResponse, RegisterRequest, RegisterResponse } from '../types'

type MockAccount = {
  id: string
  username: string
  email: string
  password: string
}

// Mock 账号只存在于当前前端运行周期，不模拟真实持久化或安全认证。
const mockAccounts: MockAccount[] = [
  {
    id: 'demo-user',
    username: 'demo',
    email: 'demo@example.com',
    password: 'Demo1234',
  },
  {
    id: 'existing-user',
    username: 'existing',
    email: 'existing@example.com',
    password: 'Existing1234',
  },
]
const MOCK_DELAY_MS = 650

async function simulateDelay(): Promise<void> {
  await new Promise((resolve) => window.setTimeout(resolve, MOCK_DELAY_MS))
}

export async function loginMockUser(request: LoginRequest): Promise<LoginResponse> {
  await simulateDelay()

  const normalizedIdentifier = request.identifier.trim().toLowerCase()
  const account = mockAccounts.find(
    (candidate) => candidate.email === normalizedIdentifier
      || candidate.username.toLowerCase() === normalizedIdentifier,
  )

  if (!account || account.password !== request.password) {
    throw new ApiError('账号或密码不正确，请检查后重试。', 401)
  }

  return {
    accessToken: crypto.randomUUID(),
    tokenType: 'Bearer',
    expiresIn: 3600,
    user: {
      id: account.id,
      username: account.username,
      email: account.email,
    },
  }
}

export async function logoutMockUser(): Promise<void> {
  await simulateDelay()
}

// Mock 通过现有本地会话返回同形 DTO，仅用于显式开启 Mock 时保持路由守卫行为一致。
export async function getCurrentMockUser(): Promise<AuthUser> {
  await simulateDelay()
  const user = getSession()?.user
  if (!user) {
    throw new ApiError('登录状态无效或已过期', 401, 'UNAUTHENTICATED')
  }
  return user
}

export async function registerMockUser(request: RegisterRequest): Promise<RegisterResponse> {
  await simulateDelay()

  const normalizedEmail = request.email.trim().toLowerCase()
  if (mockAccounts.some((account) => account.email === normalizedEmail)) {
    throw new ApiError('该邮箱已注册，请直接登录或更换邮箱。', 409)
  }

  const account: MockAccount = {
    id: crypto.randomUUID(),
    username: request.username.trim(),
    email: normalizedEmail,
    password: request.password,
  }
  mockAccounts.push(account)

  return {
    accessToken: crypto.randomUUID(),
    tokenType: 'Bearer',
    expiresIn: 3600,
    user: {
      id: account.id,
      username: account.username,
      email: normalizedEmail,
    },
  }
}
