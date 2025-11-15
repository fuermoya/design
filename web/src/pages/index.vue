<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '~/stores/user'
import { formatDate } from '~/utils/format'
import request from '~/utils/request'

const router = useRouter()
const userStore = useUserStore()

// 文章类型定义
interface Article {
  ID: number
  title: string
  summary: string
  coverImage?: string
  viewCount: number
  likeCount: number
  CreatedAt: string
  UpdatedAt: string
}

// 统计数据
const stats = ref({
  totalUsers: 0,
  totalArticles: 0,
  totalViews: 0,
  totalLikes: 0,
  totalMessages: 0,
  totalCategories: 0,
  userChange: 0,
  articleChange: 0,
  viewChange: 0,
  likeChange: 0,
  messageChange: 0,
  categoryChange: 0,
})

// 最新文章
const latestArticles = ref<Article[]>([])

// 系统状态
const systemStatus = ref({
  cpu: 0,
  memory: 0,
  disk: 0,
  uptime: 0,
})

// 图表数据
const _chartData = ref({
  views: [],
  articles: [],
  users: [],
})

// 获取进度条颜色
function getProgressColor(percentage: number) {
  if (percentage < 50)
    return '#67c23a'
  if (percentage < 80)
    return '#e6a23c'
  return '#f56c6c'
}

// 加载统计数据
async function loadStats() {
  try {
    const res: any = await request.get('/dashboard/stats')
    if (res.code === 0) {
      stats.value = res.data
    }
  }
  catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 加载最新文章
async function loadLatestArticles() {
  try {
    const res: any = await request.get('/portal/articles', {
      page: 1,
      pageSize: 5,
    })
    if (res.code === 0) {
      latestArticles.value = res.data.list
    }
  }
  catch (error) {
    console.error('加载最新文章失败:', error)
  }
}

// 加载系统状态
async function loadSystemStatus() {
  try {
    const res: any = await request.get('/dashboard/status')
    if (res.code === 0) {
      systemStatus.value = res.data
    }
  }
  catch (error) {
    console.error('加载系统状态失败:', error)
  }
}

// 格式化时间
function formatUptime(seconds: number) {
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return `${days}天 ${hours}小时 ${minutes}分钟`
}

// 格式化百分比变化
function formatPercentageChange(change: number) {
  const sign = change >= 0 ? '+' : ''
  return `${sign}${change.toFixed(1)}%`
}

// 获取趋势样式类
function getTrendClass(change: number) {
  return change >= 0 ? 'trend-up' : 'trend-down'
}

// 查看文章详情
function viewArticle(article: Article) {
  try {
    router.push(`/portal/article/${article.ID}`)
  }
  catch (error) {
    console.error('跳转文章详情失败:', error)
  }
}

// 页面初始化
onMounted(() => {
  loadStats()
  loadLatestArticles()
  loadSystemStatus()
})

// 路由跳转函数
function navigateTo(path: string) {
  try {
    router.push(path)
    // 等待路由跳转完成后检查标签页
    setTimeout(() => {
      console.warn('当前标签页:', userStore.tabs)
      console.warn('激活的标签页:', userStore.activeTab)
    }, 100)
  }
  catch (error) {
    console.error('路由跳转失败:', error)
  }
}
</script>

<template>
  <div class="dashboard-container">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">
          <el-icon><Document /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-number">
            {{ stats.totalArticles }}
          </div>
          <div class="stat-label">
            总文章数
          </div>
        </div>
        <div class="stat-trend">
          <span :class="getTrendClass(stats.articleChange)">{{ formatPercentageChange(stats.articleChange) }}</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <el-icon><View /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-number">
            {{ stats.totalViews }}
          </div>
          <div class="stat-label">
            总浏览量
          </div>
        </div>
        <div class="stat-trend">
          <span :class="getTrendClass(stats.viewChange)">{{ formatPercentageChange(stats.viewChange) }}</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <el-icon><Star /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-number">
            {{ stats.totalLikes }}
          </div>
          <div class="stat-label">
            总点赞数
          </div>
        </div>
        <div class="stat-trend">
          <span :class="getTrendClass(stats.likeChange)">{{ formatPercentageChange(stats.likeChange) }}</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <el-icon><Message /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-number">
            {{ stats.totalMessages }}
          </div>
          <div class="stat-label">
            总留言数
          </div>
        </div>
        <div class="stat-trend">
          <span :class="getTrendClass(stats.messageChange)">{{ formatPercentageChange(stats.messageChange) }}</span>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 左侧内容 -->
      <div class="left-content">
        <!-- 系统状态 -->
        <div class="status-card">
          <div class="card-header">
            <h3 class="card-title">
              系统状态
            </h3>
            <el-tag type="success" size="small">
              运行正常
            </el-tag>
          </div>
          <div class="status-grid">
            <div class="status-item">
              <div class="status-label">
                CPU使用率
              </div>
              <div class="status-value">
                {{ systemStatus.cpu }}%
              </div>
              <el-progress :percentage="systemStatus.cpu" :color="getProgressColor(systemStatus.cpu)" />
            </div>
            <div class="status-item">
              <div class="status-label">
                内存使用率
              </div>
              <div class="status-value">
                {{ systemStatus.memory }}%
              </div>
              <el-progress :percentage="systemStatus.memory" :color="getProgressColor(systemStatus.memory)" />
            </div>
            <div class="status-item">
              <div class="status-label">
                磁盘使用率
              </div>
              <div class="status-value">
                {{ systemStatus.disk }}%
              </div>
              <el-progress :percentage="systemStatus.disk" :color="getProgressColor(systemStatus.disk)" />
            </div>
            <div class="status-item">
              <div class="status-label">
                运行时间
              </div>
              <div class="status-value">
                {{ formatUptime(systemStatus.uptime) }}
              </div>
            </div>
          </div>
        </div>

        <!-- 最新文章 -->
        <div class="articles-card">
          <div class="card-header">
            <h3 class="card-title">
              最新文章
            </h3>
            <el-button type="primary" size="small" @click="navigateTo('/portal')">
              查看更多
            </el-button>
          </div>
          <div class="articles-list">
            <div
              v-for="article in latestArticles"
              :key="article.ID"
              class="article-item"
              @click="viewArticle(article)"
            >
              <div class="article-cover">
                <img
                  v-if="article.coverImage"
                  :src="article.coverImage"
                  :alt="article.title"
                  class="cover-image"
                >
                <div v-else class="placeholder-image">
                  <el-icon><Document /></el-icon>
                </div>
              </div>
              <div class="article-info">
                <h4 class="article-title">
                  {{ article.title }}
                </h4>
                <p class="article-summary">
                  {{ article.summary }}
                </p>
                <div class="article-meta">
                  <span class="meta-item">
                    <el-icon><Calendar /></el-icon>
                    {{ formatDate(article.CreatedAt) }}
                  </span>
                  <span class="meta-item">
                    <el-icon><View /></el-icon>
                    {{ article.viewCount }}
                  </span>
                  <span class="meta-item">
                    <el-icon><Star /></el-icon>
                    {{ article.likeCount }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧内容 -->
      <div class="right-content">
        <!-- 快速操作 -->
        <div class="quick-actions-card">
          <div class="card-header">
            <h3 class="card-title">
              快速操作
            </h3>
          </div>
          <div class="actions-grid">
            <div class="action-item" @click="navigateTo('/admin/user/user')">
              <div class="action-icon">
                <el-icon><User /></el-icon>
              </div>
              <div class="action-text">
                用户管理
              </div>
            </div>
            <div class="action-item" @click="navigateTo('/substance/article/article')">
              <div class="action-icon">
                <el-icon><Document /></el-icon>
              </div>
              <div class="action-text">
                文章管理
              </div>
            </div>
            <div class="action-item" @click="navigateTo('/admin/authority/authority')">
              <div class="action-icon">
                <el-icon><Star /></el-icon>
              </div>
              <div class="action-text">
                权限管理
              </div>
            </div>
            <div class="action-item" @click="navigateTo('/substance/category/category')">
              <div class="action-icon">
                <el-icon><View /></el-icon>
              </div>
              <div class="action-text">
                分类管理
              </div>
            </div>
            <div class="action-item" @click="navigateTo('/portal')">
              <div class="action-icon">
                <el-icon><Message /></el-icon>
              </div>
              <div class="action-text">
                门户网站
              </div>
            </div>
            <div class="action-item" @click="navigateTo('/admin/operation/sysOperationRecord')">
              <div class="action-icon">
                <el-icon><Calendar /></el-icon>
              </div>
              <div class="action-text">
                操作日志
              </div>
            </div>
          </div>
        </div>

        <!-- 系统信息 -->
        <div class="system-info-card">
          <div class="card-header">
            <h3 class="card-title">
              系统信息
            </h3>
          </div>
          <div class="info-list">
            <div class="info-item">
              <span class="info-label">系统版本</span>
              <span class="info-value">v1.0.0</span>
            </div>
            <div class="info-item">
              <span class="info-label">Vue版本</span>
              <span class="info-value">3.5.13</span>
            </div>
            <div class="info-item">
              <span class="info-label">Element Plus</span>
              <span class="info-value">2.9.0</span>
            </div>
            <div class="info-item">
              <span class="info-label">Go版本</span>
              <span class="info-value">1.21.0</span>
            </div>
            <div class="info-item">
              <span class="info-label">数据库</span>
              <span class="info-value">MySQL 8.0</span>
            </div>
            <div class="info-item">
              <span class="info-label">最后更新</span>
              <span class="info-value">{{ formatDate(new Date().toISOString()) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard-container {
  padding: 2rem;
  background: #f5f7fa;
  min-height: 100vh;
}

/* 页面标题 */
.page-header {
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: bold;
  color: #303133;
  margin-bottom: 0.5rem;
}

.page-subtitle {
  color: #606266;
  font-size: 1rem;
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  font-size: 2rem;
  color: #409eff;
  background: #ecf5ff;
  padding: 1rem;
  border-radius: 12px;
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 2rem;
  font-weight: bold;
  color: #303133;
  margin-bottom: 0.25rem;
}

.stat-label {
  color: #606266;
  font-size: 0.9rem;
}

.stat-trend {
  font-size: 0.9rem;
}

.trend-up {
  color: #67c23a;
}

.trend-down {
  color: #f56c6c;
}

/* 主要内容区域 */
.main-content {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 2rem;
}

/* 卡片通用样式 */
.status-card,
.articles-card,
.quick-actions-card,
.system-info-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  margin-bottom: 1.5rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.card-title {
  font-size: 1.2rem;
  font-weight: bold;
  color: #303133;
  margin: 0;
}

/* 系统状态 */
.status-grid {
  display: grid;
  gap: 1rem;
}

.status-item {
  margin-bottom: 1rem;
}

.status-label {
  font-size: 0.9rem;
  color: #606266;
  margin-bottom: 0.5rem;
}

.status-value {
  font-size: 1.1rem;
  font-weight: bold;
  color: #303133;
  margin-bottom: 0.5rem;
}

/* 最新文章 */
.articles-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.article-item {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.article-item:hover {
  background: #f5f7fa;
}

.article-cover {
  width: 80px;
  height: 60px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.placeholder-image {
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #c0c4cc;
  font-size: 1.5rem;
}

.article-info {
  flex: 1;
}

.article-title {
  font-size: 1rem;
  font-weight: bold;
  color: #303133;
  margin-bottom: 0.5rem;
  line-height: 1.4;
}

.article-summary {
  color: #606266;
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.8rem;
  color: #909399;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

/* 快速操作 */
.actions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-item:hover {
  background: #f5f7fa;
  transform: translateY(-2px);
}

.action-icon {
  font-size: 1.5rem;
  color: #409eff;
}

.action-text {
  font-size: 0.9rem;
  color: #606266;
}

/* 系统信息 */
.info-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  color: #606266;
  font-size: 0.9rem;
}

.info-value {
  color: #303133;
  font-weight: 500;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .main-content {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .dashboard-container {
    padding: 1rem;
  }

  .stats-grid {
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  }

  .actions-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

/* 暗色主题支持 */
html.dark .dashboard-container {
  background: #1d1e1f;
}

html.dark .stat-card,
html.dark .status-card,
html.dark .articles-card,
html.dark .quick-actions-card,
html.dark .system-info-card {
  background: #2b2b2c;
  color: #e5eaf3;
}

html.dark .page-title,
html.dark .card-title,
html.dark .stat-number,
html.dark .status-value,
html.dark .article-title,
html.dark .info-value {
  color: #e5eaf3;
}

html.dark .page-subtitle,
html.dark .stat-label,
html.dark .status-label,
html.dark .article-summary,
html.dark .action-text,
html.dark .info-label {
  color: #a8abb2;
}

html.dark .article-item:hover,
html.dark .action-item:hover {
  background: #3a3b3c;
}
</style>
