import { View } from '@/components/ui/View'
import { Text } from '@/components/ui/Text'
import { useTheme } from '@/stores/themeStore'

export function Home() {
  const { background, foreground } = useTheme()
  return (
    <View
      className="rounded-lg p-8 max-w-xl"
      style={{
        backgroundColor: background,
        color: foreground,
        border: `1px solid ${foreground}`,
      }}
    >
      <Text as="h2" className="text-2xl font-semibold mb-4">
        Welcome
      </Text>
      <Text className="opacity-90">
        Monolith Go — backend and UI in one place. Health check at /api/health.
      </Text>
    </View>
  )
}
