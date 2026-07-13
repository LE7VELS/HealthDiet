import { Theme } from '@astryxdesign/core/theme'
import { neutralTheme } from '@astryxdesign/theme-neutral/built'
import type { PropsWithChildren } from 'react'

export function AppThemeProvider({ children }: PropsWithChildren) {
  return (
    <Theme mode="light" theme={neutralTheme}>
      {children}
    </Theme>
  )
}
