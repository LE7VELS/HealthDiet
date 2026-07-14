import { Banner as AstryxBanner, type BannerProps } from '@astryxdesign/core/Banner'
import {
  ToastViewport,
  useToast,
  type ShowToastFn,
} from '@astryxdesign/core/Toast'

export function Banner(props: BannerProps) {
  return <AstryxBanner {...props} />
}

export function AppToastViewport({ children }: { children: React.ReactNode }) {
  return (
    <ToastViewport inset={{ end: 8, top: 8 }} maxVisible={3} position="topEnd">
      {children}
    </ToastViewport>
  )
}

export function useAppToast(): ShowToastFn {
  // 使用 Astryx 原生 Toast 内容结构，避免覆盖其主题和无障碍行为。
  return useToast()
}
