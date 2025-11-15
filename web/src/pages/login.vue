<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { ref } from 'vue'
import { useUserStore } from '~/stores/user'
import request from '~/utils/request'

const userStore = useUserStore()
const loading = ref(false)
const formRef = ref(null)
const formData = ref({
  username: '',
  password: '',
  remember: false,
  captcha: '',
  captchaId: '',
})

const rules = ref({
  username: [{ required: true, message: '请输入用户名' }],
  password: [{ required: true, message: '请输入密码' }],
  captcha: [{ required: true, message: '请输入验证码' }],
})

const env = import.meta.env

const captchaUrl = ref('')
const captchaId = ref('')
interface captchaInfo {
  code: number
  data: {
    username: string
    password: string
    picPath: string
    captchaId: string
  }
  msg: string
}

async function refreshCaptcha() {
  try {
    const res = (await request.post('/base/captcha', {}, {
      headers: {
        useToken: false,
      },
    })) as captchaInfo
    captchaUrl.value = res.data.picPath
    captchaId.value = res.data.captchaId
    formData.value.captchaId = res.data.captchaId
  }
  catch (error) {
    console.error('更新密码失败:', error)
    ElMessage.error('获取验证码失败')
  }
}

async function login() {
  if (!formRef.value) {
    return
  }
  formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      formData.value.captchaId = captchaId.value
      await userStore.login(formData.value, loading)
      if (userStore.user.ID <= 0) {
        refreshCaptcha()
      }
    }
  })
}

refreshCaptcha()
</script>

<template>
  <div class="h-screen flex flex-col items-center justify-center bg-zinc-100">
    <div class="mb-12 -mt-12">
      <img src="/favicon-96x96.png" alt="">
    </div>
    <div class="rounded bg-white p-12 shadow-sm">
      <div class="mb-6 text-center">
        <p v-if="env.VITE_LOGIN_TITLE" class="mb-6 text-3xl text-zinc-800">
          {{ env.VITE_LOGIN_TITLE }}
        </p>
        <p v-if="env.VITE_LOGIN_DESC" class="mb-6 text-sm text-zinc-500">
          {{ env.VITE_LOGIN_DESC }}
        </p>
      </div>
      <div class="mx-auto w-320px">
        <el-form ref="formRef" size="large" :model="formData" :rules="rules" @submit.prevent="login">
          <el-form-item label="" prop="username">
            <el-input v-model="formData.username" placeholder="请输入用户名">
              <template #prefix>
                <i class="ri-user-line" />
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="" prop="password">
            <el-input
              v-model="formData.password" placeholder="请输入密码" type="password" show-password
              @keyup.enter="login"
            >
              <template #prefix>
                <i class="ri-lock-line" />
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="" prop="captcha">
            <div style="display: flex; align-items: center; gap: 12px;">
              <el-input v-model="formData.captcha" placeholder="请输入验证码" style="flex: 1;" />
              <div
                style="width: 120px; height: 40px; cursor: pointer; border: 1px solid #dcdfe6; border-radius: 4px; overflow: hidden;"
                @click="refreshCaptcha"
              >
                <img
                  :src="captchaUrl" alt="验证码"
                  style="width: 100%; height: 100%; object-fit: cover;"
                >
              </div>
            </div>
          </el-form-item>

          <el-form-item>
            <!-- TODO 加 loading -->
            <el-button class="w-full" type="primary" :loading="loading" @click="login">
              登 录
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
