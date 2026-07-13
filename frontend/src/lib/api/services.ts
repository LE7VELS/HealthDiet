import { getMockAppBootstrap } from '../../mocks/service'
import type { AppBootstrap } from '../../types'
import { apiRequest } from './client'
import { apiConfig } from './config'

export function getAppBootstrap(): Promise<AppBootstrap> {
  if (apiConfig.useMocks) {
    return getMockAppBootstrap()
  }

  return apiRequest<AppBootstrap>('/bootstrap')
}
