import type { Res } from '~/utils/request'
import { defineStore } from 'pinia'
import { computed, ref, type Ref } from 'vue'
import { useRouter } from 'vue-router'
import request from '~/utils/request'

// 定义菜单项类型
interface MenuItem {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  parentId: string
  path: string
  name: string
  hidden: boolean
  sort: number
  meta: {
    title: string
    icon: string
  }
  authoritys: any
  menuId: string
  children: MenuItem[] | null
}

// 定义标签页类型
interface TabItem {
  name: string
  path: string
  title: string
  icon?: string
  closable?: boolean
}

const user = ref({
  ID: 0,
  userName: '',
  nickName: '',
  headerImg: '',
  authorityId: 0,
  authority: {
    authorityId: 0,
    authorityName: '',
    parentId: 0,
    defaultRouter: '',
  },
  enable: 0,
  menus: [] as MenuItem[],
  site: {
    title: '',
    logo: '',
  },
})

export const useUserStore = defineStore('user', () => {
  const router = useRouter()

  // 标签页相关状态 - 从localStorage初始化
  const getInitialTabs = (): TabItem[] => {
    try {
      const savedTabs = localStorage.getItem('user-tabs')
      return savedTabs ? JSON.parse(savedTabs) : []
    }
    catch (error) {
      console.error('Failed to parse saved tabs:', error)
      return []
    }
  }

  const getInitialActiveTab = (): string => {
    try {
      return localStorage.getItem('user-active-tab') || ''
    }
    catch (error) {
      console.error('Failed to get active tab:', error)
      return ''
    }
  }

  const tabs = ref<TabItem[]>(getInitialTabs())
  const activeTab = ref(getInitialActiveTab())

  // 保存标签页到localStorage
  const saveTabsToStorage = () => {
    try {
      localStorage.setItem('user-tabs', JSON.stringify(tabs.value))
      localStorage.setItem('user-active-tab', activeTab.value)
    }
    catch (error) {
      console.error('Failed to save tabs to storage:', error)
    }
  }

  // 监听标签页变化，自动保存
  const saveTabs = () => {
    saveTabsToStorage()
  }

  const site = computed(() => {
    return user.value.site
  })

  const token = ref('')
  if (localStorage.getItem('x-token')) {
    token.value = localStorage.getItem('x-token') as string
  }

  // 添加标签页
  const addTab = (menuItem: any) => {
    const existingTab = tabs.value.find(tab => tab.path === menuItem.path)
    if (!existingTab) {
      const newTab: TabItem = {
        name: menuItem.name,
        path: menuItem.path,
        title: menuItem.label,
        icon: menuItem.icon,
        closable: true,
      }
      tabs.value.push(newTab)
    }
    activeTab.value = menuItem.path
    saveTabs() // 保存到localStorage
  }

  // 手动添加标签页
  const addManualTab = (tabItem: TabItem) => {
    const existingTab = tabs.value.find(tab => tab.path === tabItem.path)
    if (!existingTab) {
      tabs.value.push(tabItem)
    }
    activeTab.value = tabItem.path
    saveTabs() // 保存到localStorage
  }

  // 关闭标签页
  const closeTab = (targetPath: string) => {
    const tabsList = tabs.value
    let activeIndex = tabsList.findIndex((tab: TabItem) => tab.path === activeTab.value)

    for (let i = 0; i < tabsList.length; i++) {
      if (tabsList[i].path === targetPath) {
        tabsList.splice(i, 1)
        break
      }
    }

    // 如果关闭的是当前激活的标签页，需要切换到其他标签页
    if (targetPath === activeTab.value) {
      if (tabsList.length && activeIndex >= tabsList.length) {
        activeIndex = tabsList.length - 1
      }
      if (tabsList.length) {
        const nextPath = tabsList[activeIndex].path
        activeTab.value = nextPath
        // 确保路径以 / 开头
        const normalizedPath = nextPath.startsWith('/') ? nextPath : `/${nextPath}`
        router.push(normalizedPath)
      }
      else {
        activeTab.value = ''
        router.push('/')
      }
    }
    saveTabs() // 保存到localStorage
  }

  // 切换标签页
  const switchTab = (path: string) => {
    // 检查 path 是否有效
    if (!path) {
      console.error('switchTab: path 参数无效:', path)
      return
    }

    activeTab.value = path
    // 确保路径以 / 开头
    const normalizedPath = path.startsWith('/') ? path : `/${path}`

    // 调试信息
    if (import.meta.env.DEV) {
      console.warn('切换标签页:', {
        originalPath: path,
        normalizedPath,
        currentRoute: router.currentRoute.value.path,
      })
    }

    router.push(normalizedPath)
    saveTabs() // 保存到localStorage
  }

  // 关闭其他标签页
  const closeOtherTabs = (path: string) => {
    tabs.value = tabs.value.filter(tab => tab.path === path || !tab.closable)
    activeTab.value = path
    saveTabs() // 保存到localStorage
  }

  // 关闭所有标签页
  const closeAllTabs = () => {
    tabs.value = tabs.value.filter(tab => !tab.closable)
    activeTab.value = ''
    router.push('/')
    saveTabs() // 保存到localStorage
  }

  // 清除标签页存储
  const clearTabsStorage = () => {
    localStorage.removeItem('user-tabs')
    localStorage.removeItem('user-active-tab')
  }

  // 初始化标签页状态
  const initTabs = () => {
    // 从localStorage恢复标签页
    const savedTabs = getInitialTabs()
    const savedActiveTab = getInitialActiveTab()

    if (savedTabs.length > 0) {
      tabs.value = savedTabs
      if (savedActiveTab && savedTabs.find(tab => tab.path === savedActiveTab)) {
        activeTab.value = savedActiveTab
        // 确保路由也同步到当前激活的标签页
        router.push(savedActiveTab)
      }
    }
  }

  // 转换菜单数据格式，适配前端组件
  const transformMenuData = (menus: MenuItem[]): any[] => {
    return menus.map((menu) => {
      // 确保路径以 / 开头
      const finalPath = menu.path.startsWith('/') ? menu.path : `/${menu.path}`

      return {
        label: menu.meta.title,
        key: menu.name,
        icon: menu.meta.icon,
        name: menu.name,
        path: finalPath,
        hidden: menu.hidden,
        children: menu.children ? transformMenuData(menu.children) : undefined,
      }
    })
  }

  // {label, key, icon, children, actions}
  const menu = computed(() => {
    return transformMenuData(user.value.menus || [])
  })

  // console.log(menu)
  // 异步拉取用户信息
  const checkUserInfo = async (mute: boolean = true) => {
    if (!token.value) {
      return false
    }
    return new Promise((resolve) => {
      request
        .get(
          '/user/getUserInfo',
          {},
          {
            headers: {
              mute,
            },
          },
        )
        .then(async (res: any) => {
          user.value = res.data.userInfo
          user.value.site = {
            title: '',
            logo: '',
          }
          user.value.menus = res.data.menus

          resolve(true)
        })
        .catch((_err) => {
          resolve(false)
        })
    })
  }

  // 登录
  const login = async (formData: any, loading: Ref<boolean>) => {
    loading.value = true
    const res = (await request.post('/base/login', formData, {
      headers: {
        useToken: false,
      },
    })) as Res
    console.warn('Login response:', res)
    if (res.code !== 0) {
      loading.value = false
      return
    }
    token.value = res.data.token
    user.value = res.data.user
    localStorage.setItem('x-token', token.value)
    if (token.value) {
      router.push('/')
    }
  }

  // 重置用户信息
  const resetUserInfo = () => {
    user.value = {
      ID: 0,
      userName: '',
      nickName: '',
      headerImg: '',
      authorityId: 0,
      authority: {
        authorityId: 0,
        authorityName: '',
        parentId: 0,
        defaultRouter: '',
      },
      enable: 0,
      menus: [],
      site: {
        title: '',
        logo: '',
      },
    }
    token.value = ''
    // 清除标签页存储
    clearTabsStorage()
    tabs.value = []
    activeTab.value = ''
  }

  // topbar - 由于用户数据结构中没有topbar，返回空数组
  const topbar = computed(() => {
    return []
  })

  return {
    user,
    site,
    menu,
    topbar,
    login,
    checkUserInfo,
    resetUserInfo,
    tabs,
    activeTab,
    addTab,
    addManualTab,
    closeTab,
    switchTab,
    closeOtherTabs,
    closeAllTabs,
    clearTabsStorage,
    initTabs,
  }
})
