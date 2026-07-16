import type { AppBootstrap } from '../types'
import { mockAppBootstrap } from './data'

const MOCK_DELAY_MS = 250

export async function getMockAppBootstrap(): Promise<AppBootstrap> {
  await new Promise((resolve) => window.setTimeout(resolve, MOCK_DELAY_MS))
  // 返回副本，避免页面修改数据时污染后续 Mock 请求共享的原始对象。
  return structuredClone(mockAppBootstrap)
}
