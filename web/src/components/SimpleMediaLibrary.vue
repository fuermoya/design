<script setup>
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'
import { formatDate } from '~/utils/format'
import request from '~/utils/request'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  type: {
    type: String,
    default: 'all',
  },
  multiple: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:visible', 'select'])

const mediaFiles = ref([])
const loading = ref(false)
const selectedFiles = ref([])

// 分页相关
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

// 筛选相关
const filterType = ref('JPG')
const searchKeyword = ref('')

// 文件类型选项（按后缀名）
const typeOptions = [
  { label: 'JPG', value: 'jpg' },
  { label: 'PNG', value: 'png' },
  { label: 'GIF', value: 'gif' },
  { label: 'BMP', value: 'bmp' },
  { label: 'WebP', value: 'webp' },
  { label: 'SVG', value: 'svg' },
  { label: 'MP4', value: 'mp4' },
  { label: 'AVI', value: 'avi' },
  { label: 'MOV', value: 'mov' },
  { label: 'WMV', value: 'wmv' },
  { label: 'FLV', value: 'flv' },
  { label: 'WebM', value: 'webm' },
  { label: 'MKV', value: 'mkv' },
  { label: 'M4V', value: 'm4v' },
]
const baseUrl = import.meta.env.VITE_BASE_API || '/api'
// 加载媒体文件列表
async function loadMediaFiles() {
  loading.value = true
  try {
    // 构建请求参数
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
    }

    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }

    const res = await request.get('/fileUploadAndDownload/getFileList', params)
    if (res.code === 0) {
      // 转换文件数据格式以适配媒体库组件，只保留图片和视频
      const allFiles = res.data.list
        .map(file => ({
          ID: file.ID,
          name: file.name,
          url: `${baseUrl}/fileUploadAndDownload/preview?id=${file.ID}`,
          type: getFileType(file.name),
          CreatedAt: file.CreatedAt || file.UpdatedAt || new Date().toISOString(),
        }))

      // 应用筛选条件
      let filteredFiles = allFiles

      // 按文件后缀筛选
      if (filterType.value !== 'all') {
        filteredFiles = filteredFiles.filter((file) => {
          const fileName = file.name.toLowerCase()
          return fileName.endsWith(filterType.value.toLowerCase())
        })
      }

      // 按关键词搜索
      if (searchKeyword.value) {
        const keyword = searchKeyword.value.toLowerCase()
        filteredFiles = filteredFiles.filter(file =>
          file.name.toLowerCase().includes(keyword),
        )
      }

      // 使用后端返回的数据
      mediaFiles.value = filteredFiles
      total.value = res.data.total || filteredFiles.length
    }
  }
  catch (error) {
    ElMessage.error('加载媒体文件失败')
    console.error('加载媒体文件失败:', error)
  }
  finally {
    loading.value = false
  }
}

// 处理分页变化
function handlePageChange(page) {
  currentPage.value = page
  loadMediaFiles()
}

// 处理页面大小变化
function handleSizeChange(size) {
  pageSize.value = size
  currentPage.value = 1
  loadMediaFiles()
}

// 处理筛选变化
function handleFilterChange() {
  currentPage.value = 1
  loadMediaFiles()
}

// 重置筛选
function resetFilter() {
  filterType.value = 'jpg' // 根据用户设置的默认值
  searchKeyword.value = ''
  currentPage.value = 1
  loadMediaFiles()
}

// 根据文件名判断文件类型
function getFileType(fileName) {
  const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp', '.svg']
  const videoExtensions = ['.mp4', '.avi', '.mov', '.wmv', '.flv', '.webm', '.mkv', '.m4v']

  const extension = fileName.toLowerCase().substring(fileName.lastIndexOf('.'))

  if (imageExtensions.includes(extension)) {
    return 'image'
  }
  else if (videoExtensions.includes(extension)) {
    return 'video'
  }
  else {
    return 'other'
  }
}

// 选择文件
function selectFile(file) {
  if (props.multiple) {
    const index = selectedFiles.value.findIndex(f => f.ID === file.ID)
    if (index > -1) {
      selectedFiles.value.splice(index, 1)
    }
    else {
      selectedFiles.value.push(file)
    }
  }
  else {
    selectedFiles.value = [file]
  }
}

// 确认选择
function confirmSelection() {
  if (selectedFiles.value.length === 0) {
    ElMessage.warning('请选择文件')
    return
  }
  emit('select', selectedFiles.value)
  emit('update:visible', false)
  selectedFiles.value = []
}

// 取消选择
function cancelSelection() {
  emit('update:visible', false)
  selectedFiles.value = []
}

// 处理图片加载错误
function handleImageError(event) {
  const img = event.target
  img.style.display = 'none'
  // 隐藏加载动画
  const preview = img.closest('.media-preview')
  if (preview) {
    preview.classList.add('loaded')
  }
  // 可以在这里添加默认图片
  const container = img.parentElement
  container.innerHTML = '<div class="image-error">图片加载失败</div>'
}

// 处理图片加载成功
function handleImageLoad(event) {
  const img = event.target
  img.classList.add('loaded')
  // 隐藏加载动画
  const preview = img.closest('.media-preview')
  if (preview) {
    preview.classList.add('loaded')
  }
}

onMounted(() => {
  loadMediaFiles()
})
</script>

<template>
  <el-dialog
    :model-value="visible"
    title="媒体库"
    width="80%"
    :close-on-click-modal="false"
    @update:model-value="$emit('update:visible', $event)"
    @close="cancelSelection"
  >
    <div class="media-library">
      <!-- 筛选工具栏 -->
      <div class="filter-toolbar">
        <div class="filter-left">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文件名..."
            clearable
            style="width: 200px; margin-right: 12px;"
            @keyup.enter="handleFilterChange"
            @clear="handleFilterChange"
          />
          <el-select
            v-model="filterType"
            placeholder="选择文件格式"
            style="width: 140px; margin-right: 12px;"
            @change="handleFilterChange"
          >
            <el-option
              v-for="option in typeOptions"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>
          <el-button @click="resetFilter">
            重置
          </el-button>
        </div>
        <div class="filter-right">
          <span class="file-count">
            共 {{ total }} 个文件
          </span>
        </div>
      </div>

      <!-- 媒体文件网格 -->
      <div class="media-grid">
        <div
          v-for="file in mediaFiles"
          :key="file.ID"
          class="media-item"
          :class="{
            selected: selectedFiles.some(f => f.ID === file.ID),
            video: file.type === 'video',
          }"
          :data-type="file.type.toUpperCase()"
          @click="selectFile(file)"
        >
          <div class="media-preview">
            <img
              v-if="file.type === 'image'"
              :src="file.url"
              :alt="file.name"
              class="preview-image"
              @error="handleImageError"
              @load="handleImageLoad"
            >
            <video
              v-else-if="file.type === 'video'"
              :src="file.url"
              class="preview-video"
              preload="metadata"
            />
            <div class="file-info">
              <div class="file-name">
                {{ file.name }}
              </div>
              <div class="file-meta">
                {{ formatDate(file.CreatedAt) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[12, 24, 48, 96]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <template #footer>
      <el-button @click="cancelSelection">
        取消
      </el-button>
      <el-button type="primary" @click="confirmSelection">
        确定选择
      </el-button>
    </template>
  </el-dialog>
</template>

<style scoped>
.media-library {
  max-height: 60vh;
  overflow-y: auto;
}

/* 筛选工具栏 */
.filter-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #e4e7ed;
  background-color: #fafafa;
}

.filter-left {
  display: flex;
  align-items: center;
}

.filter-right {
  display: flex;
  align-items: center;
}

.file-count {
  color: #606266;
  font-size: 14px;
}

/* 媒体文件网格 */
.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 20px;
  padding: 20px;
  min-height: 300px;
}

.media-item {
  border: 2px solid #e4e7ed;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  background: white;
  position: relative;
}

.media-item:hover {
  border-color: #409eff;
  transform: translateY(-4px) scale(1.02);
  box-shadow: 0 8px 25px rgba(64, 158, 255, 0.15);
}

.media-item.selected {
  border-color: #409eff;
  background: linear-gradient(135deg, #f0f9ff 0%, #e6f7ff 100%);
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.2);
}

.media-item.selected::before {
  content: '✓';
  position: absolute;
  top: 8px;
  right: 8px;
  width: 24px;
  height: 24px;
  background: #409eff;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: bold;
  z-index: 10;
}

.media-preview {
  position: relative;
  aspect-ratio: 16/9;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  overflow: hidden;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition:
    transform 0.3s ease,
    opacity 0.3s ease;
  opacity: 0;
}

.preview-image.loaded {
  opacity: 1;
}

.preview-image:hover {
  transform: scale(1.05);
}

/* 图片加载错误样式 */
.image-error {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
  color: #6c757d;
  font-size: 12px;
  font-weight: 500;
}

/* 图片加载中样式 */
.media-preview::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 30px;
  height: 30px;
  border: 2px solid #e4e7ed;
  border-top: 2px solid #409eff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  opacity: 0.6;
  transition: opacity 0.3s ease;
}

/* 图片加载完成后隐藏加载动画 */
.media-preview.loaded::before {
  opacity: 0;
  pointer-events: none;
}

@keyframes spin {
  0% {
    transform: translate(-50%, -50%) rotate(0deg);
  }
  100% {
    transform: translate(-50%, -50%) rotate(360deg);
  }
}

.preview-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.preview-video:hover {
  transform: scale(1.05);
}

/* 视频播放图标 - 默认隐藏 */
.media-preview::after {
  content: '▶';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 40px;
  height: 40px;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  opacity: 0;
  transition: opacity 0.3s ease;
  pointer-events: none;
}

/* 只对视频显示播放图标 */
.media-item.video .media-preview::after {
  opacity: 0.8;
}

.media-item.video:hover .media-preview::after {
  opacity: 1;
}

.file-info {
  padding: 12px;
  background: linear-gradient(180deg, transparent 0%, rgba(0, 0, 0, 0.8) 30%, rgba(0, 0, 0, 0.9) 100%);
  color: white;
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  backdrop-filter: blur(4px);
}

.file-name {
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
}

.file-meta {
  font-size: 11px;
  opacity: 0.9;
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-meta::before {
  content: '📄';
  font-size: 10px;
}

/* 文件类型标签 */
.media-item::after {
  content: attr(data-type);
  position: absolute;
  top: 8px;
  left: 8px;
  padding: 2px 6px;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 500;
  text-transform: uppercase;
  backdrop-filter: blur(4px);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .media-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 12px;
    padding: 12px;
  }

  .media-item {
    border-radius: 8px;
  }

  .file-info {
    padding: 8px;
  }

  .file-name {
    font-size: 11px;
  }

  .file-meta {
    font-size: 10px;
  }
}

@media (max-width: 480px) {
  .media-grid {
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 8px;
    padding: 8px;
  }

  .filter-toolbar {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }

  .filter-left {
    flex-wrap: wrap;
    gap: 8px;
  }

  .filter-left .el-input,
  .filter-left .el-select {
    width: 100% !important;
    margin-right: 0 !important;
  }
}

/* 深色模式支持 */
@media (prefers-color-scheme: dark) {
  .media-item {
    background: #2c2c2c;
    border-color: #4a4a4a;
  }

  .media-item:hover {
    border-color: #409eff;
    box-shadow: 0 8px 25px rgba(64, 158, 255, 0.25);
  }

  .media-item.selected {
    background: linear-gradient(135deg, #1a3a5f 0%, #1e4a7a 100%);
  }

  .media-preview {
    background: linear-gradient(135deg, #3a3a3a 0%, #2a2a2a 100%);
  }

  .image-error {
    background: linear-gradient(135deg, #3a3a3a 0%, #2a2a2a 100%);
    color: #b0b0b0;
  }
}
</style>
