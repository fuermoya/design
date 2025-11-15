<script setup lang="ts">
import type { SysCategory } from '~/types/portal'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'
import { formatDate } from '~/utils/format'
import request from '~/utils/request'

// 扩展 SysCategory 类型以包含 selected 属性
interface CategoryWithSelection extends SysCategory {
  selected?: boolean
}

// 响应式数据
const categories = ref<CategoryWithSelection[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 表单数据
const formData = ref({
  id: 0,
  name: '',
  description: '',
  sort: 0,
  status: 1,
})

// 状态选项
const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 2 },
]

// 加载分类列表
async function loadCategories() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
    }
    const res: any = await request.get('/sysCategory/getSysCategoryList', params)
    if (res.code === 0) {
      categories.value = res.data.list
      total.value = res.data.total
    }
  }
  catch (error) {
    ElMessage.error('加载分类列表失败')
    console.error('加载分类列表失败:', error)
  }
  finally {
    loading.value = false
  }
}

// 打开新增对话框
function openAddDialog() {
  dialogTitle.value = '新增分类'
  isEdit.value = false
  formData.value = {
    id: 0,
    name: '',
    description: '',
    sort: 0,
    status: 1,
  }
  dialogVisible.value = true
}

// 打开编辑对话框
function openEditDialog(row: SysCategory) {
  dialogTitle.value = '编辑分类'
  isEdit.value = true
  formData.value = {
    id: row.ID,
    name: row.name,
    description: row.description,
    sort: row.sort,
    status: row.status,
  }
  dialogVisible.value = true
}

// 保存分类
async function saveCategory() {
  try {
    if (isEdit.value) {
      await request.put('/sysCategory/updateSysCategory', formData.value)
      ElMessage.success('更新成功')
    }
    else {
      await request.post('/sysCategory/createSysCategory', formData.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadCategories()
  }
  catch (error) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    console.error('保存分类失败:', error)
  }
}

// 删除分类
async function deleteCategory(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除这个分类吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await request.del(`/sysCategory/deleteSysCategory?id=${id}`, {})

    loadCategories()
  }
  catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error('删除分类失败:', error)
    }
  }
}

// 批量删除
async function batchDelete() {
  const selectedIds = categories.value
    .filter((category: CategoryWithSelection) => category.selected)
    .map((category: CategoryWithSelection) => category.ID)

  if (selectedIds.length === 0) {
    ElMessage.warning('请选择要删除的分类')
    return
  }

  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedIds.length} 个分类吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await request.del('/sysCategory/deleteSysCategoryByIds', {
      ids: selectedIds,
    })

    loadCategories()
  }
  catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
      console.error('批量删除失败:', error)
    }
  }
}

// 获取状态文本
function getStatusText(status: number) {
  const statusMap: Record<number, string> = {
    1: '启用',
    2: '禁用',
  }
  return statusMap[status] || '未知'
}

// 获取状态类型
function getStatusType(status: number) {
  const typeMap: Record<number, 'success' | 'danger'> = {
    1: 'success',
    2: 'danger',
  }
  return typeMap[status] || 'info'
}

// 处理表格选择变化
function handleSelectionChange(selection: CategoryWithSelection[]) {
  categories.value.forEach((item: CategoryWithSelection) => item.selected = selection.includes(item))
}

// 页面初始化
onMounted(() => {
  loadCategories()
})
</script>

<template>
  <div class="category-container">
    <!-- 分类列表 -->
    <div class="table-container">
      <div class="table-header">
        <div class="table-actions">
          <el-button type="primary" @click="openAddDialog">
            新增分类
          </el-button>
          <el-button type="danger" @click="batchDelete">
            批量删除
          </el-button>
        </div>
      </div>
      <el-table
        v-loading="loading"
        :data="categories"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="分类名称" min-width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" class="status-tag">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="CreatedAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right" class-name="table-actions-column">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button type="primary" size="small" @click="openEditDialog(row)">
                编辑
              </el-button>
              <el-button type="danger" size="small" @click="deleteCategory(row.ID)">
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
          @size-change="loadCategories"
          @current-change="loadCategories"
        />
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      :close-on-click-modal="false"
      class="category-dialog"
    >
      <el-form :model="formData" label-width="100px">
        <el-form-item label="分类名称" required>
          <el-input v-model="formData.name" placeholder="请输入分类名称" />
        </el-form-item>

        <el-form-item label="分类描述">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入分类描述"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="排序">
              <el-input-number v-model="formData.sort" :min="0" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
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
        </el-row>
      </el-form>

      <template #footer>
        <div class="button-group">
          <el-button @click="dialogVisible = false">
            取消
          </el-button>
          <el-button type="primary" @click="saveCategory">
            保存
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.category-container {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: calc(100vh - 120px);
}

/* 卡片样式 */
.upload-card,
.list-card {
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  border: none;
  background: #fff;
  margin-bottom: 24px;
  transition: all 0.3s ease;
}

.upload-card:hover,
.list-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.upload-card :deep(.el-card__body),
.list-card :deep(.el-card__body) {
  padding: 24px;
}

/* 表格操作区域样式 */
.table-actions {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

/* 对话框样式 */
.category-dialog :deep(.el-dialog) {
  border-radius: 8px;
  overflow: hidden;
}

.category-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px;
}

.category-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.category-dialog :deep(.el-dialog__body) {
  padding: 24px;
}

.category-dialog :deep(.el-dialog__footer) {
  padding: 16px 24px;
  border-top: 1px solid #ebeef5;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .category-container {
    padding: 10px;
  }

  .upload-card,
  .list-card {
    margin-bottom: 16px;
  }

  .upload-card :deep(.el-card__body),
  .list-card :deep(.el-card__body) {
    padding: 16px;
  }

  .table-actions {
    gap: 8px;
    margin-bottom: 16px;
  }
}

/* 暗色主题样式 */
html.dark .category-container {
  background-color: #1d1e1f;
}

html.dark .upload-card,
html.dark .list-card {
  background: #2b2b2c;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
}

html.dark .upload-card:hover,
html.dark .list-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}
</style>
