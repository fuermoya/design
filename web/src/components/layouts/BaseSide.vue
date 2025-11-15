<script lang="ts" setup>
// 导入用户存储
import { useUserStore } from '~/stores/user'

// 获取用户存储实例
const userStore = useUserStore()
// console.log(userStore.menu)
// 菜单打开/关闭处理函数
function handleOpen(_key: string, _keyPath: string[]) {
  // 菜单打开处理逻辑
  // 可以在这里添加菜单打开时的逻辑
}

function handleClose(_key: string, _keyPath: string[]) {
  // 菜单关闭处理逻辑
  // 可以在这里添加菜单关闭时的逻辑
}

// 菜单点击处理函数
function handleMenuClick(index: string) {
  // 查找对应的菜单项
  const findMenuItem = (menus: any[], path: string): any => {
    // console.log(menus, path)
    for (const menu of menus) {
      if (menu.path === path) {
        return menu
      }
      if (menu.children) {
        const found = findMenuItem(menu.children, path)
        if (found)
          return found
      }
    }
    return null
  }

  const menuItem = findMenuItem(userStore.menu, index)
  if (menuItem) {
    // 添加标签页并切换（这会自动处理路由导航）
    userStore.addTab(menuItem)
    userStore.switchTab(menuItem.path)
  }
  else {
    console.warn('未找到菜单项:', index)
  }
}
</script>

<template>
  <el-menu
    :default-active="$route.path"
    class="el-menu-vertical-demo"
    @open="handleOpen"
    @close="handleClose"
    @select="handleMenuClick"
  >
    <!-- 根据用户存储中的菜单数据动态渲染 -->
    <template v-for="item in userStore.menu" :key="item.key">
      <!-- 有子菜单的情况 -->
      <el-sub-menu v-if="item.children && item.children.length > 0 && !item.hidden" :index="item.path">
        <template #title>
          <el-icon v-if="item.icon">
            <component :is="item.icon" />
          </el-icon>
          <span>{{ item.label }}</span>
        </template>
        <!-- 渲染子菜单 -->
        <template v-for="child in item.children" :key="child.key">
          <el-sub-menu v-if="child.children && child.children.length > 0 && !child.hidden" :index="child.path">
            <template #title>
              <el-icon v-if="child.icon">
                <component :is="child.icon" />
              </el-icon>
              <span>{{ child.label }}</span>
            </template>
            <el-menu-item
              v-for="grandChild in child.children"
              v-show="!grandChild.hidden"
              :key="grandChild.key"
              :index="grandChild.path"
            >
              <el-icon v-if="grandChild.icon">
                <component :is="grandChild.icon" />
              </el-icon>
              <span>{{ grandChild.label }}</span>
            </el-menu-item>
          </el-sub-menu>
          <!-- 普通子菜单项 -->
          <el-menu-item
            v-else-if="!child.hidden"
            :key="child.key"
            :index="child.path"
          >
            <el-icon v-if="child.icon">
              <component :is="child.icon" />
            </el-icon>
            <span>{{ child.label }}</span>
          </el-menu-item>
        </template>
      </el-sub-menu>
      <!-- 普通菜单项 -->
      <el-menu-item
        v-else-if="!item.hidden"
        :key="item.key"
        :index="item.path"
      >
        <el-icon v-if="item.icon">
          <component :is="item.icon" />
        </el-icon>
        <span>{{ item.label }}</span>
      </el-menu-item>
    </template>
  </el-menu>
</template>
