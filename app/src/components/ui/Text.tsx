import { type CSSProperties, type ReactNode } from 'react'

type TextProps = {
  children: ReactNode
  className?: string
  style?: CSSProperties
  as?: 'p' | 'span' | 'h1' | 'h2' | 'h3' | 'label'
}

export function Text({
  children,
  className = '',
  style,
  as: Tag = 'p',
}: TextProps) {
  return (
    <Tag className={className} style={style}>
      {children}
    </Tag>
  )
}
