<script lang="ts" setup>
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import request from '~/utils/request'

// 用户数据接口
interface User {
  ID: number
  uuid: string
  userName: string
  nickName: string
  headerImg: string
  authorityId: number
  authority: {
    authorityId: number
    authorityName: string
  }
  enable: number
  CreatedAt: string
  UpdatedAt: string
}

// 用户表单
interface UserForm {
  userName: string
  passWord: string
  nickName: string
  headerImg: string
  authorityId: number
  enable: number
  authorityIds: number[]
}

// 权限接口
interface Authority {
  authorityId: number
  authorityName: string
}

// 响应式数据
const loading = ref(false)
const userList = ref<User[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 对话框控制
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const currentUserId = ref(0)

// 表单数据
const userForm = reactive<UserForm>({
  userName: '',
  passWord: '',
  nickName: '',
  headerImg: 'https://qmplusimg.henrongyi.top/gva_header.jpg',
  authorityId: 888,
  enable: 1,
  authorityIds: [],
})

// 权限列表
const authorityList = ref<Array<{ authorityId: number, authorityName: string }>>([])

// 表单规则
const rules = {
  userName: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' },
  ],
  passWord: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' },
  ],
  nickName: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
  ],
  authorityId: [
    { required: true, message: '请选择用户权限', trigger: 'change' },
  ],
}

const formRef = ref()

// 获取用户列表
async function fetchUserList() {
  loading.value = true
  try {
    const res = await request.post('/user/getUserList', {
      page: currentPage.value,
      pageSize: pageSize.value,
    })
    if (res.code === 0) {
      userList.value = res.data.list
      total.value = res.data.total
    }
  }
  catch (error) {
    console.error('获取用户列表失败:', error)
  }
  finally {
    loading.value = false
  }
}

// 获取权限列表
async function fetchAuthorityList() {
  try {
    const res = await request.post('/authority/getAuthorityList', { page: 1, pageSize: 999 })
    if (res.code === 0) {
      authorityList.value = res.data.list
    }
  }
  catch (error) {
    console.error('获取权限列表失败:', error)
  }
}

// 打开添加用户对话框
function openAddDialog() {
  dialogTitle.value = '添加用户'
  isEdit.value = false
  currentUserId.value = 0
  resetForm()
  dialogVisible.value = true
}

// 打开编辑用户对话框
function openEditDialog(user: User) {
  dialogTitle.value = '编辑用户'
  isEdit.value = true
  currentUserId.value = user.ID
  userForm.userName = user.userName
  userForm.nickName = user.nickName
  userForm.headerImg = user.headerImg
  userForm.authorityId = user.authorityId
  userForm.enable = user.enable
  userForm.authorityIds = [user.authorityId]
  userForm.passWord = ''
  dialogVisible.value = true
}

// 重置表单
function resetForm() {
  userForm.userName = ''
  userForm.passWord = ''
  userForm.nickName = ''
  userForm.headerImg = 'https://qmplusimg.henrongyi.top/gva_header.jpg'
  userForm.authorityId = 888
  userForm.enable = 1
  userForm.authorityIds = []
  formRef.value?.resetFields()
}

// 提交表单
async function submitForm() {
  if (!formRef.value)
    return

  try {
    await formRef.value.validate()

    if (isEdit.value) {
      // 编辑用户
      const res = await request.put('/user/setUserInfo', {
        id: currentUserId.value,
        nickName: userForm.nickName,
        headerImg: userForm.headerImg,
        enable: userForm.enable,
        authorityIds: userForm.authorityIds,
      })

      if (res.code === 0) {
        ElMessage.success('编辑成功')
        dialogVisible.value = false
        fetchUserList()
      }
    }
    else {
      // 添加用户
      const res = await request.post('/user/admin_register', {
        userName: userForm.userName,
        passWord: userForm.passWord,
        nickName: userForm.nickName,
        headerImg: userForm.headerImg,
        authorityId: userForm.authorityId,
        enable: userForm.enable,
        authorityIds: userForm.authorityIds,
      })

      if (res.code === 0) {
        ElMessage.success('添加成功')
        dialogVisible.value = false
        fetchUserList()
      }
    }
  }
  catch (error) {
    console.error('提交失败:', error)
  }
}

// 删除用户
function deleteUser(user: User) {
  ElMessageBox.confirm(
    `确定要删除用户 "${user.nickName}" 吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    },
  ).then(async () => {
    try {
      const res = await request.del('/user/deleteUser', { ID: user.ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功',
        })
        fetchUserList()
      }
      else {
        ElMessage({
          type: 'error',
          message: res.msg || '删除失败',
        })
      }
    }
    catch (error) {
      ElMessage({
        type: 'error',
        message: '删除失败',
      })
      console.error('删除失败:', error)
    }
  }).catch((error: any) => {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  })
}

// 重置密码
function resetPassword(user: User) {
  ElMessageBox.confirm(
    '是否将此用户密码重置为123456?',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    },
  ).then(async () => {
    try {
      const res = await request.post('/user/resetPassword', { ID: user.ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: res.msg || '密码重置成功',
        })
      }
      else {
        ElMessage({
          type: 'error',
          message: res.msg || '密码重置失败',
        })
      }
    }
    catch (error) {
      ElMessage({
        type: 'error',
        message: '密码重置失败',
      })
      console.error('重置密码失败:', error)
    }
  }).catch((error: any) => {
    if (error !== 'cancel') {
      console.error('重置密码失败:', error)
    }
  })
}

// 切换用户状态
function toggleUserStatus(user: User) {
  const newStatus = user.enable === 1 ? 2 : 1
  const statusText = newStatus === 1 ? '启用' : '禁用'

  ElMessageBox.confirm(
    `确定要${statusText}用户 "${user.nickName}" 吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    },
  ).then(async () => {
    try {
      const res = await request.put('/user/setUserInfo', {
        ID: user.ID,
        nickName: user.nickName,
        headerImg: user.headerImg,
        enable: newStatus,
        authorityIds: [user.authorityId],
      })

      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: `${statusText}成功`,
        })
        fetchUserList()
      }
      else {
        ElMessage({
          type: 'error',
          message: res.msg || `${statusText}失败`,
        })
      }
    }
    catch (error) {
      ElMessage({
        type: 'error',
        message: `${statusText}失败`,
      })
      console.error('状态切换失败:', error)
    }
  }).catch((error: any) => {
    if (error !== 'cancel') {
      console.error('状态切换失败:', error)
    }
  })
}

// 分页变化
function handleCurrentChange(page: number) {
  currentPage.value = page
  fetchUserList()
}

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  fetchUserList()
}

// 格式化时间
function formatTime(time: string) {
  return new Date(time).toLocaleString()
}

// 获取状态文本
function getStatusText(enable: number) {
  return enable === 1 ? '正常' : '禁用'
}

// 获取状态类型
function getStatusType(enable: number) {
  return enable === 1 ? 'success' : 'danger'
}

// 初始化
onMounted(() => {
  fetchUserList()
  fetchAuthorityList()
})
</script>

<template>
  <div class="table-container">
    <div class="table-header">
      <h2 class="table-title">
        用户管理
      </h2>
      <div class="table-actions">
        <el-button type="primary" @click="openAddDialog">
          <el-icon><Plus /></el-icon>
          添加用户
        </el-button>
      </div>
    </div>

    <!-- 用户列表 -->
    <el-table
      v-loading="loading"
      :data="userList"
      style="width: 100%"
      border
    >
      <el-table-column prop="ID" label="ID" min-width="80" />
      <el-table-column prop="userName" label="用户名" min-width="120" />
      <el-table-column prop="nickName" label="昵称" min-width="120" />
      <el-table-column label="头像" min-width="80">
        <template #default="{ row }">
          <el-avatar :src="row.headerImg" :size="40" />
        </template>
      </el-table-column>
      <el-table-column prop="authority.authorityName" label="权限" min-width="120" />
      <el-table-column label="状态" min-width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.enable)" class="status-tag">
            {{ getStatusText(row.enable) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="CreatedAt" label="创建时间" min-width="180">
        <template #default="{ row }">
          {{ formatTime(row.CreatedAt) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" min-width="300" fixed="right" class-name="table-actions-column">
        <template #default="{ row }">
          <div class="action-buttons">
            <el-button type="primary" size="small" @click="openEditDialog(row)">
              编辑
            </el-button>
            <el-button
              :type="row.enable === 1 ? 'warning' : 'success'"
              size="small"
              @click="toggleUserStatus(row)"
            >
              {{ row.enable === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button type="info" size="small" @click="resetPassword(row)">
              重置密码
            </el-button>
            <el-button type="danger" size="small" @click="deleteUser(row)">
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
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 用户表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="userForm"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="userName">
          <el-input
            v-model="userForm.userName"
            placeholder="请输入用户名"
            :disabled="isEdit"
          />
        </el-form-item>

        <el-form-item v-if="!isEdit" label="密码" prop="passWord">
          <el-input
            v-model="userForm.passWord"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>

        <el-form-item label="昵称" prop="nickName">
          <el-input
            v-model="userForm.nickName"
            placeholder="请输入昵称"
          />
        </el-form-item>

        <el-form-item label="头像">
          <el-input
            v-model="userForm.headerImg"
            placeholder="请输入头像链接"
          />
        </el-form-item>

        <el-form-item v-if="!isEdit" label="用户权限" prop="authorityId">
          <el-select
            v-model="userForm.authorityId"
            placeholder="请选择用户权限"
            style="width: 100%"
          >
            <el-option
              v-for="item in authorityList"
              :key="item.authorityId"
              :label="item.authorityName"
              :value="item.authorityId"
            />
          </el-select>
        </el-form-item>

        <el-form-item v-if="isEdit" label="用户权限">
          <el-select
            v-model="userForm.authorityIds"
            multiple
            placeholder="请选择用户权限"
            style="width: 100%"
          >
            <el-option
              v-for="item in authorityList"
              :key="item.authorityId"
              :label="item.authorityName"
              :value="item.authorityId"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="状态">
          <el-radio-group v-model="userForm.enable">
            <el-radio :label="1">
              正常
            </el-radio>
            <el-radio :label="2">
              禁用
            </el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="button-group">
          <el-button @click="dialogVisible = false">
            取消
          </el-button>
          <el-button type="primary" @click="submitForm">
            {{ isEdit ? '更新' : '确定' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
/* 使用全局样式，无需自定义样式 */
</style>
