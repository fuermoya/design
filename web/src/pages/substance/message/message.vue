<script setup lang="ts">
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, ref } from 'vue'
import { formatDate } from '~/utils/format'
import request from '~/utils/request'

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const messageList = ref([])
const total = ref(0)
const stats = ref({
  total: 0,
  unread: 0,
  replied: 0,
})

const pageInfo = ref({
  page: 1,
  pageSize: 20,
})

const searchForm = ref({
  name: '',
  status: 0,
})

// 对话框相关
const messageDialogVisible = ref(false)
const replyDialogVisible = ref(false)
const currentMessage = ref(null)

// 回复表单
const replyFormRef = ref()
const replyForm = ref({
  id: 0,
  reply: '',
})

const replyRules = {
  reply: [
    { required: true, message: '请输入回复内容', trigger: 'blur' },
    { min: 5, max: 500, message: '回复内容长度在 5 到 500 个字符', trigger: 'blur' },
  ],
}

// 加载留言列表
async function loadMessages() {
  try {
    loading.value = true
    const params = {
      ...pageInfo.value,
      ...searchForm.value,
    }
    const res: any = await request.get('/sysMessage/getMessageList', params)
    if (res.code === 0) {
      messageList.value = res.data.list
      total.value = res.data.total
    }
  }
  catch (error) {
    ElMessage.error('加载留言列表失败')
    console.error('加载留言列表失败:', error)
  }
  finally {
    loading.value = false
  }
}

// 加载统计数据
async function loadStats() {
  try {
    const res: any = await request.get('/sysMessage/getMessageStats')
    if (res.code === 0) {
      stats.value = res.data
    }
  }
  catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 查看留言详情
async function viewMessage(message: any) {
  try {
    const res: any = await request.get(`/sysMessage/getMessageById?id=${message.ID}`)
    if (res.code === 0) {
      currentMessage.value = res.data
      messageDialogVisible.value = true
    }
  }
  catch (error) {
    ElMessage.error('获取留言详情失败')
    console.error('获取留言详情失败:', error)
  }
}

// 标记为已读
async function markAsRead(message: any) {
  try {
    const res: any = await request.post('/sysMessage/markAsRead', { id: message.ID })
    if (res.code === 0) {
      ElMessage.success('标记成功')
      loadMessages()
      loadStats()
    }
  }
  catch (error) {
    ElMessage.error('标记失败')
    console.error('标记失败:', error)
  }
}

// 回复留言
function replyMessage(message: any) {
  currentMessage.value = message
  replyForm.value = {
    id: message.ID,
    reply: '',
  }
  replyDialogVisible.value = true
}

// 提交回复
async function submitReply() {
  if (!replyFormRef.value)
    return

  try {
    await replyFormRef.value.validate()
    submitting.value = true

    const res: any = await request.post('/sysMessage/replyMessage', replyForm.value)
    if (res.code === 0) {
      ElMessage.success('回复成功')
      replyDialogVisible.value = false
      loadMessages()
      loadStats()
    }
  }
  catch (error) {
    ElMessage.error('回复失败')
    console.error('回复失败:', error)
  }
  finally {
    submitting.value = false
  }
}

// 删除留言
async function deleteMessage(message: any) {
  try {
    await ElMessageBox.confirm(
      '确定要删除这条留言吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      },
    )

    const res: any = await request.del('/sysMessage/deleteMessage', { data: { ID: message.ID } })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      loadMessages()
      loadStats()
    }
  }
  catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error('删除失败:', error)
    }
  }
}

// 选择变化
function handleSelectionChange(_selection: any[]) {
  // 处理选择变化
}

// 获取状态类型
function getStatusType(status: number) {
  switch (status) {
    case 1:
      return 'success'
    case 2:
      return 'warning'
    case 3:
      return 'info'
    default:
      return 'info'
  }
}

// 页面初始化
onMounted(() => {
  loadMessages()
  loadStats()
})
</script>

<template>
  <div class="message-management">
    <div class="header">
      <h2>留言管理</h2>
      <div class="stats">
        <el-card class="stat-card">
          <div class="stat-number">
            {{ stats.total || 0 }}
          </div>
          <div class="stat-label">
            总留言
          </div>
        </el-card>
        <el-card class="stat-card unread">
          <div class="stat-number">
            {{ stats.unread || 0 }}
          </div>
          <div class="stat-label">
            未读留言
          </div>
        </el-card>
        <el-card class="stat-card replied">
          <div class="stat-number">
            {{ stats.replied || 0 }}
          </div>
          <div class="stat-label">
            已回复
          </div>
        </el-card>
      </div>
    </div>

    <el-card class="message-list">
      <template #header>
        <div class="card-header">
          <span>留言列表</span>
          <div class="header-actions">
            <el-input
              v-model="searchForm.name"
              placeholder="搜索姓名"
              style="width: 200px; margin-right: 10px"
              clearable
              @clear="loadMessages"
              @keyup.enter="loadMessages"
            />
            <el-select
              v-model="searchForm.status"
              placeholder="状态筛选"
              style="width: 120px; margin-right: 10px"
              clearable
              @change="loadMessages"
            >
              <el-option label="全部" :value="0" />
              <el-option label="未读" :value="2" />
              <el-option label="已读" :value="1" />
              <el-option label="已回复" :value="3" />
            </el-select>
            <el-button type="primary" @click="loadMessages">
              搜索
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="messageList"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="phone" label="电话" width="130" />
        <el-table-column prop="content" label="留言内容" min-width="300">
          <template #default="{ row }">
            <div class="content-cell">
              {{ row.content }}
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="getStatusType(row.status)"
              size="small"
            >
              {{ row.statusText }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="CreatedAt" label="留言时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              size="small"
              @click="viewMessage(row)"
            >
              查看
            </el-button>
            <el-button
              v-if="row.status === 2"
              size="small"
              type="success"
              @click="markAsRead(row)"
            >
              标记已读
            </el-button>
            <el-button
              size="small"
              type="warning"
              @click="replyMessage(row)"
            >
              回复
            </el-button>
            <el-button
              size="small"
              type="danger"
              @click="deleteMessage(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pageInfo.page"
          v-model:page-size="pageInfo.pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadMessages"
          @current-change="loadMessages"
        />
      </div>
    </el-card>

    <!-- 留言详情对话框 -->
    <el-dialog
      v-model="messageDialogVisible"
      title="留言详情"
      width="600px"
    >
      <div v-if="currentMessage" class="message-detail">
        <div class="detail-item">
          <label>姓名：</label>
          <span>{{ currentMessage.name }}</span>
        </div>
        <div class="detail-item">
          <label>邮箱：</label>
          <span>{{ currentMessage.email }}</span>
        </div>
        <div class="detail-item">
          <label>电话：</label>
          <span>{{ currentMessage.phone || '未填写' }}</span>
        </div>
        <div class="detail-item">
          <label>IP地址：</label>

          <span>{{ currentMessage.ip }}</span>
        </div>
        <div class="detail-item">
          <label>留言时间：</label>
          <span>{{ formatDate(currentMessage.CreatedAt) }}</span>
        </div>
        <div class="detail-item">
          <label>状态：</label>
          <el-tag :type="getStatusType(currentMessage.status)">
            {{ currentMessage.statusText }}
          </el-tag>
        </div>
        <div class="detail-item">
          <label>留言内容：</label>
          <div class="content-box">
            {{ currentMessage.content }}
          </div>
        </div>
        <div v-if="currentMessage.reply" class="detail-item">
          <label>回复内容：</label>
          <div class="content-box reply">
            {{ currentMessage.reply }}
          </div>
          <div class="reply-info">
            回复时间：{{ formatDate(currentMessage.ReplyTime) }}
            <span v-if="currentMessage.ReplyUser">
              | 回复人：{{ currentMessage.ReplyUser.nickname || currentMessage.ReplyUser.username }}
            </span>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 回复对话框 -->
    <el-dialog
      v-model="replyDialogVisible"
      title="回复留言"
      width="500px"
    >
      <el-form
        ref="replyFormRef"
        :model="replyForm"
        :rules="replyRules"
        label-width="80px"
      >
        <el-form-item label="留言内容" prop="content">
          <div class="original-message">
            <strong>原留言：</strong>
            <p>{{ currentMessage?.content }}</p>
          </div>
        </el-form-item>
        <el-form-item label="回复内容" prop="reply">
          <el-input
            v-model="replyForm.reply"
            type="textarea"
            :rows="4"
            placeholder="请输入回复内容"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="replyDialogVisible = false">取消</el-button>
          <el-button
            type="primary"
            :loading="submitting"
            @click="submitReply"
          >
            提交回复
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.message-management {
  padding: 20px;
}

.header {
  margin-bottom: 20px;
}

.header h2 {
  margin: 0 0 20px 0;
  color: #303133;
}

.stats {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  flex: 1;
  text-align: center;
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.stat-card.unread {
  border-left: 4px solid #e6a23c;
}

.stat-card.replied {
  border-left: 4px solid #67c23a;
}

.stat-number {
  font-size: 2rem;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 5px;
}

.stat-label {
  color: #606266;
  font-size: 0.9rem;
}

.message-list {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  align-items: center;
}

.content-cell {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pagination-wrapper {
  margin-top: 20px;
  text-align: right;
}

.message-detail {
  padding: 20px;
}

.detail-item {
  margin-bottom: 15px;
  display: flex;
  align-items: flex-start;
}

.detail-item label {
  font-weight: bold;
  min-width: 80px;
  color: #606266;
}

.content-box {
  background: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  margin-top: 5px;
  line-height: 1.6;
}

.content-box.reply {
  background: #e1f3d8;
  border-left: 3px solid #67c23a;
}

.reply-info {
  font-size: 0.9rem;
  color: #909399;
  margin-top: 5px;
}

.original-message {
  background: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 15px;
}

.original-message p {
  margin: 5px 0 0 0;
  color: #606266;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
