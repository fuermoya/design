<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useUserStore } from '~/stores/user'

const userStore = useUserStore()

// 组件挂载时初始化标签页
onMounted(() => {
  userStore.initTabs()
})

// 监听路由变化，自动添加标签页
// watch(() => route.path, (newPath: string) => {
//   if (newPath === '/login')
//     return

//   // 查找对应的菜单项
//   const findMenuItem = (menus: any[], path: string): any => {
//     for (const menu of menus) {
//       if (menu.path === path) {
//         return menu
//       }
//       if (menu.children) {
//         const found = findMenuItem(menu.children, path)
//         if (found)
//           return found
//       }
//     }
//     return null
//   }

//   const menuItem = findMenuItem(userStore.menu, newPath.replace('/', ''))
//   if (menuItem) {
//     userStore.addTab(menuItem)
//   }
// }, { immediate: true })

// 标签页点击事件
function handleTabClick(tab: any) {
  // 从 userStore.tabs 中找到对应的标签页数据
  const tabData = userStore.tabs.find(t => t.path === tab.paneName)
  if (tabData) {
    userStore.switchTab(tabData.path)
  }
  else {
    console.error('未找到标签页数据:', tab.paneName)
  }
}

// 标签页关闭事件
function handleTabRemove(targetPath: string) {
  userStore.closeTab(targetPath)
}

// 右键菜单事件
function handleContextMenu(e: MouseEvent, _tab: any) {
  e.preventDefault()
  // console.log('右键菜单:', tab)
}
</script>

<template>
  <div class="base-tabs">
    <el-tabs
      v-model="userStore.activeTab"
      type="card"
      closable
      @tab-click="handleTabClick"
      @tab-remove="handleTabRemove"
    >
      <el-tab-pane
        v-for="tab in userStore.tabs"
        :key="tab.path"
        :label="tab.title"
        :name="tab.path"
        :closable="tab.closable"
        @contextmenu="(e) => handleContextMenu(e, tab)"
      >
        <template #label>
          <el-icon v-if="tab.icon" class="tab-icon">
            <component :is="tab.icon" />
          </el-icon>
          <span>{{ tab.title }}</span>
        </template>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style scoped>
.base-tabs {
  background: var(--ep-bg-color);
  /* border-bottom: 1px solid var(--ep-border-color-light); */
  padding: 0 0px;
}

.tab-icon {
  margin-right: 4px;
  font-size: 14px;
}

:deep(.el-tabs__header) {
  margin: 0;
}

:deep(.el-tabs__nav-wrap) {
  padding: 0;
}

:deep(.el-tabs__item) {
  height: 32px;
  line-height: 32px;
  font-size: 13px;
}

:deep(.el-tabs__content) {
  display: none;
}
</style>
