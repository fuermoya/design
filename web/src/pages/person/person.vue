<script setup lang="ts">
import type { Res } from '~/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import { useUserStore } from '~/stores/user'
import request from '~/utils/request'

const userStore = useUserStore()

// 表单数据
const formData = reactive({
  nickName: '',
  password: '',
  newPassword: '',
  confirmPassword: '',
})

// 头像上传相关
const avatarFile = ref<File | null>(null)
const avatarPreview = ref('')
const avatarLoading = ref(false)
const avatarUploadRef = ref()

// 媒体库相关
const mediaLibraryVisible = ref(false)
const mediaLibraryLoading = ref(false)
const mediaFileList = ref<any[]>([])
const mediaSearchForm = reactive({
  name: '',
  tag: '',
  page: 1,
  pageSize: 10,
})
const mediaTotal = ref(0)

// 表单验证规则
const rules = {
  nickName: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '昵称长度在 2 到 20 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入原密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' },
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: any) => {
        if (value !== formData.newPassword) {
          callback(new Error('两次输入密码不一致'))
        }
        else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
}

// 表单引用
const nickNameFormRef = ref()
const passwordFormRef = ref()

// 加载状态
const nickNameLoading = ref(false)
const passwordLoading = ref(false)

// 初始化用户信息
onMounted(() => {
  formData.nickName = userStore.user.nickName || ''
})

// 处理头像文件选择
function handleAvatarChange(uploadFile: any) {
  if (uploadFile.raw) {
    const file = uploadFile.raw

    // 验证文件类型
    if (!file.type.startsWith('image/')) {
      ElMessage.error('请选择图片文件')
      return
    }

    // 验证文件大小（限制为2MB）
    if (file.size > 2 * 1024 * 1024) {
      ElMessage.error('图片大小不能超过2MB')
      return
    }

    avatarFile.value = file

    // 创建预览
    const reader = new FileReader()
    reader.onload = (e) => {
      avatarPreview.value = e.target?.result as string
    }
    reader.readAsDataURL(file)
  }
}

// 打开媒体库
function openMediaLibrary() {
  mediaLibraryVisible.value = true
  getMediaFileList()
}

// 获取媒体库文件列表
async function getMediaFileList() {
  mediaLibraryLoading.value = true
  try {
    const res: any = await request.get('/fileUploadAndDownload/getFileList', {
      name: mediaSearchForm.name,
      tag: mediaSearchForm.tag,
      page: mediaSearchForm.page,
      pageSize: mediaSearchForm.pageSize,
    })
    if (res.code === 0) {
      mediaFileList.value = res.data.list
      mediaTotal.value = res.data.total
    }
  }
  catch (error) {
    console.error('获取媒体库文件列表失败:', error)
  }
  finally {
    mediaLibraryLoading.value = false
  }
}

// 从媒体库选择头像
function selectFromMediaLibrary(file: any) {
  try {
    // 验证是否为图片文件
    const imageExtensions = ['jpg', 'jpeg', 'png', 'gif', 'webp']
    const fileExtension = file.name.split('.').pop()?.toLowerCase()

    if (!imageExtensions.includes(fileExtension)) {
      ElMessage.error('请选择图片文件')
      return
    }

    // 设置头像预览
    avatarPreview.value = `${import.meta.env.VITE_BASE_API || '/api'}${file.url}`
    avatarFile.value = null // 清空本地文件，因为使用媒体库文件

    // 关闭媒体库弹窗
    mediaLibraryVisible.value = false

    ElMessage.success('已选择图片，请点击确认设置头像')
  }
  catch (error) {
    console.error('选择媒体库文件失败:', error)
    ElMessage.error('选择文件失败')
  }
}

// 上传新图片到媒体库
async function uploadToMediaLibrary() {
  if (!avatarFile.value) {
    ElMessage.warning('请先选择图片文件')
    return
  }

  try {
    avatarLoading.value = true

    const formData = new FormData()
    formData.append('file', avatarFile.value)

    const res = await request.upload('/fileUploadAndDownload/upload', formData) as Res

    if (res.code === 0) {
      ElMessage.success('图片上传到媒体库成功')
      // 刷新媒体库列表
      getMediaFileList()
      // 设置预览为新上传的图片
      avatarPreview.value = res.data.url || res.data.headerImg
      avatarFile.value = null
      // 重置文件输入框
      if (avatarUploadRef.value) {
        avatarUploadRef.value.value = ''
      }
      // 关闭媒体库弹窗
      mediaLibraryVisible.value = false
    }
    else {
      ElMessage.error(res.msg || '图片上传失败')
    }
  }
  catch (error) {
    console.error('上传图片到媒体库失败:', error)
    ElMessage.error('图片上传失败')
  }
  finally {
    avatarLoading.value = false
  }
}

// 确认设置头像
async function confirmSetAvatar() {
  if (!avatarPreview.value) {
    ElMessage.warning('请先选择或上传头像')
    return
  }

  try {
    avatarLoading.value = true

    const res = await request.put('/user/setSelfInfo', {
      headerImg: avatarPreview.value,
    }) as Res

    if (res.code === 0) {
      ElMessage.success('头像设置成功')
      // 更新本地用户信息
      userStore.user.headerImg = avatarPreview.value
      // 清空预览
      avatarPreview.value = ''
      avatarFile.value = null
      // 重置文件输入框
      if (avatarUploadRef.value) {
        avatarUploadRef.value.value = ''
      }
      // 关闭媒体库
      mediaLibraryVisible.value = false
    }
    else {
      ElMessage.error(res.msg || '头像设置失败')
    }
  }
  catch (error) {
    console.error('设置头像失败:', error)
    ElMessage.error('头像设置失败')
  }
  finally {
    avatarLoading.value = false
  }
}

// 媒体库搜索
function handleMediaSearch() {
  mediaSearchForm.page = 1
  getMediaFileList()
}

// 媒体库重置搜索
function handleMediaReset() {
  mediaSearchForm.name = ''
  mediaSearchForm.tag = ''
  mediaSearchForm.page = 1
  getMediaFileList()
}

// 媒体库分页
function handleMediaCurrentChange(page: number) {
  mediaSearchForm.page = page
  getMediaFileList()
}

function handleMediaSizeChange(size: number) {
  mediaSearchForm.pageSize = size
  mediaSearchForm.page = 1
  getMediaFileList()
}

// 更新昵称
async function updateNickName() {
  try {
    await nickNameFormRef.value.validate()
    nickNameLoading.value = true

    const res = await request.put('/user/setSelfInfo', {
      nickName: formData.nickName,
    }) as Res

    if (res.code === 0) {
      ElMessage.success('昵称更新成功')
      // 更新本地用户信息
      userStore.user.nickName = formData.nickName
    }
    else {
      ElMessage.error(res.msg || '昵称更新失败')
    }
  }
  catch (error) {
    console.error('更新昵称失败:', error)
    ElMessage.error('昵称更新失败')
  }
  finally {
    nickNameLoading.value = false
  }
}

// 更新密码
async function updatePassword() {
  try {
    await passwordFormRef.value.validate()
    passwordLoading.value = true

    const res = await request.post('/user/changePassword', {
      password: formData.password,
      newPassword: formData.newPassword,
    }) as Res

    if (res.code === 0) {
      ElMessage.success('密码更新成功，请重新登录')
      // 清空密码表单
      formData.password = ''
      formData.newPassword = ''
      formData.confirmPassword = ''
      passwordFormRef.value.resetFields()

      // 退出登录
      setTimeout(() => {
        userStore.resetUserInfo()
        localStorage.removeItem('x-token')
        window.location.href = '/login'
      }, 1500)
    }
    else {
      ElMessage.error(res.msg || '密码更新失败')
    }
  }
  catch (error) {
    console.error('更新密码失败:', error)
    ElMessage.error('密码更新失败')
  }
  finally {
    passwordLoading.value = false
  }
}

// 重置昵称表单
function resetNickNameForm() {
  formData.nickName = userStore.user.nickName || ''
  nickNameFormRef.value.resetFields()
}

// 重置密码表单
function resetPasswordForm() {
  formData.password = ''
  formData.newPassword = ''
  formData.confirmPassword = ''
  passwordFormRef.value.resetFields()
}
</script>

<template>
  <div class="person-container">
    <el-card class="person-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <i class="el-icon-user" style="margin-right: 8px;" />
            个人信息管理
          </span>
        </div>
      </template>

      <el-row :gutter="24">
        <!-- 基本信息 -->
        <el-col :span="12">
          <el-card class="info-card" shadow="never">
            <template #header>
              <div class="info-header">
                <span>基本信息</span>
              </div>
            </template>

            <div class="user-info">
              <div class="avatar-section">
                <div class="avatar-container">
                  <el-avatar
                    :size="80"
                    :src="avatarPreview || userStore.user.headerImg"
                    :icon="(avatarPreview || userStore.user.headerImg) ? undefined : 'el-icon-user'"
                  />
                  <div class="avatar-overlay">
                    <el-button
                      type="primary"
                      size="small"
                      icon="el-icon-camera"
                      @click="openMediaLibrary"
                    >
                      更换头像
                    </el-button>
                  </div>
                </div>
                <div class="user-details">
                  <h3>{{ userStore.user.userName }}</h3>
                  <p class="user-role">
                    {{ userStore.user.authority?.authorityName || '普通用户' }}
                  </p>
                </div>
              </div>

              <!-- 头像预览和操作按钮 -->
              <div v-if="avatarPreview" class="avatar-preview-section">
                <el-divider />
                <div class="preview-container">
                  <div class="preview-info">
                    <p class="preview-text">
                      头像预览：
                    </p>
                    <el-image
                      :src="avatarPreview"
                      style="width: 60px; height: 60px; border-radius: 50%;"
                      fit="cover"
                    />
                  </div>
                  <div class="preview-actions">
                    <el-button
                      type="primary"
                      size="small"
                      :loading="avatarLoading"
                      @click="confirmSetAvatar"
                    >
                      确认设置头像
                    </el-button>
                    <el-button
                      size="small"
                      @click="() => { avatarPreview = ''; avatarFile = null; if (avatarUploadRef.value) avatarUploadRef.value.value = '' }"
                    >
                      取消
                    </el-button>
                  </div>
                </div>
              </div>

              <el-divider />

              <div class="info-item">
                <label>用户ID：</label>
                <span>{{ userStore.user.ID }}</span>
              </div>
              <div class="info-item">
                <label>用户名：</label>
                <span>{{ userStore.user.userName }}</span>
              </div>
              <div class="info-item">
                <label>当前昵称：</label>
                <span>{{ userStore.user.nickName || '未设置' }}</span>
              </div>
              <div class="info-item">
                <label>状态：</label>
                <el-tag :type="userStore.user.enable ? 'success' : 'danger'">
                  {{ userStore.user.enable ? '启用' : '禁用' }}
                </el-tag>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 修改昵称 -->
        <el-col :span="12">
          <el-card class="form-card" shadow="never">
            <template #header>
              <div class="form-header">
                <span>修改昵称</span>
              </div>
            </template>

            <el-form
              ref="nickNameFormRef"
              :model="formData"
              :rules="rules"
              label-width="80px"
              size="default"
            >
              <el-form-item label="昵称" prop="nickName">
                <el-input
                  v-model="formData.nickName"
                  placeholder="请输入新的昵称"
                  maxlength="20"
                  show-word-limit
                />
              </el-form-item>

              <el-form-item>
                <el-button
                  type="primary"
                  :loading="nickNameLoading"
                  @click="updateNickName"
                >
                  更新昵称
                </el-button>
                <el-button @click="resetNickNameForm">
                  重置
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-col>
      </el-row>

      <!-- 修改密码 -->
      <el-card class="password-card" shadow="never" style="margin-top: 20px;">
        <template #header>
          <div class="form-header">
            <span>修改密码</span>
          </div>
        </template>

        <el-form
          ref="passwordFormRef"
          :model="formData"
          :rules="rules"
          label-width="100px"
          size="default"
          style="max-width: 500px;"
        >
          <el-form-item label="原密码" prop="password">
            <el-input
              v-model="formData.password"
              type="password"
              placeholder="请输入原密码"
              show-password
              maxlength="20"
            />
          </el-form-item>

          <el-form-item label="新密码" prop="newPassword">
            <el-input
              v-model="formData.newPassword"
              type="password"
              placeholder="请输入新密码"
              show-password
              maxlength="20"
            />
          </el-form-item>

          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input
              v-model="formData.confirmPassword"
              type="password"
              placeholder="请再次输入新密码"
              show-password
              maxlength="20"
            />
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              :loading="passwordLoading"
              @click="updatePassword"
            >
              更新密码
            </el-button>
            <el-button @click="resetPasswordForm">
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </el-card>

    <!-- 媒体库选择器对话框 -->
    <el-dialog
      v-model="mediaLibraryVisible"
      title="选择头像"
      width="80%"
      :close-on-click-modal="false"
    >
      <div class="media-library-container">
        <!-- 上传新图片到媒体库 -->
        <el-card class="upload-card" shadow="hover" style="margin-bottom: 20px;">
          <template #header>
            <div class="card-header">
              <span>上传新图片到媒体库</span>
            </div>
          </template>

          <el-form label-width="80px">
            <el-form-item label="选择图片">
              <el-upload
                class="upload-demo"
                :auto-upload="false"
                :show-file-list="false"
                :on-change="handleAvatarChange"
                accept="image/*"
              >
                <el-button type="primary">
                  选择图片
                </el-button>
                <template #tip>
                  <div class="el-upload__tip">
                    支持 JPG、PNG、GIF、WebP 格式，文件大小不超过 2MB
                  </div>
                </template>
              </el-upload>
            </el-form-item>

            <el-form-item>
              <el-button
                type="success"
                :loading="avatarLoading"
                :disabled="!avatarFile"
                @click="uploadToMediaLibrary"
              >
                上传到媒体库
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 从媒体库选择 -->
        <el-card class="search-card" shadow="hover" style="margin-bottom: 20px;">
          <template #header>
            <div class="card-header">
              <span>从媒体库选择图片</span>
            </div>
          </template>

          <el-form :model="mediaSearchForm" inline>
            <el-form-item label="文件名">
              <el-input
                v-model="mediaSearchForm.name"
                placeholder="请输入文件名"
                clearable
                @keyup.enter="handleMediaSearch"
              />
            </el-form-item>

            <el-form-item label="文件类型">
              <el-input
                v-model="mediaSearchForm.tag"
                placeholder="请输入文件类型"
                clearable
                @keyup.enter="handleMediaSearch"
              />
            </el-form-item>

            <el-form-item>
              <el-button type="primary" @click="handleMediaSearch">
                搜索
              </el-button>
              <el-button @click="handleMediaReset">
                重置
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 媒体库文件列表 -->
        <el-card class="list-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>媒体库文件列表</span>
              <el-button type="primary" :loading="mediaLibraryLoading" @click="getMediaFileList">
                刷新
              </el-button>
            </div>
          </template>

          <el-table
            v-loading="mediaLibraryLoading"
            :data="mediaFileList"
            style="width: 100%"
          >
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column prop="name" label="文件名" min-width="200" />
            <el-table-column prop="tag" label="文件类型" width="100" />
            <el-table-column label="上传时间" width="180">
              <template #default="scope">
                {{ new Date(scope.row.created_at).toLocaleString() }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="scope">
                <el-button
                  type="primary"
                  size="small"
                  @click="selectFromMediaLibrary(scope.row)"
                >
                  选择
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination-wrapper" style="margin-top: 20px;">
            <el-pagination
              v-model:current-page="mediaSearchForm.page"
              v-model:page-size="mediaSearchForm.pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="mediaTotal"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleMediaSizeChange"
              @current-change="handleMediaCurrentChange"
            />
          </div>
        </el-card>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="mediaLibraryVisible = false">
            取消
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<style lang="scss" scoped>
// 个人信息页面使用全局样式，这里只保留特定的自定义样式
</style>
