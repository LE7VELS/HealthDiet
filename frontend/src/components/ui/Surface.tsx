import { Card, type CardProps } from '@astryxdesign/core/Card'

type SurfaceProps = CardProps

export function Surface({ children, className = '', ...props }: SurfaceProps) {
  return (
    <Card className={`ui-surface ${className}`.trim()} padding={0} {...props}>
      {children}
    </Card>
  )
}
