import type { AppBootstrap } from '../types'
import { mockAppBootstrap } from './data'

const MOCK_DELAY_MS = 250

export async function getMockAppBootstrap(): Promise<AppBootstrap> {
  await new Promise((resolve) => window.setTimeout(resolve, MOCK_DELAY_MS))
  return structuredClone(mockAppBootstrap)
}
