import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useState, type PropsWithChildren } from 'react'
import { AppThemeProvider, AppToastViewport } from '../components/ui'

export function AppProviders({ children }: PropsWithChildren) {
  const [queryClient] = useState(
    () =>
      new QueryClient({
        defaultOptions: {
          queries: {
            retry: 1,
            staleTime: 30_000,
          },
        },
      }),
  )

  return (
    <QueryClientProvider client={queryClient}>
      <AppThemeProvider>
        <AppToastViewport>{children}</AppToastViewport>
      </AppThemeProvider>
    </QueryClientProvider>
  )
}
