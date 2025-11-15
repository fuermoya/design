import type { UserModule } from '~/types'
import { createPinia } from 'pinia'
import { useUserStore } from '~/stores/user'

export const install: UserModule = ({ router, app }) => {
  // 初始化 Pinia
  const pinia = createPinia()
  app.use(pinia)

  // 添加一个标志位来跟踪是否已经初始化过用户信息
  let isUserInfoInitialized = false

  // 监听 localStorage 变化，当 token 被删除时重置初始化状态
  const originalRemoveItem = localStorage.removeItem
  localStorage.removeItem = function (key) {
    if (key === 'x-token') {
      isUserInfoInitialized = false
    }
    return originalRemoveItem.call(this, key)
  }

  router.beforeEach(async (to, from, next) => {
    const token = localStorage.getItem('x-token')

    // 如果是登录页面，直接放行
    if (to.path === '/login') {
      // 如果已经登录，重定向到首页
      if (token) {
        next('/')
        return
      }
      next()
      return
    }

    // 如果是前台页面，直接放行（不需要登录验证）
    if (to.path.startsWith('/portal')) {
      next()
      return
    }

    // 检查是否有token
    if (!token) {
      next('/login')
      return
    }

    // 只有在未初始化用户信息时才调用 checkUserInfo
    if (!isUserInfoInitialized) {
      const userStore = useUserStore()
      const userInfoValid = await userStore.checkUserInfo(false)

      if (!userInfoValid) {
        // 用户信息无效，清除token并跳转到登录页
        localStorage.removeItem('x-token')
        next('/login')
        return
      }

      // 标记用户信息已初始化
      isUserInfoInitialized = true
    }

    // 用户信息有效，允许访问
    next()
  })
}
