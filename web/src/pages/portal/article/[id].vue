<script setup lang="ts">
import type { SysArticle } from '~/types/portal'
import { ElMessage, ElSkeleton } from 'element-plus'
import { onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import PortalLayout from '~/components/layouts/PortalLayout.vue'
import { formatDate } from '~/utils/format'
import request from '~/utils/request'

const route = useRoute()

// 响应式数据
const article = ref<SysArticle | null>(null)
const loading = ref(false)
const isLiked = ref(false)
const viewCount = ref(0)
const likeCount = ref(0)

// 获取文章详情
async function loadArticle() {
  const articleId = route.params.id as string
  if (!articleId) {
    ElMessage.error('文章ID不能为空')
    return
  }

  loading.value = true
  try {
    const res: any = await request.get(`/portal/articles/${articleId}`)
    if (res.code === 0) {
      article.value = res.data.article
      viewCount.value = res.data.viewCount || 0
      likeCount.value = res.data.likeCount || 0
      isLiked.value = res.data.isLiked || false
    }
    else {
      ElMessage.error('加载文章详情失败')
    }
  }
  catch (error) {
    ElMessage.error('加载文章详情失败')
    console.error('加载文章详情失败:', error)
  }
  finally {
    loading.value = false
  }
}

// 点赞文章
async function likeArticle() {
  if (!article.value || isLiked.value)
    return

  try {
    const res: any = await request.post(`/portal/articles/${article.value.ID}/like`)
    if (res.code === 0) {
      likeCount.value++
      isLiked.value = true
      ElMessage.success('点赞成功')
    }
  }
  catch (error) {
    ElMessage.error('点赞失败')
    console.error('点赞失败:', error)
  }
}

// 分享文章
function shareArticle() {
  if (!article.value)
    return

  const url = window.location.href
  const title = article.value.title

  if (navigator.share) {
    navigator.share({
      title,
      url,
    })
  }
  else {
    // 复制链接到剪贴板
    navigator.clipboard.writeText(url).then(() => {
      ElMessage.success('链接已复制到剪贴板')
    }).catch(() => {
      ElMessage.error('复制链接失败')
    })
  }
}
// 监听路由变化
watch(() => route.params.id, () => {
  loadArticle()
})

// 页面初始化
onMounted(() => {
  loadArticle()
})
</script>

<template>
  <PortalLayout>
    <div class="article-container">
      <div v-if="article" class="article-content">
        <!-- 文章头部 -->
        <div class="article-header">
          <h1 class="article-title">
            {{ article.title }}
          </h1>
          <div class="article-meta">
            <span class="meta-item">
              <el-icon><User /></el-icon>
              {{ article.author?.nickname || '管理员' }}
            </span>
            <span class="meta-item">
              <el-icon><Calendar /></el-icon>
              {{ formatDate(article.CreatedAt) }}
            </span>
            <span class="meta-item">
              <el-icon><View /></el-icon>
              {{ viewCount }} 次浏览
            </span>
            <span class="meta-item">
              <el-icon><Star /></el-icon>
              {{ likeCount }} 次点赞
            </span>
          </div>
          <div class="article-tags">
            <el-tag
              v-for="tag in article.tags"
              :key="tag.ID"
              :color="tag.color"
              size="small"
            >
              {{ tag.name }}
            </el-tag>
          </div>
        </div>

        <!-- 文章封面图 -->
        <div v-if="article.coverImage" class="article-cover">
          <img :src="article.coverImage" :alt="article.title" class="cover-image">
        </div>

        <!-- 文章摘要 -->
        <div v-if="article.summary" class="article-summary">
          <p>{{ article.summary }}</p>
        </div>

        <!-- 文章正文 -->
        <div class="article-body" v-html="article.content" />

        <!-- 文章底部 -->
        <div class="article-footer">
          <div class="article-category">
            <span>分类：</span>
            <el-tag v-if="article.category" :color="article.category.color" size="small">
              {{ article.category.name }}
            </el-tag>
          </div>
          <div class="article-actions">
            <el-button
              :type="isLiked ? 'success' : 'primary'"
              size="small"
              :disabled="isLiked"
              @click="likeArticle"
            >
              <el-icon><Star /></el-icon>
              {{ isLiked ? '已点赞' : '点赞' }} ({{ likeCount }})
            </el-button>
            <el-button size="small" @click="shareArticle">
              <el-icon><Share /></el-icon>
              分享
            </el-button>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-else-if="loading" class="loading-container">
        <ElSkeleton :rows="10" animated />
      </div>

      <!-- 错误状态 -->
      <div v-else class="error-container">
        <el-empty description="文章不存在或已被删除" />
      </div>
    </div>
  </PortalLayout>
</template>

<style scoped>
.article-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 2rem;
}

.article-content {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.article-header {
  padding: 2rem 2rem 1rem;
  border-bottom: 1px solid #f0f0f0;
}

.article-title {
  font-size: 2rem;
  font-weight: bold;
  color: #333;
  margin: 0 0 1rem;
  line-height: 1.4;
}

.article-meta {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 1rem;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #666;
  font-size: 0.9rem;
}

.article-tags {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.article-cover {
  padding: 0 2rem;
}

.cover-image {
  width: 100%;
  height: 300px;
  object-fit: cover;
  border-radius: 8px;
}

.article-summary {
  padding: 1rem 2rem;
  background-color: #f8f9fa;
  border-left: 4px solid var(--portal-primary-color, #409eff);
  margin: 1rem 2rem;
}

.article-summary p {
  margin: 0;
  color: #666;
  font-style: italic;
  line-height: 1.6;
}

.article-body {
  padding: 2rem;
  line-height: 1.8;
  color: #333;
}

.article-body :deep(h1),
.article-body :deep(h2),
.article-body :deep(h3),
.article-body :deep(h4),
.article-body :deep(h5),
.article-body :deep(h6) {
  margin: 1.5rem 0 1rem;
  color: #333;
}

.article-body :deep(p) {
  margin: 1rem 0;
}

.article-body :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
  margin: 1rem 0;
}

.article-body :deep(blockquote) {
  border-left: 4px solid var(--portal-primary-color, #409eff);
  padding-left: 1rem;
  margin: 1rem 0;
  color: #666;
  font-style: italic;
}

.article-body :deep(code) {
  background-color: #f5f5f5;
  padding: 0.2rem 0.4rem;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
}

.article-body :deep(pre) {
  background-color: #f5f5f5;
  padding: 1rem;
  border-radius: 4px;
  overflow-x: auto;
  margin: 1rem 0;
}

.article-body :deep(pre code) {
  background: none;
  padding: 0;
}

.article-footer {
  padding: 1rem 2rem 2rem;
  border-top: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.article-category {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #666;
}

.article-actions {
  display: flex;
  gap: 1rem;
}

.loading-container {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.error-container {
  background: white;
  border-radius: 8px;
  padding: 4rem 2rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  text-align: center;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .article-container {
    padding: 0 1rem;
  }

  .article-header {
    padding: 1.5rem 1.5rem 1rem;
  }

  .article-title {
    font-size: 1.5rem;
  }

  .article-meta {
    gap: 1rem;
  }

  .article-body {
    padding: 1.5rem;
  }

  .article-footer {
    padding: 1rem 1.5rem 1.5rem;
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
