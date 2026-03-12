package cli

import (
	"os"
	"path/filepath"
)

func writeAppFiles(destDir string) error {
	appDir := filepath.Join(destDir, "app")
	dirs := []string{
		appDir,
		filepath.Join(appDir, "src"),
		filepath.Join(appDir, "public"),
		filepath.Join(appDir, "src", "layout"),
		filepath.Join(appDir, "src", "pages"),
		filepath.Join(appDir, "src", "components", "ui"),
		filepath.Join(appDir, "src", "stores"),
		filepath.Join(appDir, "src", "lib"),
	}
	for _, d := range dirs {
		if err := mkdirAll(d); err != nil {
			return err
		}
	}
	files := map[string]string{
		"package.json":     appPackageJSON,
		"index.html":       appIndexHTML,
		"vite.config.ts":   appViteConfig,
		"tailwind.config.js": appTailwindConfig,
		"postcss.config.js":  appPostcssConfig,
		"tsconfig.json":     appTsconfig,
		"src/main.tsx":     appMainTsx,
		"src/App.tsx":      appAppTsx,
		"src/index.css":    appIndexCSS,
		"src/vite-env.d.ts": appViteEnvDts,
		"src/types.ts":     appTypesTs,
		"src/lib/http.ts":  appLibHttpTs,
		"src/stores/themeStore.ts": appThemeStoreTs,
		"src/layout/Layout.tsx":   appLayoutTsx,
		"src/layout/Header.tsx":  appHeaderTsx,
		"src/pages/Home.tsx":     appHomeTsx,
		"src/components/ui/Text.tsx":  appTextTsx,
		"src/components/ui/View.tsx":  appViewTsx,
	}
	for name, content := range files {
		p := filepath.Join(appDir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(p, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

func mkdirAll(path string) error {
	return os.MkdirAll(path, 0755)
}

var (
	appPackageJSON = `{"name":"monolith-app","private":true,"version":"0.0.0","type":"module","scripts":{"dev":"vite","build":"tsc -b && vite build","preview":"vite preview"},"dependencies":{"react":"^18.3.1","react-dom":"^18.3.1","zustand":"^5.0.2"},"devDependencies":{"@types/react":"^18.3.12","@types/react-dom":"^18.3.1","@vitejs/plugin-react":"^4.3.4","autoprefixer":"^10.4.20","postcss":"^8.4.49","tailwindcss":"^3.4.15","typescript":"~5.6.2","vite":"^6.0.1"}}
`
	appIndexHTML = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Monolith</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
`
	appViteConfig = `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { fileURLToPath } from 'url'

export default defineConfig({
  plugins: [react()],
  resolve: { alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) } },
  server: { port: 5173, proxy: { '/api': 'http://localhost:8080' } },
})
`
	appTailwindConfig = `/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: { extend: { colors: { background: 'var(--color-background)', foreground: 'var(--color-foreground)' } } },
  plugins: [],
}
`
	appPostcssConfig = `export default { plugins: { tailwindcss: {}, autoprefixer: {} } }
`
	appTsconfig = `{"compilerOptions":{"target":"ES2020","useDefineForClassFields":true,"lib":["ES2020","DOM","DOM.Iterable"],"module":"ESNext","skipLibCheck":true,"moduleResolution":"bundler","allowImportingTsExtensions":true,"resolveJsonModule":true,"isolatedModules":true,"noEmit":true,"jsx":"react-jsx","strict":true,"noUnusedLocals":true,"noUnusedParameters":true,"noFallthroughCasesInSwitch":true,"baseUrl":".","paths":{"@/*":["src/*"]}},"include":["src"]}
`
	appMainTsx = `import React from 'react'
import ReactDOM from 'react-dom/client'
import { App } from './App'
import './index.css'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode><App /></React.StrictMode>
)
`
	appAppTsx = `import { useEffect } from 'react'
import { Layout } from '@/layout/Layout'
import { Home } from '@/pages/Home'
import { useThemeStore } from '@/stores/themeStore'

export function App() {
  const { background, foreground } = useThemeStore()
  useEffect(() => {
    document.documentElement.style.setProperty('--color-bg', background)
    document.documentElement.style.setProperty('--color-fg', foreground)
  }, [background, foreground])
  return (<Layout><Home /></Layout>)
}
`
	appIndexCSS = `@tailwind base;@tailwind components;@tailwind utilities;
:root{--color-background:#0a0a0a;--color-foreground:#fafafa;--color-bg:var(--color-background);--color-fg:var(--color-foreground)}
body{margin:0;font-family:system-ui,sans-serif}
`
	appViteEnvDts = `/// <reference types="vite/client" />
interface ImportMetaEnv { readonly VITE_API_URL?: string }
interface ImportMeta { readonly env: ImportMetaEnv }
`
	appTypesTs = `export interface HealthResponse { status: string }
`
	appLibHttpTs = `const base = import.meta.env.VITE_API_URL ?? '/api'
export async function get<T>(path: string): Promise<T> {
  const res = await fetch(base + path)
  if (!res.ok) throw new Error((await res.json().catch(() => ({}))).error ?? 'Request failed')
  return res.headers.get('content-type')?.includes('application/json') ? res.json() : res.text()
}
`
	appThemeStoreTs = `import { create } from 'zustand'
export type Theme = 'light'|'dark'
const dark = { background:'#0a0a0a',foreground:'#fafafa',accent:'#3b82f6',muted:'#737373' }
const light = { background:'#fafafa',foreground:'#171717',accent:'#2563eb',muted:'#a3a3a3' }
export const useThemeStore = create((set)=>({ theme:'dark',...dark, setTheme:(t)=>set({ theme:t,...(t==='light'?light:dark) }) }))
export function useTheme() { return useThemeStore(s=>({ background:s.background,foreground:s.foreground,accent:s.accent,muted:s.muted })) }
`
	appLayoutTsx = `import { type ReactNode } from 'react'
import { View } from '@/components/ui/View'
import { Header } from '@/layout/Header'
export function Layout({ children }: { children: ReactNode }) {
  return (<View className="min-h-screen flex flex-col" style={{ backgroundColor:'var(--color-bg)', color:'var(--color-fg)' }}><Header /><View as="main" className="flex-1 p-6">{children}</View></View>)
}
`
	appHeaderTsx = `import { View } from '@/components/ui/View'
import { Text } from '@/components/ui/Text'
import { useThemeStore } from '@/stores/themeStore'
export function Header() {
  const { theme, setTheme, foreground } = useThemeStore()
  return (<View as="header" className="border-b px-6 py-4 flex items-center justify-between"><Text as="h1" className="text-xl font-semibold">Monolith</Text><button type="button" onClick={()=>setTheme(theme==='dark'?'light':'dark')} className="px-3 py-1 rounded text-sm" style={{ color:foreground, border:'1px solid '+foreground }}>{theme==='dark'?'light':'dark'}</button></View>)
}
`
	appHomeTsx = `import { View } from '@/components/ui/View'
import { Text } from '@/components/ui/Text'
import { useTheme } from '@/stores/themeStore'
export function Home() {
  const { background, foreground } = useTheme()
  return (<View className="rounded-lg p-8 max-w-xl" style={{ backgroundColor:background, color:foreground, border:'1px solid '+foreground }}><Text as="h2" className="text-2xl font-semibold mb-4">Welcome</Text><Text className="opacity-90">Monolith Go — health at /api/health.</Text></View>)
}
`
	appTextTsx = `import { type CSSProperties, type ReactNode } from 'react'
export function Text({ children, className='', style, as: Tag='p' }: { children: ReactNode; className?: string; style?: CSSProperties; as?: 'p'|'span'|'h1'|'h2'|'h3'|'label' }) { return <Tag className={className} style={style}>{children}</Tag> }
`
	appViewTsx = `import { type CSSProperties, type ReactNode } from 'react'
export function View({ children, className='', style, as: Tag='div' }: { children: ReactNode; className?: string; style?: CSSProperties; as?: 'div'|'section'|'main'|'header'|'footer'|'article' }) { return <Tag className={className} style={style}>{children}</Tag> }
`
)
