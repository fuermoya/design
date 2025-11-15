import type { UserModule } from '~/types'
import ElementPlus from 'element-plus'
// @ts-expect-error - Element Plus locale module type issue
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'

export const install: UserModule = ({ app }) => {
  // 配置 Element Plus 使用中文
  app.use(ElementPlus, {
    locale: zhCn,
  })
}
