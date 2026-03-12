import { create } from 'zustand'

export type Theme = 'light' | 'dark'

type ThemeState = {
  theme: Theme
  background: string
  foreground: string
  accent: string
  muted: string
  setTheme: (theme: Theme) => void
}

const lightTokens = {
  background: '#fafafa',
  foreground: '#171717',
  accent: '#2563eb',
  muted: '#a3a3a3',
}

const darkTokens = {
  background: '#0a0a0a',
  foreground: '#fafafa',
  accent: '#3b82f6',
  muted: '#737373',
}

export const useThemeStore = create<ThemeState>((set) => ({
  theme: 'dark',
  ...darkTokens,
  setTheme: (theme) =>
    set({
      theme,
      ...(theme === 'light' ? lightTokens : darkTokens),
    }),
}))

export function useTheme() {
  const { background, foreground, accent, muted } = useThemeStore()
  return { background, foreground, accent, muted }
}
