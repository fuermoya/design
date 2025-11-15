<script setup lang="ts">
import type { SysArticle, SysCategory, SysTag, SysTheme, ThemeConfig } from '~/types/portal'
import { ElMessage } from 'element-plus'
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import PortalLayout from '~/components/layouts/PortalLayout.vue'
import { formatDate } from '~/utils/format'
import request from '~/utils/request'

const route = useRoute()
const router = useRouter()

// 响应式数据
const articles = ref<SysArticle[]>([])
const categories = ref<SysCategory[]>([])
const tags = ref<SysTag[]>([])
const themes = ref<SysTheme[]>([])
const currentTheme = ref<SysTheme | null>(null)
const themeConfig = ref<ThemeConfig | null>(null)
const currentCategory = ref<SysCategory | null>(null)

const currentPage = ref(1)
const pageSize = ref(6)
const total = ref(0)
const selectedTag = ref<number | null>(null)
const loading = ref(false)

// 计算属性
const filteredArticles = computed(() => {
  let filtered = articles.value

  if (selectedTag.value) {
    filtered = filtered.filter(article =>
      article.tags?.some(tag => tag.ID === selectedTag.value),
    )
  }

  return filtered
})

// 获取分类ID
const categoryId = computed(() => {
  return Number.parseInt(route.params.id as string)
})

// 加载文章列表
async function loadArticles() {
  if (!categoryId.value)
    return

  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      categoryId: categoryId.value,
    }
    const res: any = await request.get('/portal/articles', params)
    if (res.code === 0) {
      articles.value = res.data.list
      total.value = res.data.total
    }
  }
  catch (error) {
    ElMessage.error('加载文章列表失败')
    console.error('加载文章列表失败:', error)
  }
  finally {
    loading.value = false
  }
}

// 加载分类列表
async function loadCategories() {
  try {
    const res: any = await request.get('/portal/categories')
    if (res.code === 0) {
      categories.value = res.data.list
      // 设置当前分类
      currentCategory.value = categories.value.find(cat => cat.ID === categoryId.value) || null
    }
  }
  catch (error) {
    console.error('加载分类列表失败:', error)
  }
}

// 加载标签列表
async function loadTags() {
  try {
    const res: any = await request.get('/portal/tags')
    if (res.code === 0) {
      tags.value = res.data.list
    }
  }
  catch (error) {
    console.error('加载标签列表失败:', error)
  }
}

// 加载主题列表
async function loadThemes() {
  try {
    const res: any = await request.get('/portal/themes')
    if (res.code === 0) {
      themes.value = res.data.list
    }
  }
  catch (error) {
    console.error('加载主题列表失败:', error)
  }
}

// 获取当前激活的主题
async function loadActiveTheme() {
  try {
    const res: any = await request.get('/portal/theme')
    if (res.code === 0) {
      themeConfig.value = res.data.theme
    }
  }
  catch (error) {
    console.error('加载当前主题失败:', error)
  }
}

// 过滤文章
function filterArticles() {
  currentPage.value = 1
  loadArticles()
}

// 查看文章详情
function viewArticle(article: SysArticle) {
  router.push(`/portal/article/${article.ID}`)
}

// 处理分页大小变化
function handleSizeChange(newSize: number) {
  pageSize.value = newSize
  currentPage.value = 1
  loadArticles()
}

// 处理当前页变化
function handleCurrentChange(newPage: number) {
  currentPage.value = newPage
  loadArticles()
}

// 监听路由变化
watch(() => route.params.id, () => {
  if (categoryId.value) {
    loadArticles()
    loadCategories()
  }
})

// 页面初始化
onMounted(() => {
  loadArticles()
  loadCategories()
  loadTags()
  loadThemes()
  loadActiveTheme()
})
</script>

<template>
  <PortalLayout>
    <div class="category-content">
      <!-- 页面头部 -->
      <div class="page-header">
        <div class="header-content">
          <h1 class="page-title">
            {{ currentCategory?.name || '分类文章' }}
          </h1>
          <p class="page-subtitle">
            {{ currentCategory?.description || '浏览该分类下的所有文章' }}
          </p>
        </div>
      </div>

      <!-- 文章列表区域 -->
      <section class="articles-section">
        <div class="section-container">
          <!-- <div class="articles-header">
            <h3 class="section-title">
              {{ currentCategory?.name || '分类文章' }}
            </h3>
            <div class="articles-filter">
              <el-select v-model="selectedTag" placeholder="选择标签" clearable @change="filterArticles">
                <el-option
                  v-for="tag in tags"
                  :key="tag.ID"
                  :label="tag.name"
                  :value="tag.ID"
                />
              </el-select>
            </div>
          </div> -->

          <!-- 文章列表 -->
          <div v-if="loading" class="loading-container">
            <el-skeleton :rows="6" animated />
          </div>

          <div v-else-if="filteredArticles.length > 0" class="articles-grid">
            <article
              v-for="article in filteredArticles"
              :key="article.ID"
              class="article-card"
              @click="viewArticle(article)"
            >
              <div class="article-image">
                <img
                  v-if="article.coverImage"
                  :src="article.coverImage"
                  :alt="article.title"
                  class="cover-image"
                >
                <div v-else class="placeholder-image">
                  <el-icon><Picture /></el-icon>
                </div>
              </div>
              <div class="article-content">
                <div class="article-meta">
                  <span class="article-category">{{ article.category?.name }}</span>
                  <span class="article-date">{{ formatDate(article.CreatedAt) }}</span>
                </div>
                <h4 class="article-title">
                  {{ article.title }}
                </h4>
                <p class="article-summary">
                  {{ article.summary }}
                </p>
                <div class="article-tags">
                  <el-tag
                    v-for="tag in article.tags"
                    :key="tag.ID"
                    :color="tag.color"
                    size="small"
                    class="article-tag"
                  >
                    {{ tag.name }}
                  </el-tag>
                </div>
                <div class="article-stats">
                  <span class="stat-item">
                    <el-icon><View /></el-icon>
                    {{ article.viewCount }}
                  </span>
                  <span class="stat-item">
                    <el-icon><Star /></el-icon>
                    {{ article.likeCount }}
                  </span>
                </div>
              </div>
            </article>
          </div>

          <div v-else class="empty-container">
            <el-empty description="暂无文章" />
          </div>

          <!-- 分页 -->
          <div v-if="total > pageSize" class="pagination-wrapper">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="total"
              :page-sizes="[6, 12, 24, 48]"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </div>
      </section>
    </div>
  </PortalLayout>
</template>

<style scoped>
.category-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

/* 页面头部 */
.page-header {
  background: linear-gradient(
    135deg,
    var(--portal-primary-color, #409eff) 0%,
    var(--portal-secondary-color, #67c23a) 100%
  );
  color: white;
  padding: 3rem 2rem;
  border-radius: 12px;
  margin-bottom: 2rem;
  text-align: center;
}

.header-content {
  max-width: 800px;
  margin: 0 auto;
}

.page-title {
  font-size: 2.5rem;
  font-weight: bold;
  margin: 0 0 1rem;
}

.page-subtitle {
  font-size: 1.1rem;
  opacity: 0.9;
  margin: 0;
  line-height: 1.6;
}

/* 文章区域 */
.articles-section {
  margin-bottom: 3rem;
}

.section-container {
  max-width: 1200px;
  margin: 0 auto;
}

.articles-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.section-title {
  font-size: 1.8rem;
  font-weight: bold;
  color: var(--portal-text-color, #303133);
  margin: 0;
}

.articles-filter {
  display: flex;
  gap: 1rem;
  align-items: center;
}

/* 文章网格 */
.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.article-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
}

.article-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.article-image {
  height: 200px;
  overflow: hidden;
  position: relative;
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.article-card:hover .cover-image {
  transform: scale(1.05);
}

.placeholder-image {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
  font-size: 3rem;
}

.article-content {
  padding: 1.5rem;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  margin-bottom: 1rem;
  font-size: 0.9rem;
  color: #909399;
}

.article-category {
  color: var(--portal-primary-color, #409eff);
  font-weight: 500;
}

.article-title {
  font-size: 1.2rem;
  font-weight: bold;
  color: var(--portal-text-color, #303133);
  margin: 0 0 1rem;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-summary {
  color: #606266;
  line-height: 1.6;
  margin: 0 0 1rem;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-tags {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.article-tag {
  font-size: 0.8rem;
}

.article-stats {
  display: flex;
  gap: 1rem;
  color: #909399;
  font-size: 0.9rem;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.loading-container {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.empty-container {
  background: white;
  border-radius: 12px;
  padding: 4rem 2rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .category-content {
    padding: 0 1rem;
  }

  .page-header {
    padding: 2rem 1rem;
  }

  .page-title {
    font-size: 2rem;
  }

  .articles-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .articles-grid {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }

  .article-content {
    padding: 1rem;
  }
}
</style>
