import { type CSSProperties, type ReactNode } from 'react'

type ViewProps = {
  children: ReactNode
  className?: string
  style?: CSSProperties
  as?: 'div' | 'section' | 'main' | 'header' | 'footer' | 'article'
}

export function View({
  children,
  className = '',
  style,
  as: Tag = 'div',
}: ViewProps) {
  return (
    <Tag className={className} style={style}>
      {children}
    </Tag>
  )
}
