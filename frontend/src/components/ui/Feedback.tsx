import { Banner as AstryxBanner, type BannerProps } from '@astryxdesign/core/Banner'
import { ToastViewport, useToast } from '@astryxdesign/core/Toast'

export function Banner(props: BannerProps) {
  return <AstryxBanner {...props} />
}

export function AppToastViewport({ children }: { children: React.ReactNode }) {
  return (
    <ToastViewport maxVisible={3} position="topEnd">
      {children}
    </ToastViewport>
  )
}

export function useAppToast() {
  return useToast()
}
