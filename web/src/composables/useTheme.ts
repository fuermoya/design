import type { ThemeConfig } from '~/types/portal'
import { ref, watch } from 'vue'
import request from '~/utils/request'
import { applyThemeConfig, defaultThemeConfig, parseThemeConfig } from '~/utils/theme'

// 主题状态管理
const currentTheme = ref<ThemeConfig | null>(null)
const isLoading = ref(false)

// 加载当前激活的主题
export async function loadActiveTheme() {
  isLoading.value = true
  try {
    const res: any = await request.get('/portal/theme')
    if (res.code === 0 && res.data.theme) {
      const theme = res.data.theme
      if (theme.config) {
        currentTheme.value = parseThemeConfig(theme.config)
        applyThemeConfig(currentTheme.value)
      }
    }
  }
  catch (error) {
    console.error('加载当前主题失败:', error)
    // 使用默认主题
    currentTheme.value = defaultThemeConfig
    applyThemeConfig(defaultThemeConfig)
  }
  finally {
    isLoading.value = false
  }
}

// 预览主题配置
export function previewTheme(config: ThemeConfig) {
  applyThemeConfig(config)
}

// 重置主题配置
export function resetTheme() {
  currentTheme.value = defaultThemeConfig
  applyThemeConfig(defaultThemeConfig)
}

// 监听主题变化
watch(currentTheme, (newTheme) => {
  if (newTheme) {
    applyThemeConfig(newTheme)
  }
}, { deep: true })

// 导出主题状态
export function useTheme() {
  return {
    currentTheme,
    isLoading,
    loadActiveTheme,
    previewTheme,
    resetTheme,
  }
}
