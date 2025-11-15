<template>
  <el-config-provider namespace="ep">
    <!-- 登录页面和门户页面不显示头部和侧边栏 -->
    <template v-if="$route.path === '/login' || $route.path.startsWith('/portal')">
      <RouterView />
    </template>
    <!-- 其他页面显示完整布局 -->
    <template v-else>
      <BaseHeader />
      <div class="main-container flex">
        <BaseSide />
        <div class="content-area" w="full">
          <BaseTabs />
          <div class="main-content" py="4">
            <RouterView />
          </div>
        </div>
      </div>
    </template>
  </el-config-provider>
</template>

<style>
#app {
  text-align: center;
  color: var(--ep-text-color-primary);
}

.main-container {
  height: calc(100vh - var(--ep-menu-item-height) - 4px);
}

.content-area {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.main-content {
  flex: 1;
  overflow: auto;
}
</style>
