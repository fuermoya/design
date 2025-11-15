<script lang="ts" setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, onUnmounted, reactive, ref } from 'vue'
import { formatDate } from '~/utils/format'
import request from '~/utils/request'

// 媒体库上传页面组件
interface FileItem {
  ID: number
  name: string
  url: string
  tag: string
  key: string
  created_at: string
  updated_at: string
}

interface SearchForm {
  name: string
  tag: string
  page: number
  pageSize: number
}

const fileList = ref<FileItem[]>([])
const total = ref(0)
const loading = ref(false)
const uploadLoading = ref(false)

// 预览相关状态
const previewVisible = ref(false)
const previewUrl = ref('')
const previewType = ref('')
const previewFileName = ref('')

const searchForm = reactive<SearchForm>({
  name: '',
  tag: '',
  page: 1,
  pageSize: 10,
})

const uploadForm = reactive({
  file: null as File | null,
})

// 获取文件列表
async function getFileList() {
  loading.value = true
  try {
    const res: any = await request.get('/fileUploadAndDownload/getFileList', {
      name: searchForm.name,
      tag: searchForm.tag,
      page: searchForm.page,
      pageSize: searchForm.pageSize,
    })
    if (res.code === 0) {
      fileList.value = res.data.list
      total.value = res.data.total
    }
  }
  catch (error) {
    console.error('获取文件列表失败:', error)
  }
  finally {
    loading.value = false
  }
}

// 上传文件
async function handleUpload() {
  if (!uploadForm.file) {
    ElMessage.warning('请选择要上传的文件')
    return
  }

  uploadLoading.value = true
  try {
    const formData = new FormData()
    formData.append('file', uploadForm.file)

    const res: any = await request.upload('/fileUploadAndDownload/upload', formData)
    if (res.code === 0) {
      ElMessage.success('文件上传成功')
      uploadForm.file = null
      getFileList()
    }
  }
  catch (error) {
    console.error('文件上传失败:', error)
  }
  finally {
    uploadLoading.value = false
  }
}

// 删除文件
async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除这个文件吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    const res: any = await request.del(`/fileUploadAndDownload/deleteFile?id=${id}`, { })
    if (res.code === 0) {
      getFileList()
    }
  }
  catch (error) {
    if (error !== 'cancel') {
      console.error('删除文件失败:', error)
    }
  }
}

// 下载文件
async function handleDownload(id: number, name: string) {
  try {
    // 使用下载API而不是直接使用文件URL
    const response = await fetch(`${import.meta.env.VITE_BASE_API || '/api'}/fileUploadAndDownload/download?id=${id}`, {
      method: 'GET',
      headers: {
        'x-token': localStorage.getItem('x-token') || '',
      },
    })

    if (!response.ok) {
      throw new Error('下载失败')
    }

    // 获取文件blob
    const blob = await response.blob()

    // 创建下载链接
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = name
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    // 清理URL对象
    window.URL.revokeObjectURL(url)

    ElMessage.success('下载成功')
  }
  catch (error) {
    console.error('下载文件失败:', error)
    ElMessage.error('下载失败')
  }
}

// 文件选择
function handleFileChange(uploadFile: any) {
  if (uploadFile.raw) {
    uploadForm.file = uploadFile.raw
  }
}

// 搜索
function handleSearch() {
  searchForm.page = 1
  getFileList()
}

// 重置搜索
function handleReset() {
  searchForm.name = ''
  searchForm.tag = ''
  searchForm.page = 1
  getFileList()
}

// 分页
function handleCurrentChange(page: number) {
  searchForm.page = page
  getFileList()
}

function handleSizeChange(size: number) {
  searchForm.pageSize = size
  searchForm.page = 1
  getFileList()
}

onMounted(() => {
  getFileList()
})

// 判断是否为图片文件
function isImageFile(file: File): boolean {
  return file.type.startsWith('image/')
}

// 判断文件URL是否为图片
function isImageUrl(url: string): boolean {
  const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg']
  return imageExtensions.some(ext => url.toLowerCase().includes(ext))
}

// 获取预览URL
function getPreviewUrl(file: File): string {
  return URL.createObjectURL(file)
}

// 获取服务器文件预览URL
function getServerPreviewUrl(id: number): string {
  const baseUrl = import.meta.env.VITE_BASE_API || '/api'
  return `${baseUrl}/fileUploadAndDownload/preview?id=${id}`
}

// 预览文件
function handlePreview(fileItem: FileItem) {
  previewFileName.value = fileItem.name
  // 使用预览接口
  previewUrl.value = getServerPreviewUrl(fileItem.ID)

  // 判断文件类型
  if (isImageUrl(fileItem.url)) {
    previewType.value = 'image'
  }
  else if (fileItem.tag === 'video' || fileItem.url.includes('.mp4') || fileItem.url.includes('.avi')) {
    previewType.value = 'video'
  }
  else if (fileItem.tag === 'audio' || fileItem.url.includes('.mp3') || fileItem.url.includes('.wav')) {
    previewType.value = 'audio'
  }
  else {
    previewType.value = 'other'
  }

  previewVisible.value = true
}

// 关闭预览
function closePreview() {
  previewVisible.value = false
  previewUrl.value = ''
  previewType.value = ''
  previewFileName.value = ''
}

// 复制到剪贴板
async function copyToClipboard(text: string) {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('链接已复制到剪贴板')
  }
  catch {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    ElMessage.success('链接已复制到剪贴板')
  }
}

// 在新窗口中打开文件
function openInNewWindow() {
  window.open(previewUrl.value, '_blank')
}

// 格式化文件大小
function formatFileSize(bytes: number): string {
  if (bytes === 0)
    return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return `${Number.parseFloat((bytes / k ** i).toFixed(2))} ${sizes[i]}`
}

// 组件卸载时清理资源
onUnmounted(() => {
  // 清理资源
})
</script>

<template>
  <div class="upload-container">
    <!-- 上传区域 -->
    <el-card class="upload-card" shadow="hover">
      <el-form :model="uploadForm" label-width="80px">
        <el-form-item label="选择文件">
          <el-upload
            class="upload-demo"
            :auto-upload="false"
            :show-file-list="false"
            :on-change="handleFileChange"
            accept="*/*"
          >
            <el-button type="primary">
              选择文件
            </el-button>
            <template #tip>
              <div class="el-upload__tip">
                支持任意格式文件上传
              </div>
            </template>
          </el-upload>
        </el-form-item>

        <el-form-item>
          <el-button
            type="success"
            :loading="uploadLoading"
            :disabled="!uploadForm.file"
            @click="handleUpload"
          >
            上传文件
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 文件预览区域 -->
      <div v-if="uploadForm.file" class="preview-section">
        <h4>文件预览</h4>
        <div class="preview-content">
          <!-- 图片预览 -->
          <div v-if="isImageFile(uploadForm.file)" class="image-preview">
            <img :src="getPreviewUrl(uploadForm.file)" alt="文件预览" class="preview-image">
          </div>

          <!-- 其他文件类型 -->
          <div v-else class="file-info">
            <el-descriptions :column="1" border>
              <el-descriptions-item label="文件名">
                {{ uploadForm.file.name }}
              </el-descriptions-item>
              <el-descriptions-item label="文件大小">
                {{ formatFileSize(uploadForm.file.size) }}
              </el-descriptions-item>
              <el-descriptions-item label="文件类型">
                {{ uploadForm.file.type || '未知类型' }}
              </el-descriptions-item>
              <el-descriptions-item label="最后修改时间">
                {{ new Date(uploadForm.file.lastModified).toLocaleString() }}
              </el-descriptions-item>
            </el-descriptions>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 搜索区域 -->
    <el-card class="search-card" shadow="hover">
      <el-form :model="searchForm" inline>
        <el-form-item label="文件名">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入文件名"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>

        <el-form-item label="文件类型">
          <el-input
            v-model="searchForm.tag"
            placeholder="请输入文件类型"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            搜索
          </el-button>
          <el-button @click="handleReset">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 文件列表 -->
    <el-card class="list-card" shadow="hover">
      <el-table
        v-loading="loading"
        :data="fileList"
        style="width: 100%"
      >
        <el-table-column prop="ID" label="ID" width="80" />
        <el-table-column prop="name" label="文件名" min-width="200" />
        <el-table-column prop="tag" label="文件类型" width="100" />
        <el-table-column label="预览" width="100">
          <template #default="scope">
            <el-button
              v-if="isImageUrl(scope.row.url)"
              type="primary"
              size="small"
              @click="handlePreview(scope.row)"
            >
              预览
            </el-button>
            <el-button
              v-else
              type="info"
              size="small"
              @click="handlePreview(scope.row)"
            >
              查看
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="上传时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              @click="handleDownload(scope.row.ID, scope.row.name)"
            >
              下载
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(scope.row.ID)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="searchForm.page"
          v-model:page-size="searchForm.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 预览对话框 -->
    <el-dialog
      v-model="previewVisible"
      :title="`文件预览 - ${previewFileName}`"
      width="80%"
      :before-close="closePreview"
      destroy-on-close
    >
      <div class="preview-dialog-content">
        <!-- 图片预览 -->
        <div v-if="previewType === 'image'" class="image-preview-container">
          <img :src="previewUrl" :alt="previewFileName" class="preview-dialog-image">
        </div>

        <!-- 视频预览 -->
        <div v-else-if="previewType === 'video'" class="video-preview-container">
          <video :src="previewUrl" controls class="preview-dialog-video">
            您的浏览器不支持视频播放
          </video>
        </div>

        <!-- 音频预览 -->
        <div v-else-if="previewType === 'audio'" class="audio-preview-container">
          <audio :src="previewUrl" controls class="preview-dialog-audio">
            您的浏览器不支持音频播放
          </audio>
        </div>

        <!-- 其他文件类型 -->
        <div v-else class="other-file-preview">
          <div class="file-info-card">
            <el-icon class="file-icon" size="64">
              <Document />
            </el-icon>
            <h3>{{ previewFileName }}</h3>
            <p class="file-url">
              {{ previewUrl }}
            </p>
            <div class="preview-actions">
              <el-button type="primary" @click="openInNewWindow">
                在新窗口中打开
              </el-button>
              <el-button @click="copyToClipboard(previewUrl)">
                复制链接
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closePreview">
            关闭
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
/* 样式已移至 src/styles/index.scss */

/* 文件预览样式 */
.preview-section {
  margin-top: 20px;
}

.preview-section h4 {
  margin: 0 0 15px 0;
  color: #303133;
  font-size: 16px;
  font-weight: 500;
}

.preview-content {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
  padding: 20px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  background-color: #fafafa;
}

.image-preview {
  max-width: 100%;
  text-align: center;
}

.preview-image {
  max-width: 100%;
  max-height: 400px;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.file-info {
  width: 100%;
  max-width: 500px;
}

.upload-container {
  padding: 20px;
}

.upload-card,
.search-card,
.list-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 预览对话框样式 */
.preview-dialog-content {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  max-height: 80vh;
  overflow: auto;
}

.image-preview-container,
.video-preview-container {
  width: 100%;
  text-align: center;
}

.preview-dialog-image {
  max-width: 100%;
  max-height: 70vh;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.preview-dialog-video {
  max-width: 100%;
  max-height: 70vh;
  border-radius: 8px;
}

.audio-preview-container {
  width: 100%;
  text-align: center;
  padding: 20px;
}

.preview-dialog-audio {
  width: 100%;
  max-width: 500px;
}

.other-file-preview {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.file-info-card {
  text-align: center;
  padding: 40px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background-color: #fafafa;
  max-width: 400px;
  width: 100%;
}

.file-icon {
  color: #909399;
  margin-bottom: 16px;
}

.file-info-card h3 {
  margin: 16px 0;
  color: #303133;
  font-size: 18px;
  font-weight: 500;
}

.file-url {
  margin: 16px 0;
  color: #606266;
  font-size: 14px;
  word-break: break-all;
  background-color: #f5f7fa;
  padding: 8px;
  border-radius: 4px;
}

.preview-actions {
  margin-top: 24px;
  display: flex;
  gap: 12px;
  justify-content: center;
}

.dialog-footer {
  text-align: center;
}
</style>
