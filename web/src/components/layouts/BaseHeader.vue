<script setup lang="ts">
import { ElMessage } from 'element-plus'

import { useRouter } from 'vue-router'
import { toggleDark } from '~/composables'
import { useUserStore } from '~/stores/user'

import request from '~/utils/request'

const router = useRouter()
const userStore = useUserStore()

function goToHomePage() {
  router.push('/')
}

function logout() {
  // 调用退出登录接口
  request.post('/jwt/jsonInBlacklist')
  // 重置用户信息
  userStore.resetUserInfo()
  localStorage.removeItem('x-token')
  ElMessage.success('已退出登录')
  router.push('/login')
}

function goToPersonInfo() {
  router.push('/person')
}

function goToPortal() {
  // 在新标签页中打开门户网站，通过代理访问
  // 使用完整的URL确保不会影响当前页面的路由
  const currentOrigin = window.location.origin
  window.open(`${currentOrigin}/portal`, '_blank')
}

// 处理下拉菜单命令
function handleCommand(command: string) {
  switch (command) {
    case 'person':
      goToPersonInfo()
      break
    case 'logout':
      logout()
      break
  }
}
</script>

<template>
  <el-menu class="el-menu-demo" mode="horizontal" :ellipsis="false">
    <el-menu-item index="goToHomePage" @click="goToHomePage">
      <div class="flex items-center justify-center gap-2">
        <div class="text-xl" i-ep-element-plus />
        <span>首页</span>
      </div>
    </el-menu-item>

    <div class="flex-1" />
    <!-- 门户网站跳转 -->
    <el-menu-item index="goToPortal" @click="goToPortal">
      <div class="flex items-center justify-center gap-2">
        <i class="el-icon-house" />
        <span>门户网站</span>
      </div>
    </el-menu-item>
    <el-menu-item h="full" index="toggleDark" @click="toggleDark()">
      <button class="w-full cursor-pointer border-none bg-transparent" style="height: var(--ep-menu-item-height)">
        <i inline-flex i="dark:ep-moon ep-sunny" />
      </button>
    </el-menu-item>
    <!-- 用户信息下拉菜单 -->
    <el-dropdown trigger="click" @command="handleCommand">
      <!-- <el-menu-item class="user-menu-item"> -->
      <div class="flex items-center gap-2">
        <el-avatar
          :size="32" :src="userStore.user.headerImg"
          :icon="userStore.user.headerImg ? undefined : 'el-icon-user'"
        />
        <span class="user-name">{{ userStore.user.nickName || userStore.user.userName }}</span>
        <i class="el-icon-arrow-down" />
      </div>
      <!-- </el-menu-item> -->

      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="person">
            <i class="el-icon-user" style="margin-right: 8px;" />
            个人信息
          </el-dropdown-item>
          <el-dropdown-item divided command="logout">
            <i class="el-icon-switch-button" style="margin-right: 8px;" />
            退出登录
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </el-menu>
</template>

<style lang="scss">
.el-menu-demo {
  &.ep-menu--horizontal > .ep-menu-item:nth-child(1) {
    margin-right: auto;
  }

  .user-menu-item {
    padding: 0 16px;

    .user-name {
      max-width: 120px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}

// 下拉菜单样式
.el-dropdown-menu {
  .el-dropdown-item {
    display: flex;
    align-items: center;
    padding: 8px 16px;

    i {
      font-size: 16px;
    }
  }
}
</style>
