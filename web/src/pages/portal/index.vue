<script setup lang="ts">
import type { SysArticle, SysCategory, SysTag, SysTheme, ThemeConfig } from '~/types/portal'
import { ElMessage } from 'element-plus'
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import PortalLayout from '~/components/layouts/PortalLayout.vue'
import { formatDate } from '~/utils/format'
import request from '~/utils/request'

const router = useRouter()

// 响应式数据
const articles = ref<SysArticle[]>([])
const categories = ref<SysCategory[]>([])
const tags = ref<SysTag[]>([])
const themes = ref<SysTheme[]>([])
const currentTheme = ref<SysTheme | null>(null)
const themeConfig = ref<ThemeConfig | null>(null)

const currentPage = ref(1)
const pageSize = ref(6)
const total = ref(0)
const selectedCategory = ref<number | null>(null)
const selectedTag = ref<SysTag | null>(null)

// 留言相关数据
const messageFormRef = ref()
const submitting = ref(false)
const messageForm = ref({
  name: '',
  email: '',
  phone: '',
  content: '',
})

// 留言表单验证规则
const messageRules = {
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '姓名长度在 2 到 20 个字符', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email' as const, message: '请输入正确的邮箱地址', trigger: 'blur' },
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' },
  ],
  content: [
    { required: true, message: '请输入留言内容', trigger: 'blur' },
    { min: 10, max: 500, message: '留言内容长度在 10 到 500 个字符', trigger: 'blur' },
  ],
}

// 提交留言
async function submitMessage() {
  if (!messageFormRef.value)
    return

  try {
    await messageFormRef.value.validate()
    submitting.value = true

    const res: any = await request.post('/portal/message', messageForm.value)
    if (res.code === 0) {
      ElMessage.success('留言提交成功，我们会尽快回复您！')
      resetMessageForm()
    }
    else {
      ElMessage.error(res.message || '留言提交失败')
    }
  }
  catch (error) {
    ElMessage.error('留言提交失败')
    console.error('提交留言失败:', error)
  }
  finally {
    submitting.value = false
  }
}

// 重置留言表单
function resetMessageForm() {
  if (messageFormRef.value) {
    messageFormRef.value.resetFields()
  }
}

// 计算属性
const filteredArticles = computed(() => {
  let filtered = articles.value

  if (selectedCategory.value) {
    filtered = filtered.filter(article => article.categoryId === selectedCategory.value)
  }

  if (selectedTag.value) {
    filtered = filtered.filter(article =>
      article.tags?.some(tag => tag.ID === selectedTag.value?.ID),
    )
  }

  return filtered
})

// 加载文章列表
async function loadArticles() {
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
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
}

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
      currentTheme.value = res.data.theme
      if (currentTheme.value?.config) {
        themeConfig.value = JSON.parse(currentTheme.value.config)
      }
    }
  }
  catch (error) {
    console.error('加载当前主题失败:', error)
  }
}

// 切换主题
async function switchTheme() {
  try {
    const nextTheme = themes.value.find(theme => !theme.isActive)
    if (nextTheme) {
      const res: any = await request.post('/portal/switchTheme', {
        themeId: nextTheme.ID,
      })
      if (res.code === 0) {
        ElMessage.success('主题切换成功')
        await loadActiveTheme()
      }
    }
  }
  catch (error) {
    ElMessage.error('主题切换失败')
    console.error('切换主题失败:', error)
  }
}

// 选择分类
function selectCategory(category: SysCategory) {
  selectedCategory.value = category.ID
  selectedTag.value = null
  currentPage.value = 1
  loadArticles()
}

// 选择标签
function selectTag(tag: SysTag) {
  selectedTag.value = tag
  selectedCategory.value = null
  currentPage.value = 1
  loadArticles()
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
    <div class="portal-content">
      <!-- 轮播图区域 -->
      <section class="hero-section">
        <div class="hero-content">
          <h2 class="hero-title">
            {{ themeConfig?.siteName || '欢迎来到企业门户' }}
          </h2>
          <p class="hero-subtitle">
            {{ themeConfig?.siteDescription || '专业的服务，卓越的品质' }}
          </p>
          <el-button type="primary" size="large" class="hero-btn">
            了解更多
          </el-button>
        </div>
      </section>

      <!-- 企业服务区域 -->
      <section class="services-section">
        <div class="section-container">
          <h3 class="section-title">
            我们的服务
          </h3>
          <div class="services-grid">
            <div class="service-card">
              <div class="service-icon">
                <el-icon><Star /></el-icon>
              </div>
              <h4 class="service-title">
                技术咨询
              </h4>
              <p class="service-desc">
                为企业提供专业的技术解决方案，助力数字化转型
              </p>
            </div>
            <div class="service-card">
              <div class="service-icon">
                <el-icon><User /></el-icon>
              </div>
              <h4 class="service-title">
                人才培养
              </h4>
              <p class="service-desc">
                定制化培训课程，提升团队专业技能和创新能力
              </p>
            </div>
            <div class="service-card">
              <div class="service-icon">
                <el-icon><Folder /></el-icon>
              </div>
              <h4 class="service-title">
                项目管理
              </h4>
              <p class="service-desc">
                全流程项目管理服务，确保项目高质量交付
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- 企业数据展示 -->
      <section class="stats-section">
        <div class="section-container">
          <div class="stats-grid">
            <div class="stat-card">
              <div class="stat-number">
                500+
              </div>
              <div class="stat-label">
                服务客户
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-number">
                50+
              </div>
              <div class="stat-label">
                专业团队
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-number">
                10年+
              </div>
              <div class="stat-label">
                行业经验
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-number">
                98%
              </div>
              <div class="stat-label">
                客户满意度
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 最新动态 -->
      <section v-if="total > 0" class="news-section">
        <div class="section-container">
          <div class="news-grid">
            <article
              v-for="article in filteredArticles"
              :key="article.ID"
              class="news-card"
              @click="viewArticle(article)"
            >
              <div class="news-image">
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
              <div class="news-content">
                <div class="news-meta">
                  <span class="news-category">{{ article.category?.name }}</span>
                  <span class="news-date">{{ formatDate(article.CreatedAt) }}</span>
                </div>
                <h4 class="news-title">
                  {{ article.title }}
                </h4>
                <p class="news-summary">
                  {{ article.summary }}
                </p>
                <div class="news-tags">
                  <el-tag
                    v-for="tag in article.tags"
                    :key="tag.ID"
                    :color="tag.color"
                    size="small"
                    class="news-tag"
                  >
                    {{ tag.name }}
                  </el-tag>
                </div>
                <div class="news-stats">
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

          <!-- 分页 -->
          <div class="pagination-wrapper">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="total"
              :page-sizes="[6, 12, 24, 48]"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="loadArticles"
              @current-change="loadArticles"
            />
          </div>
        </div>
      </section>

      <!-- 合作伙伴 -->
      <section class="partners-section">
        <div class="section-container">
          <h3 class="section-title">
            合作伙伴
          </h3>
          <div class="partners-grid">
            <div class="partner-card">
              <div class="partner-logo">
                合作伙伴A
              </div>
              <p class="partner-desc">
                行业领先的技术解决方案提供商
              </p>
            </div>
            <div class="partner-card">
              <div class="partner-logo">
                合作伙伴B
              </div>
              <p class="partner-desc">
                专业的咨询服务公司
              </p>
            </div>
            <div class="partner-card">
              <div class="partner-logo">
                合作伙伴C
              </div>
              <p class="partner-desc">
                创新型企业孵化平台
              </p>
            </div>
            <div class="partner-card">
              <div class="partner-logo">
                合作伙伴D
              </div>
              <p class="partner-desc">
                国际化的技术交流平台
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- 留言区域 -->
      <section class="message-section">
        <div class="section-container">
          <h3 class="section-title">
            联系我们
          </h3>
          <p class="section-subtitle">
            如果您有任何问题或建议，请随时与我们联系
          </p>
          <div class="message-form-wrapper">
            <el-form
              ref="messageFormRef"
              :model="messageForm"
              :rules="messageRules"
              label-width="80px"
              class="message-form"
            >
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="姓名" prop="name">
                    <el-input
                      v-model="messageForm.name"
                      placeholder="请输入您的姓名"
                      clearable
                    />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="邮箱" prop="email">
                    <el-input
                      v-model="messageForm.email"
                      placeholder="请输入您的邮箱"
                      clearable
                    />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item label="电话" prop="phone">
                <el-input
                  v-model="messageForm.phone"
                  placeholder="请输入您的电话号码（可选）"
                  clearable
                />
              </el-form-item>
              <el-form-item label="留言内容" prop="content">
                <el-input
                  v-model="messageForm.content"
                  type="textarea"
                  :rows="4"
                  placeholder="请输入您的留言内容"
                  maxlength="500"
                  show-word-limit
                />
              </el-form-item>
              <el-form-item>
                <el-button
                  type="primary"
                  size="large"
                  :loading="submitting"
                  @click="submitMessage"
                >
                  提交留言
                </el-button>
                <el-button size="large" @click="resetMessageForm">
                  重置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
      </section>
    </div>
  </PortalLayout>
</template>

<style scoped>
.portal-container {
  min-height: 100vh;
  background-color: var(--portal-bg-color, #f5f7fa);
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
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
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
  border-radius: 4px;
  transition: all 0.3s ease;
}

.nav-item:hover,
.nav-item.active {
  background-color: rgba(255, 255, 255, 0.2);
}

/* 主要内容区域 */
.portal-main {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

/* 轮播图区域 */
.hero-section {
  background: linear-gradient(
    135deg,
    var(--portal-primary-color, #409eff) 0%,
    var(--portal-secondary-color, #67c23a) 100%
  );
  color: white;
  padding: 4rem 2rem;
  border-radius: 12px;
  text-align: center;
  margin-bottom: 3rem;
}

.hero-title {
  font-size: 3rem;
  font-weight: bold;
  margin-bottom: 1rem;
}

.hero-subtitle {
  font-size: 1.2rem;
  margin-bottom: 2rem;
  opacity: 0.9;
}

.hero-btn {
  font-size: 1.1rem;
  padding: 0.8rem 2rem;
}

/* 分类区域 */
.categories-section {
  margin-bottom: 3rem;
}

.section-container {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.section-title {
  font-size: 1.5rem;
  font-weight: bold;
  margin-bottom: 1.5rem;
  color: var(--portal-text-color, #303133);
}

.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.category-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 1.5rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.category-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  border-color: var(--portal-primary-color, #409eff);
}

.category-icon {
  font-size: 2rem;
  color: var(--portal-primary-color, #409eff);
  margin-bottom: 1rem;
}

.category-name {
  font-size: 1.2rem;
  font-weight: bold;
  margin-bottom: 0.5rem;
  color: var(--portal-text-color, #303133);
}

.category-desc {
  color: #606266;
  font-size: 0.9rem;
}

/* 标签区域 */
.tags-section {
  margin-bottom: 3rem;
}

.tags-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag-item {
  cursor: pointer;
  transition: all 0.3s ease;
}

.tag-item:hover {
  transform: scale(1.05);
}

/* 企业服务区域 */
.services-section {
  margin-bottom: 3rem;
}

.services-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.service-card {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  text-align: center;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.service-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.service-icon {
  font-size: 3rem;
  color: var(--portal-primary-color, #409eff);
  margin-bottom: 1rem;
}

.service-title {
  font-size: 1.3rem;
  font-weight: bold;
  margin-bottom: 1rem;
  color: var(--portal-text-color, #303133);
}

.service-desc {
  color: #606266;
  line-height: 1.6;
}

/* 企业数据展示 */
.stats-section {
  margin-bottom: 3rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 2rem;
}

.stat-card {
  background: linear-gradient(
    135deg,
    var(--portal-primary-color, #409eff) 0%,
    var(--portal-secondary-color, #67c23a) 100%
  );
  color: white;
  border-radius: 12px;
  padding: 2rem;
  text-align: center;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.stat-number {
  font-size: 2.5rem;
  font-weight: bold;
  margin-bottom: 0.5rem;
}

.stat-label {
  font-size: 1rem;
  opacity: 0.9;
}

/* 最新动态 */
.news-section {
  margin-bottom: 3rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.news-filter {
  display: flex;
  gap: 1rem;
}

.news-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.news-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
}

.news-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.news-image {
  height: 200px;
  overflow: hidden;
}

.news-content {
  padding: 1.5rem;
}

.news-meta {
  display: flex;
  gap: 1rem;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
  color: #909399;
}

.news-title {
  font-size: 1.2rem;
  font-weight: bold;
  margin-bottom: 0.5rem;
  color: var(--portal-text-color, #303133);
  line-height: 1.4;
}

.news-summary {
  color: #606266;
  margin-bottom: 1rem;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.news-tags {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.news-tag {
  font-size: 0.8rem;
}

.news-stats {
  display: flex;
  gap: 1rem;
  font-size: 0.9rem;
  color: #909399;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
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
  font-size: 3rem;
  color: #c0c4cc;
}

/* 合作伙伴 */
.partners-section {
  margin-bottom: 3rem;
}

.partners-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
}

.partner-card {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  text-align: center;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.partner-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.partner-logo {
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--portal-primary-color, #409eff);
  margin-bottom: 1rem;
}

.partner-desc {
  color: #606266;
  line-height: 1.6;
}

/* 留言区域 */
.message-section {
  margin-bottom: 3rem;
}

.message-form-wrapper {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.message-form {
  max-width: 800px;
  margin: 0 auto;
}

.section-subtitle {
  color: #606266;
  margin-bottom: 2rem;
  font-size: 1rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header-container {
    flex-direction: column;
    gap: 1rem;
  }

  .main-nav {
    gap: 1rem;
  }

  .hero-title {
    font-size: 2rem;
  }

  .hero-subtitle {
    font-size: 1rem;
  }

  .articles-header {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }

  .articles-grid {
    grid-template-columns: 1fr;
  }

  .footer-content {
    grid-template-columns: 1fr;
  }
}

/* 暗色主题支持 */
html.dark .portal-container {
  background-color: #1d1e1f;
}

html.dark .section-container,
html.dark .article-card {
  background: #2b2b2c;
  color: #e5eaf3;
}

html.dark .category-card {
  background: #3a3b3c;
}

html.dark .article-title,
html.dark .category-name,
html.dark .section-title {
  color: #e5eaf3;
}

html.dark .article-summary,
html.dark .category-desc {
  color: #a8abb2;
}
</style>
