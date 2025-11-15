<script setup lang="ts">
import type { SysArticle, SysCategory, SysTag } from '~/types/portal'
import { ElMessage, ElMessageBox } from 'element-plus'
import { nextTick, onMounted, ref, watch } from 'vue'
import HtmlEditor from '~/components/HtmlEditor.vue'
import SimpleMediaLibrary from '~/components/SimpleMediaLibrary.vue'
import { formatDate } from '~/utils/format'
import request from '~/utils/request'

// 响应式数据
const articles = ref<SysArticle[]>([])
const categories = ref<SysCategory[]>([])
const tags = ref<SysTag[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 媒体库相关
const mediaLibraryVisible = ref(false)
const mediaType = ref<'image' | 'video'>('image')

// 搜索条件
const searchForm = ref({
  keyword: '',
  categoryId: undefined as number | undefined,
  status: undefined as number | undefined,
})

// 表单数据
const formData = ref({
  id: 0,
  title: '',
  content: '',
  summary: '',
  coverImage: '',
  categoryId: undefined as number | undefined,
  tagIds: [] as number[],
  status: 2, // 默认草稿状态
  sort: 0,
})

// 重置表单数据
function resetFormData() {
  formData.value = {
    id: 0,
    title: '',
    content: '',
    summary: '',
    coverImage: '',
    categoryId: undefined,
    tagIds: [],
    status: 2,
    sort: 0,
  }
}

// 状态选项
const statusOptions = [
  { label: '草稿', value: 2 },
  { label: '已发布', value: 1 },
  { label: '已下架', value: 3 },
]

// 加载文章列表
async function loadArticles() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      ...searchForm.value,
    }
    const res: any = await request.get('/sysArticle/getSysArticleList', params)
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
    const res: any = await request.get('/sysCategory/getSysCategoryList', {
      page: 1,
      pageSize: 1000,
    })
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
    const res: any = await request.get('/sysTag/getSysTagList', {
      page: 1,
      pageSize: 1000,
    })
    if (res.code === 0) {
      tags.value = res.data.list
    }
  }
  catch (error) {
    console.error('加载标签列表失败:', error)
  }
}

// 自动生成摘要
function generateSummary(content: string) {
  // 移除HTML标签，获取纯文本
  const textContent = content.replace(/<[^>]*>/g, '')
  // 截取前200个字符作为摘要
  return textContent.length > 200 ? `${textContent.substring(0, 200)}...` : textContent
}

// 监听内容变化，自动生成摘要
watch(() => formData.value.content, (newContent: string) => {
  if (newContent && !formData.value.summary) {
    formData.value.summary = generateSummary(newContent)
  }
})

// 打开媒体库选择图片
function openMediaLibrary(type: 'image' | 'video') {
  mediaType.value = type
  mediaLibraryVisible.value = true
}

// 处理媒体文件选择
function handleMediaSelect(files: any[]) {
  if (files.length > 0) {
    const file = files[0]
    if (mediaType.value === 'image') {
      formData.value.coverImage = file.url
    }
  }
}

// 打开新增对话框
function openAddDialog() {
  dialogTitle.value = '新增文章'
  isEdit.value = false
  resetFormData()
  dialogVisible.value = true
}

// 打开编辑对话框
function openEditDialog(row: SysArticle) {
  dialogTitle.value = '编辑文章'
  isEdit.value = true

  // 先打开对话框
  dialogVisible.value = true

  // 使用nextTick确保对话框已打开后再设置数据
  nextTick(() => {
    // 确保所有字段都正确回显
    formData.value = {
      id: row.ID,
      title: row.title || '',
      content: row.content || '',
      summary: row.summary || '',
      coverImage: row.coverImage || '',
      categoryId: row.categoryId ? Number(row.categoryId) : undefined,
      tagIds: row.tags?.map(tag => tag.ID) || [],
      status: row.status || 2,
      sort: row.sort || 0,
    }
  })
}

// 保存文章
async function saveArticle() {
  try {
    if (isEdit.value) {
      await request.put('/sysArticle/updateSysArticle', formData.value)
      ElMessage.success('更新成功')
    }
    else {
      await request.post('/sysArticle/createSysArticle', formData.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadArticles()
  }
  catch (error) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    console.error('保存文章失败:', error)
  }
}

// 删除文章
async function deleteArticle(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await request.del(`/sysArticle/deleteSysArticle?id=${id}`, {})

    loadArticles()
  }
  catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error('删除文章失败:', error)
    }
  }
}

// 批量删除
async function batchDelete() {
  const selectedIds = articles.value
    .filter((article: SysArticle) => article.selected)
    .map((article: SysArticle) => article.ID)

  if (selectedIds.length === 0) {
    ElMessage.warning('请选择要删除的文章')
    return
  }

  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedIds.length} 篇文章吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await request.del('/sysArticle/deleteSysArticleByIds', {
      ids: selectedIds,
    })

    loadArticles()
  }
  catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 获取状态文本
function getStatusText(status: number) {
  const statusMap: Record<number, string> = {
    1: '已发布',
    2: '草稿',
    3: '已下架',
  }
  return statusMap[status] || '未知'
}

// 获取状态类型
function getStatusType(status: number) {
  const typeMap: Record<number, 'success' | 'warning' | 'danger'> = {
    1: 'success',
    2: 'warning',
    3: 'danger',
  }
  return typeMap[status] || 'info'
}

// 处理表格选择变化
function handleSelectionChange(selection: SysArticle[]) {
  articles.value.forEach((item: SysArticle) => {
    item.selected = selection.includes(item)
  })
}

// 页面初始化
onMounted(() => {
  loadArticles()
  loadCategories()
  loadTags()
})
</script>

<template>
  <div class="upload-container">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入文章标题或内容关键词"
            clearable
            @keyup.enter="loadArticles"
          />
        </el-form-item>
        <el-form-item label="分类">
          <el-select v-model="searchForm.categoryId" placeholder="请选择分类" clearable>
            <el-option
              v-for="category in categories"
              :key="category.ID"
              :label="category.name"
              :value="category.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option
              v-for="option in statusOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadArticles">
            搜索
          </el-button>
          <el-button @click="searchForm = { keyword: '', categoryId: undefined, status: undefined }">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 文章列表 -->
    <el-card class="list-card">
      <div class="table-header">
        <div class="table-actions">
          <el-button type="primary" @click="openAddDialog">
            新增文章
          </el-button>
          <el-button type="danger" @click="batchDelete">
            批量删除
          </el-button>
        </div>
      </div>
      <el-table
        v-loading="loading"
        :data="articles"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="category.name" label="分类" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="viewCount" label="阅读量" width="100" />
        <el-table-column prop="likeCount" label="点赞数" width="100" />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="CreatedAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button type="primary" size="small" @click="openEditDialog(row)">
                编辑
              </el-button>
              <el-button type="danger" size="small" @click="deleteArticle(row.ID)">
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadArticles"
          @current-change="loadArticles"
        />
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="80%"
      :close-on-click-modal="false"
    >
      <el-form :model="formData" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="16">
            <el-form-item label="文章标题" required>
              <el-input v-model="formData.title" placeholder="请输入文章标题" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="分类" required>
              <el-select
                v-model="formData.categoryId"
                placeholder="请选择分类"
              >
                <el-option
                  v-for="category in categories"
                  :key="category.ID"
                  :label="category.name"
                  :value="category.ID"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="状态">
              <el-select v-model="formData.status">
                <el-option
                  v-for="option in statusOptions"
                  :key="option.value"
                  :label="option.label"
                  :value="option.value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="排序">
              <el-input-number v-model="formData.sort" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="标签">
              <el-select v-model="formData.tagIds" multiple placeholder="请选择标签">
                <el-option
                  v-for="tag in tags"
                  :key="tag.ID"
                  :label="tag.name"
                  :value="tag.ID"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="封面图片">
          <div class="cover-image-section">
            <el-input v-model="formData.coverImage" placeholder="请输入封面图片URL" />
            <el-button type="primary" @click="openMediaLibrary('image')">
              从媒体库选择
            </el-button>
          </div>
          <div v-if="formData.coverImage" class="cover-preview">
            <img :src="formData.coverImage" alt="封面预览" class="cover-image">
          </div>
        </el-form-item>

        <el-form-item label="文章摘要">
          <el-input
            v-model="formData.summary"
            type="textarea"
            :rows="3"
            placeholder="文章摘要将根据内容自动生成，也可手动编辑"
          />
        </el-form-item>

        <el-form-item label="文章内容" required>
          <HtmlEditor
            v-model="formData.content"
            placeholder="请输入文章内容..."
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="saveArticle">
          保存
        </el-button>
      </template>
    </el-dialog>

    <!-- 媒体库组件 -->
    <SimpleMediaLibrary
      v-model:visible="mediaLibraryVisible"
      :type="mediaType"
      :multiple="false"
      @select="handleMediaSelect"
    />
  </div>
</template>

<style scoped>
.cover-image-section {
  display: flex;
  gap: 12px;
  align-items: center;
}

.cover-image-section .el-input {
  flex: 1;
}

.cover-preview {
  margin-top: 12px;
}

.cover-image {
  max-width: 200px;
  max-height: 120px;
  border-radius: 4px;
  border: 1px solid #e4e7ed;
}
</style>
