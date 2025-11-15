import type { ThemeConfig } from '~/types/portal'

// 默认主题配置
export const defaultThemeConfig: ThemeConfig = {
  primaryColor: '#409EFF',
  secondaryColor: '#67C23A',
  backgroundColor: '#ffffff',
  textColor: '#333333',
  fontFamily: 'Arial, sans-serif',
  fontSize: '14px',
  borderRadius: '4px',
  shadow: '0 2px 4px rgba(0,0,0,0.1)',
  logo: '',
  favicon: '',
  siteName: '门户网站',
  siteDescription: '欢迎访问我们的门户网站',
  footerText: '© 2024 门户网站. All rights reserved.',
  contactPhone: '400-123-4567',
  contactEmail: 'info@company.com',
  contactAddress: '北京市朝阳区xxx街道',
}

// 应用主题配置到CSS变量
export function applyThemeConfig(config: ThemeConfig) {
  const root = document.documentElement

  // 设置CSS变量
  root.style.setProperty('--portal-primary-color', config.primaryColor || '#409eff')
  root.style.setProperty('--portal-secondary-color', config.secondaryColor || '#67c23a')
  root.style.setProperty('--portal-background-color', config.backgroundColor || '#ffffff')
  root.style.setProperty('--portal-text-color', config.textColor || '#333333')
  root.style.setProperty('--portal-font-family', config.fontFamily || 'Arial, sans-serif')
  root.style.setProperty('--portal-font-size', config.fontSize || '14px')
  root.style.setProperty('--portal-border-radius', config.borderRadius || '4px')
  root.style.setProperty('--portal-shadow', config.shadow || '0 2px 4px rgba(0,0,0,0.1)')

  // 设置网站信息
  root.style.setProperty('--portal-site-name', config.siteName || '门户网站')
  root.style.setProperty('--portal-site-description', config.siteDescription || '欢迎访问我们的门户网站')
  root.style.setProperty('--portal-footer-text', config.footerText || '© 2024 门户网站. All rights reserved.')

  // 设置Logo和Favicon
  if (config.logo) {
    root.style.setProperty('--portal-logo-url', `url(${config.logo})`)
  }
  if (config.favicon) {
    const link = document.querySelector('link[rel="icon"]') || document.createElement('link')
    link.setAttribute('rel', 'icon')
    link.setAttribute('href', config.favicon)
    if (!document.querySelector('link[rel="icon"]')) {
      document.head.appendChild(link)
    }
  }
}

// 重置主题配置为默认值
export function resetThemeConfig() {
  applyThemeConfig(defaultThemeConfig)
}

// 从JSON字符串解析主题配置
export function parseThemeConfig(configString: string): ThemeConfig {
  try {
    const config = JSON.parse(configString)
    return { ...defaultThemeConfig, ...config }
  }
  catch (error) {
    console.error('解析主题配置失败:', error)
    return defaultThemeConfig
  }
}

// 验证主题配置
export function validateThemeConfig(config: Partial<ThemeConfig>): boolean {
  const requiredFields = ['primaryColor', 'secondaryColor', 'backgroundColor', 'textColor']

  for (const field of requiredFields) {
    if (!config[field as keyof ThemeConfig]) {
      return false
    }
  }

  return true
}
