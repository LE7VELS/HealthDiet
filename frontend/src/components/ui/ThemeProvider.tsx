import { Theme } from '@astryxdesign/core/theme'
import { neutralTheme } from '@astryxdesign/theme-neutral/built'
import type { PropsWithChildren } from 'react'

export function AppThemeProvider({ children }: PropsWithChildren) {
  return (
    // Neutral 负责稳定 Astryx 组件基线，产品色由全局语义 Token 统一覆盖。
    <Theme mode="light" theme={neutralTheme}>
      {children}
    </Theme>
  )
}
