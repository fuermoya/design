<script setup>
import { ElMessage } from 'element-plus'
import { ref } from 'vue'
import { formatDate } from '~/utils/format'
import request from '~/utils/request'

defineOptions({
  name: 'SysOperationRecord',
})

// API方法
function getSysOperationRecordList(params) {
  return request.get('/sysOperationRecord/getSysOperationRecordList', params)
}
function deleteSysOperationRecord(data) {
  return request.del('/sysOperationRecord/deleteSysOperationRecord', data)
}
function deleteSysOperationRecordByIds(data) {
  return request.del('/sysOperationRecord/deleteSysOperationRecordByIds', data)
}

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const loading = ref(false)

function onReset() {
  searchInfo.value = {}
}

// 条件搜索
function onSubmit() {
  page.value = 1
  pageSize.value = 10
  if (searchInfo.value.status === '') {
    searchInfo.value.status = null
  }
  getTableData()
}

// 分页
function handleSizeChange(val) {
  pageSize.value = val
  getTableData()
}

function handleCurrentChange(val) {
  page.value = val
  getTableData()
}

// 查询
async function getTableData() {
  loading.value = true
  try {
    const table = await getSysOperationRecordList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value,
    })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }
  catch (error) {
    console.error('获取数据失败:', error)
    ElMessage.error('获取数据失败')
  }
  finally {
    loading.value = false
  }
}

getTableData()

const deleteVisible = ref(false)
const multipleSelection = ref([])

function handleSelectionChange(val) {
  multipleSelection.value = val
}

async function onDelete() {
  const ids = []
  multipleSelection.value && multipleSelection.value.forEach((item) => {
    ids.push(item.ID)
  })
  try {
    const res = await deleteSysOperationRecordByIds({ ids })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功',
      })
      if (tableData.value.length === ids.length && page.value > 1) {
        page.value--
      }
      deleteVisible.value = false
      getTableData()
    }
  }
  catch (error) {
    console.error('删除失败:', error)
    ElMessage.error('删除失败')
  }
}

async function deleteSysOperationRecordFunc(row) {
  row.visible = false
  try {
    const res = await deleteSysOperationRecord({ ID: row.ID })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功',
      })
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  }
  catch (error) {
    console.error('删除失败:', error)
    ElMessage.error('删除失败')
  }
}

function fmtBody(value) {
  try {
    return JSON.parse(value)
  }
  catch {
    return value
  }
}

// 格式化状态码显示
function formatStatus(status) {
  if (status >= 200 && status < 300) {
    return { type: 'success', text: status }
  }
  if (status >= 400 && status < 500) {
    return { type: 'warning', text: status }
  }
  if (status >= 500) {
    return { type: 'danger', text: status }
  }
  return { type: 'info', text: status }
}

// 格式化延迟时间
function formatLatency(latency) {
  if (!latency)
    return '-'
  const ms = Number.parseInt(latency)
  if (ms < 100) {
    return `${ms}ms`
  }
  if (ms < 1000) {
    return `${(ms / 1000).toFixed(1)}s`
  }
  return `${(ms / 1000).toFixed(2)}s`
}
</script>

<template>
  <div>
    <div class="table-search">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="请求方法">
          <el-select v-model="searchInfo.method" placeholder="请选择请求方法" clearable>
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
            <el-option label="PATCH" value="PATCH" />
          </el-select>
        </el-form-item>
        <el-form-item label="请求路径">
          <el-input v-model="searchInfo.path" placeholder="请输入请求路径" />
        </el-form-item>
        <el-form-item label="状态码">
          <el-input v-model="searchInfo.status" placeholder="请输入状态码" />
        </el-form-item>
        <el-form-item label="IP地址">
          <el-input v-model="searchInfo.ip" placeholder="请输入IP地址" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" :loading="loading" @click="onSubmit">
            查询
          </el-button>
          <el-button icon="Refresh" @click="onReset">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="table-container">
      <div class="table-header">
        <div class="table-actions">
          <el-button icon="Delete" type="danger" :disabled="!multipleSelection.length" @click="deleteVisible = true">
            批量删除 ({{ multipleSelection.length }})
          </el-button>
        </div>
      </div>

      <el-table
        v-loading="loading" :data="tableData" style="width: 100%" tooltip-effect="dark" row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column align="left" type="selection" width="55" />
        <el-table-column align="left" label="操作人" min-width="140">
          <template #default="scope">
            <div v-if="scope.row.user">
              {{ scope.row.user.userName }}({{ scope.row.user.nickName }})
            </div>
            <div v-else>
              -
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作时间" min-width="180">
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态码" min-width="100">
          <template #default="scope">
            <el-tag :type="formatStatus(scope.row.status).type">
              {{ formatStatus(scope.row.status).text }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="延迟" min-width="100">
          <template #default="scope">
            {{ formatLatency(scope.row.latency) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="请求IP" prop="ip" min-width="120" />
        <el-table-column align="left" label="请求方法" min-width="100">
          <template #default="scope">
            <el-tag
              :type="scope.row.method === 'GET' ? 'success' : scope.row.method === 'POST' ? 'primary' : scope.row.method === 'PUT' ? 'warning' : 'danger'"
            >
              {{ scope.row.method }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="请求路径" prop="path" min-width="240" show-overflow-tooltip />
        <el-table-column align="left" label="请求体" min-width="80">
          <template #default="scope">
            <div>
              <el-popover v-if="scope.row.body" placement="left-start" trigger="click" width="500">
                <div class="popover-box">
                  <pre>{{ fmtBody(scope.row.body) }}</pre>
                </div>
                <template #reference>
                  <el-button type="primary" link>
                    查看
                  </el-button>
                </template>
              </el-popover>
              <span v-else>-</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="响应体" min-width="80">
          <template #default="scope">
            <div>
              <el-popover v-if="scope.row.resp" placement="left-start" trigger="click" width="500">
                <div class="popover-box">
                  <pre>{{ fmtBody(scope.row.resp) }}</pre>
                </div>
                <template #reference>
                  <el-button type="primary" link>
                    查看
                  </el-button>
                </template>
              </el-popover>
              <span v-else>-</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="错误信息" min-width="120" show-overflow-tooltip>
          <template #default="scope">
            <span v-if="scope.row.error_message" class="error-message">
              {{ scope.row.error_message }}
            </span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="100" fixed="right" class-name="table-actions-column">
          <template #default="scope">
            <el-popover v-model="scope.row.visible" placement="top" width="160">
              <p>确定要删除这条操作记录吗？</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-button type="primary" link @click="scope.row.visible = false">
                  取消
                </el-button>
                <el-button type="danger" @click="deleteSysOperationRecordFunc(scope.row)">
                  确定
                </el-button>
              </div>
              <template #reference>
                <el-button icon="Delete" type="danger" link @click="scope.row.visible = true">
                  删除
                </el-button>
              </template>
            </el-popover>
          </template>
        </el-table-column>
      </el-table>

      <div class="table-pagination">
        <el-pagination
          :current-page="page" :page-size="pageSize" :page-sizes="[10, 30, 50, 100, 200, 500]"
          :total="total" layout="total, sizes, prev, pager, next, jumper" @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 批量删除确认对话框 -->
    <el-dialog v-model="deleteVisible" title="确认删除" width="400px">
      <p>确定要删除选中的 {{ multipleSelection.length }} 条操作记录吗？</p>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="deleteVisible = false">取消</el-button>
          <el-button type="danger" @click="onDelete">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style lang="scss" scoped>
:deep(.el-form-item) {
  margin-bottom: 16px;
}
</style>
