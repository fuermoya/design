<script setup lang="ts">
import type { SysCategory } from '~/types/portal'
import { onMounted, ref } from 'vue'
import { loadActiveTheme, useTheme } from '~/composables/useTheme'
import request from '~/utils/request'

// 响应式数据
const categories = ref<SysCategory[]>([])
const { currentTheme } = useTheme()

// 加载分类列表
async function loadCategories() {
  try {
    const res: any = await request.get('/portal/categories')
    if (res.code === 0) {
      categories.value = res.data.list
    }
  }
  catch (error) {
    console.error('加载分类列表失败:', error)
  }
}

// 页面初始化
onMounted(() => {
  loadCategories()
  loadActiveTheme()
})
</script>

<template>
  <div class="portal-layout">
    <!-- 头部导航 -->
    <header class="portal-header">
      <div class="header-container">
        <div class="logo-section">
          <img v-if="currentTheme?.logo" :src="currentTheme.logo" alt="Logo" class="logo">
          <h1 class="site-name">
            {{ currentTheme?.siteName || '企业门户' }}
          </h1>
        </div>
        <nav class="main-nav">
          <router-link to="/portal" class="nav-item" active-class="active">
            首页
          </router-link>
          <router-link to="/portal/about" class="nav-item" active-class="active">
            关于我们
          </router-link>
          <router-link
            v-for="category in categories"
            :key="category.ID"
            :to="`/portal/category/${category.ID}`"
            class="nav-item"
            active-class="active"
          >
            {{ category.name }}
          </router-link>
          <router-link to="/portal/contact" class="nav-item" active-class="active">
            联系我们
          </router-link>
        </nav>
      </div>
    </header>

    <!-- 主要内容区域 -->
    <main class="portal-main">
      <slot />
    </main>

    <!-- 页脚 -->
    <footer class="portal-footer">
      <div class="footer-container">
        <div class="footer-content">
          <div class="footer-section">
            <h4>关于我们</h4>
            <p>{{ currentTheme?.siteDescription || '专业的服务，卓越的品质' }}</p>
          </div>
          <div class="footer-section">
            <h4>联系方式</h4>
            <p>电话：{{ currentTheme?.contactPhone || '400-123-4567' }}</p>
            <p>邮箱：{{ currentTheme?.contactEmail || 'info@company.com' }}</p>
            <p>地址：{{ currentTheme?.contactAddress || '北京市朝阳区xxx街道' }}</p>
          </div>
          <div class="footer-section">
            <h4>快速链接</h4>
            <ul>
              <li>
                <router-link to="/portal">
                  首页
                </router-link>
              </li>
              <li>
                <router-link to="/portal/about">
                  关于我们
                </router-link>
              </li>
              <li>
                <router-link to="/portal/contact">
                  联系我们
                </router-link>
              </li>
            </ul>
          </div>
        </div>
        <div class="footer-bottom">
          <p>{{ currentTheme?.footerText || '© 2024 企业门户. All rights reserved.' }}</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.portal-layout {
  min-height: 100vh;
  background-color: var(--portal-background-color, #f5f7fa);
  display: flex;
  flex-direction: column;
  font-family: var(--portal-font-family, Arial, sans-serif);
  font-size: var(--portal-font-size, 14px);
  color: var(--portal-text-color, #333333);
}

/* 头部样式 */
.portal-header {
  background: linear-gradient(
    135deg,
    var(--portal-primary-color, #409eff) 0%,
    var(--portal-secondary-color, #67c23a) 100%
  );
  color: white;
  padding: 1rem 0;
  position: sticky;
  top: 0;
  z-index: 1000;
  box-shadow: var(--portal-shadow, 0 2px 8px rgba(0, 0, 0, 0.1));
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo-section {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.logo {
  height: 40px;
  width: auto;
  border-radius: var(--portal-border-radius, 4px);
}

.site-name {
  font-size: 1.5rem;
  font-weight: bold;
  margin: 0;
}

.main-nav {
  display: flex;
  gap: 2rem;
}

.nav-item {
  color: white;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: var(--portal-border-radius, 4px);
  transition: all 0.3s ease;
}

.nav-item:hover,
.nav-item.active {
  background-color: rgba(255, 255, 255, 0.2);
}

/* 主要内容区域 */
.portal-main {
  flex: 1;
  padding: 2rem 0;
}

/* 页脚样式 */
.portal-footer {
  background: linear-gradient(
    135deg,
    var(--portal-primary-color, #409eff) 0%,
    var(--portal-secondary-color, #67c23a) 100%
  );
  color: white;
  padding: 2rem 0 1rem;
  margin-top: auto;
}

.footer-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

.footer-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.footer-section h4 {
  margin: 0 0 1rem;
  font-size: 1.1rem;
  font-weight: bold;
}

.footer-section p {
  margin: 0.5rem 0;
  line-height: 1.6;
}

.footer-section ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.footer-section ul li {
  margin: 0.5rem 0;
}

.footer-section ul li a {
  color: white;
  text-decoration: none;
  transition: opacity 0.3s ease;
}

.footer-section ul li a:hover {
  opacity: 0.8;
}

.footer-bottom {
  border-top: 1px solid rgba(255, 255, 255, 0.2);
  padding-top: 1rem;
  text-align: center;
}

.footer-bottom p {
  margin: 0;
  opacity: 0.8;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header-container {
    flex-direction: column;
    gap: 1rem;
  }

  .main-nav {
    flex-wrap: wrap;
    justify-content: center;
  }

  .footer-content {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
}
</style>
