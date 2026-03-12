import { type ReactNode } from 'react'
import { View } from '@/components/ui/View'
import { Header } from '@/layout/Header'

type LayoutProps = { children: ReactNode }

export function Layout({ children }: LayoutProps) {
  return (
    <View
      className="min-h-screen flex flex-col"
      style={{ backgroundColor: 'var(--color-bg)', color: 'var(--color-fg)' }}
    >
      <Header />
      <View as="main" className="flex-1 p-6">
        {children}
      </View>
    </View>
  )
}
