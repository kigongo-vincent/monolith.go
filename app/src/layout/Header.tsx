import { View } from '@/components/ui/View'
import { Text } from '@/components/ui/Text'
import { useThemeStore } from '@/stores/themeStore'

export function Header() {
  const { theme, setTheme, foreground } = useThemeStore()
  const next = theme === 'dark' ? 'light' : 'dark'
  return (
    <View as="header" className="border-b px-6 py-4 flex items-center justify-between">
      <Text as="h1" className="text-xl font-semibold">
        Monolith
      </Text>
      <button
        type="button"
        onClick={() => setTheme(next)}
        className="px-3 py-1 rounded text-sm"
        style={{ color: foreground, border: `1px solid ${foreground}` }}
      >
        {next}
      </button>
    </View>
  )
}
