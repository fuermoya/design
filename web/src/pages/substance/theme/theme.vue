<script setup lang="ts">
import type { SysTheme, ThemeConfig } from '~/types/portal'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'
import request from '~/utils/request'
import { formatDate } from '~/utils/format'
import { applyThemeConfig, defaultThemeConfig, parseThemeConfig } from '~/utils/theme'

// 响应式数据
const themes = ref<SysTheme[]>([])
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
  config: '',
  preview: '',
  sort: 0,
})

// 主题配置表单
const configForm = ref<ThemeConfig>({ ...defaultThemeConfig })

// 加载主题列表
async function loadThemes() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
    }
    const res: any = await request.get('/sysTheme/getSysThemeList', params)
    if (res.code === 0) {
      themes.value = res.data.list.map((theme: SysTheme) => ({
        ...theme,
        selected: false,
      }))
      total.value = res.data.total
    }
  }
  catch (error) {
    ElMessage.error('加载主题列表失败')
    console.error('加载主题列表失败:', error)
  }
  finally {
    loading.value = false
  }
}

// 打开新增对话框
function openAddDialog() {
  dialogTitle.value = '新增主题'
  isEdit.value = false
  formData.value = {
    id: 0,
    name: '',
    description: '',
    config: JSON.stringify(defaultThemeConfig, null, 2),
    preview: '',
    sort: 0,
  }
  configForm.value = { ...defaultThemeConfig }
  dialogVisible.value = true
}

// 打开编辑对话框
function openEditDialog(row: SysTheme) {
  dialogTitle.value = '编辑主题'
  isEdit.value = true
  formData.value = {
    id: row.ID,
    name: row.name,
    description: row.description,
    config: row.config,
    preview: row.preview || '',
    sort: row.sort,
  }

  // 解析配置
  if (row.config) {
    configForm.value = parseThemeConfig(row.config)
  }

  dialogVisible.value = true
}

// 保存主题
async function saveTheme() {
  try {
    // 验证必填字段
    if (!formData.value.name.trim()) {
      ElMessage.error('请输入主题名称')
      return
    }

    // 更新配置JSON
    formData.value.config = JSON.stringify(configForm.value, null, 2)

    if (isEdit.value) {
      await request.put('/sysTheme/updateSysTheme', formData.value)
      ElMessage.success('更新成功')
    }
    else {
      await request.post('/sysTheme/createSysTheme', formData.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadThemes()
  }
  catch (error) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    console.error('保存主题失败:', error)
  }
}

// 删除主题
async function deleteTheme(id: number) {
  try {
    await ElMessageBox.confirm('确定要删除这个主题吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await request.del('/sysTheme/deleteSysTheme', {
      params: { ID: id },
    })

    loadThemes()
  }
  catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error('删除主题失败:', error)
    }
  }
}

// 激活主题
async function activateTheme(id: number) {
  try {
    await ElMessageBox.confirm('确定要激活这个主题吗？激活后门户网站将立即应用新主题。', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await request.post('/sysTheme/activateTheme', { id })
    ElMessage.success('激活成功，门户网站已应用新主题')
    loadThemes()
  }
  catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('激活失败')
      console.error('激活主题失败:', error)
    }
  }
}

// 预览主题配置
function previewTheme() {
  try {
    applyThemeConfig(configForm.value)
    ElMessage.success('主题预览已应用，请查看门户网站效果')
  }
  catch (error) {
    ElMessage.error('预览失败')
    console.error('预览主题失败:', error)
  }
}

// 重置主题配置
function resetThemeConfig() {
  configForm.value = { ...defaultThemeConfig }
  ElMessage.success('主题配置已重置为默认值')
}

// 批量删除
async function batchDelete() {
  const selectedIds = themes.value
    .filter((theme: SysTheme) => theme.selected)
    .map((theme: SysTheme) => theme.ID)

  if (selectedIds.length === 0) {
    ElMessage.warning('请选择要删除的主题')
    return
  }

  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedIds.length} 个主题吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await request.del('/sysTheme/deleteSysThemeByIds', {
      params: { 'IDs[]': selectedIds },
    })

    loadThemes()
  }
  catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
      console.error('批量删除失败:', error)
    }
  }
}

// 处理表格选择变化
function handleSelectionChange(selection: SysTheme[]) {
  themes.value.forEach((item: SysTheme) => {
    item.selected = selection.includes(item)
  })
}

// 页面初始化
onMounted(() => {
  loadThemes()
})
</script>

<template>
  <div class="theme-management">
    <!-- 主题列表 -->
    <div class="table-container">
      <div class="table-header">
        <div class="table-actions">
          <el-button type="primary" @click="openAddDialog">
            新增主题
          </el-button>
          <el-button type="danger" @click="batchDelete">
            批量删除
          </el-button>
        </div>
      </div>
      <el-table
        v-loading="loading"
        :data="themes"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="主题名称" min-width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="isActive" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.isActive ? 'success' : 'info'">
              {{ row.isActive ? '已激活' : '未激活' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="CreatedAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button type="primary" size="small" @click="openEditDialog(row)">
                编辑
              </el-button>
              <el-button
                v-if="!row.isActive"
                type="success"
                size="small"
                @click="activateTheme(row.ID)"
              >
                激活
              </el-button>
              <el-button type="danger" size="small" @click="deleteTheme(row.ID)">
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
          @size-change="loadThemes"
          @current-change="loadThemes"
        />
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="90%"
      :close-on-click-modal="false"
    >
      <el-form :model="formData" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="主题名称" required>
              <el-input v-model="formData.name" placeholder="请输入主题名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序">
              <el-input-number v-model="formData.sort" :min="0" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="主题描述">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入主题描述"
          />
        </el-form-item>

        <el-form-item label="预览图片">
          <el-input v-model="formData.preview" placeholder="请输入预览图片URL" />
        </el-form-item>

        <!-- 主题配置 -->
        <el-form-item label="主题配置">
          <div class="config-actions">
            <el-button type="primary" @click="previewTheme">
              预览主题
            </el-button>
            <el-button @click="resetThemeConfig">
              重置配置
            </el-button>
          </div>
          <el-tabs type="border-card">
            <el-tab-pane label="基础配置">
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-form-item label="主色调">
                    <el-color-picker v-model="configForm.primaryColor" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="次要色调">
                    <el-color-picker v-model="configForm.secondaryColor" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="背景色">
                    <el-color-picker v-model="configForm.backgroundColor" />
                  </el-form-item>
                </el-col>
              </el-row>

              <el-row :gutter="20">
                <el-col :span="8">
                  <el-form-item label="文字颜色">
                    <el-color-picker v-model="configForm.textColor" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="字体">
                    <el-input v-model="configForm.fontFamily" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="字体大小">
                    <el-input v-model="configForm.fontSize" />
                  </el-form-item>
                </el-col>
              </el-row>

              <el-row :gutter="20">
                <el-col :span="8">
                  <el-form-item label="圆角">
                    <el-input v-model="configForm.borderRadius" />
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="阴影">
                    <el-input v-model="configForm.shadow" />
                  </el-form-item>
                </el-col>
              </el-row>
            </el-tab-pane>

            <el-tab-pane label="网站信息">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="网站名称">
                    <el-input v-model="configForm.siteName" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="网站描述">
                    <el-input v-model="configForm.siteDescription" />
                  </el-form-item>
                </el-col>
              </el-row>

              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="Logo">
                    <el-input v-model="configForm.logo" placeholder="请输入Logo URL" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="网站图标">
                    <el-input v-model="configForm.favicon" placeholder="请输入网站图标URL" />
                  </el-form-item>
                </el-col>
              </el-row>

              <el-form-item label="页脚文字">
                <el-input v-model="configForm.footerText" />
              </el-form-item>
            </el-tab-pane>

            <el-tab-pane label="JSON配置">
              <el-form-item>
                <el-input
                  v-model="formData.config"
                  type="textarea"
                  :rows="15"
                  placeholder="请输入JSON配置"
                />
              </el-form-item>
            </el-tab-pane>
          </el-tabs>
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="button-group">
          <el-button @click="dialogVisible = false">
            取消
          </el-button>
          <el-button type="primary" @click="saveTheme">
            保存
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.theme-management {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: calc(100vh - 120px);
}

.config-actions {
  margin-bottom: 16px;
  display: flex;
  gap: 12px;
}

@media (max-width: 768px) {
  .theme-management {
    padding: 10px;
  }
}
</style>
