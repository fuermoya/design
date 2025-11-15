<script lang="ts" setup>
// 超级管理员页面组件
</script>

<template>
  <div class="super-admin-container">
    <!-- 后台管理标识 -->
    <div class="admin-badge">
      <span>后台管理</span>
    </div>

    <h1>系统管理</h1>
    <p>系统管理主页面</p>
  </div>
</template>

<style scoped>
.super-admin-container {
  padding: 20px;
  position: relative;
}

.admin-badge {
  position: fixed;
  top: 20px;
  right: 20px;
  background: linear-gradient(135deg, #409eff, #67c23a);
  color: white;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  z-index: 1000;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  animation: adminBadgePulse 2s ease-in-out infinite;
}

@keyframes adminBadgePulse {
  0%,
  100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}
</style>
