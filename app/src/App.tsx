import { useEffect } from 'react'
import { Layout } from '@/layout/Layout'
import { Home } from '@/pages/Home'
import { useThemeStore } from '@/stores/themeStore'

function applyTheme(background: string, foreground: string) {
  document.documentElement.style.setProperty('--color-bg', background)
  document.documentElement.style.setProperty('--color-fg', foreground)
}

export function App() {
  const { background, foreground } = useThemeStore()
  useEffect(() => {
    applyTheme(background, foreground)
  }, [background, foreground])
  return (
    <Layout>
      <Home />
    </Layout>
  )
}
