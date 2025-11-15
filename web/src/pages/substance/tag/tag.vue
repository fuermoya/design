<script setup lang="ts">
import type { SysTag } from '~/types/portal'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'
import { formatDate } from '~/utils/format'
import request from '~/utils/request'

// 响应式数据
const tags = ref<(SysTag & { selected?: boolean })[]>([])
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
  color: '#409EFF',
  status: 1,
})

// 状态选项
const statusOptions = [
  { label: '启用', value: 1 },
  { label: '禁用', value: 2 },
]

// 颜色选项
const colorOptions = [
  '#409EFF',
  '#67C23A',
  '#E6A23C',
  '#F56C6C',
  '#909399',
  '#9C27B0',
  '#FF9800',
  '#795548',
]

// 加载标签列表
async function loadTags() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
    }
    const res: any = await request.get('/sysTag/getSysTagList', params)
    if (res.code === 0) {
      tags.value = res.data.list.map((tag: SysTag) => ({ ...tag, selected: false }))
      total.value = res.data.total
    }
  }
  catch (error) {
    ElMessage.error('加载标签列表失败')
    console.error('加载标签列表失败:', error)
  }
  finally {
    loading.value = false
  }
}

// 打开新增对话框
function openAddDialog() {
  dialogTitle.value = '新增标签'
  isEdit.value = false
  formData.value = {
    id: 0,
    name: '',
    color: '#409EFF',
    status: 1,
  }
  dialogVisible.value = true
}

// 打开编辑对话框
function openEditDialog(row: SysTag) {
  dialogTitle.value = '编辑标签'
  isEdit.value = true
  formData.value = {
    id: row.ID,
    name: row.name,
    color: row.color,
    status: row.status,
  }
  dialogVisible.value = true
}

// 保存标签
async function saveTag() {
  try {
    if (isEdit.value) {
      await request.put('/sysTag/updateSysTag', formData.value)
      ElMessage.success('更新成功')
    }
    else {
      await request.post('/sysTag/createSysTag', formData.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadTags()
  }
  catch (error) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    console.error('保存标签失败:', error)
  }
}

// 删除标签
async function deleteTag(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除这个标签吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await request.del(`/sysTag/deleteSysTag?id=${id}`, {})

    loadTags()
  }
  catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error('删除标签失败:', error)
    }
  }
}

// 批量删除
async function batchDelete() {
  const selectedIds = tags.value
    .filter(tag => tag.selected)
    .map(tag => tag.ID)

  if (selectedIds.length === 0) {
    ElMessage.warning('请选择要删除的标签')
    return
  }

  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedIds.length} 个标签吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await request.del('/sysTag/deleteSysTagByIds', {
      ids: selectedIds,
    })

    loadTags()
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

// 页面初始化
onMounted(() => {
  loadTags()
})
</script>

<template>
  <div class="tag-container">
    <div class="table-container">
      <div class="table-header">
        <div class="table-actions">
          <el-button type="primary" @click="openAddDialog">
            新增标签
          </el-button>
          <el-button type="danger" @click="batchDelete">
            批量删除
          </el-button>
        </div>
      </div>

      <!-- 标签列表 -->
      <el-table
        v-loading="loading"
        :data="tags"
        style="width: 100%"
        @selection-change="(selection) => tags.forEach(item => item.selected = selection.includes(item))"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="标签名称" min-width="150" />
        <el-table-column prop="color" label="颜色" width="120">
          <template #default="{ row }">
            <div class="color-preview">
              <div
                class="color-block"
                :style="{ backgroundColor: row.color }"
              />
              <span>{{ row.color }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
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
              <el-button type="danger" size="small" @click="deleteTag(row.ID)">
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="table-pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadTags"
          @current-change="loadTags"
        />
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="formData" label-width="100px">
        <el-form-item label="标签名称" required>
          <el-input v-model="formData.name" placeholder="请输入标签名称" />
        </el-form-item>

        <el-form-item label="标签颜色" required>
          <div class="color-selector">
            <el-color-picker v-model="formData.color" />
            <div class="color-options">
              <div
                v-for="color in colorOptions"
                :key="color"
                class="color-option"
                :style="{ backgroundColor: color }"
                @click="formData.color = color"
              />
            </div>
          </div>
        </el-form-item>

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
      </el-form>

      <template #footer>
        <div class="button-group">
          <el-button @click="dialogVisible = false">
            取消
          </el-button>
          <el-button type="primary" @click="saveTag">
            保存
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.tag-container {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: calc(100vh - 120px);

  @media (max-width: 768px) {
    padding: 10px;
  }
}

.color-preview {
  display: flex;
  align-items: center;
  gap: 8px;
}

.color-block {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  border: 1px solid #ddd;
}

.color-selector {
  display: flex;
  align-items: center;
  gap: 16px;
}

.color-options {
  display: flex;
  gap: 8px;
}

.color-option {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: border-color 0.2s;
}

.color-option:hover {
  border-color: #409eff;
}
</style>
