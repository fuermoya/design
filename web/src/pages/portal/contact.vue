<script setup lang="ts">
import type { SysTheme, ThemeConfig } from '~/types/portal'
import { Clock, Location, Message, Phone } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, reactive, ref } from 'vue'
import PortalLayout from '~/components/layouts/PortalLayout.vue'
import request from '~/utils/request'

// 响应式数据
const themes = ref<SysTheme[]>([])
const themeConfig = ref<ThemeConfig | null>(null)

// 加载主题列表
async function loadThemes() {
  try {
    const res: any = await request.get('/portal/themes')
    if (res.code === 0) {
      themes.value = res.data.list
    }
  }
  catch (error) {
    console.error('加载主题列表失败:', error)
  }
}

// 获取当前激活的主题
async function loadActiveTheme() {
  try {
    const res: any = await request.get('/portal/theme')
    if (res.code === 0) {
      themeConfig.value = res.data.theme
    }
  }
  catch (error) {
    console.error('加载当前主题失败:', error)
  }
}

// 页面初始化
onMounted(() => {
  loadThemes()
  loadActiveTheme()
})

// 表单数据
const formData = reactive({
  name: '',
  phone: '',
  email: '',
  company: '',
  subject: '',
  message: '',
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '姓名长度在 2 到 20 个字符', trigger: 'blur' },
  ],
  phone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email' as const, message: '请输入正确的邮箱地址', trigger: 'blur' },
  ],
  company: [
    { required: true, message: '请输入公司名称', trigger: 'blur' },
  ],
  subject: [
    { required: true, message: '请选择咨询主题', trigger: 'change' },
  ],
  message: [
    { required: true, message: '请输入留言内容', trigger: 'blur' },
    { min: 10, max: 500, message: '留言内容长度在 10 到 500 个字符', trigger: 'blur' },
  ],
}

// 表单引用
const formRef = ref()

// 提交表单
async function submitForm() {
  if (!formRef.value)
    return

  try {
    await formRef.value.validate()

    // 构建提交数据
    const submitData = {
      name: formData.name,
      email: formData.email,
      phone: formData.phone,
      content: `公司：${formData.company}\n主题：${formData.subject}\n留言内容：${formData.message}`,
    }

    // 调用API提交表单数据
    const res: any = await request.post('/portal/message', submitData)

    if (res.code === 0) {
      ElMessage.success('留言提交成功，我们会尽快与您联系！')
      // 重置表单
      formRef.value.resetFields()
    }
    else {
      ElMessage.error(res.message || '提交失败，请稍后重试')
    }
  }
  catch (error: any) {
    console.error('提交留言失败:', error)
    ElMessage.error(error.response?.data?.message || '提交失败，请稍后重试')
  }
}

// 重置表单
function resetForm() {
  if (formRef.value) {
    formRef.value.resetFields()
  }
}
</script>

<template>
  <PortalLayout>
    <div class="contact-content">
      <!-- 页面头部 -->
      <div class="page-header">
        <div class="header-content">
          <h1 class="page-title">
            联系我们
          </h1>
          <p class="page-subtitle">
            我们随时准备为您提供专业的服务和支持
          </p>
        </div>
      </div>

      <!-- 主要内容区域 -->
      <div class="contact-main">
        <div class="section-container">
          <!-- 联系信息 -->
          <section class="contact-info-section">
            <div class="section-header">
              <h2 class="section-title">
                联系信息
              </h2>
              <p class="section-subtitle">
                多种方式联系我们，我们期待与您的沟通
              </p>
            </div>
            <div class="contact-info-grid">
              <div class="contact-card">
                <div class="contact-icon">
                  <el-icon><Phone /></el-icon>
                </div>
                <h3>电话咨询</h3>
                <p>{{ themeConfig?.contactPhone || '400-123-4567' }}</p>
                <p class="contact-time">
                  周一至周五 9:00-18:00
                </p>
              </div>
              <div class="contact-card">
                <div class="contact-icon">
                  <el-icon><Message /></el-icon>
                </div>
                <h3>邮箱联系</h3>
                <p>{{ themeConfig?.contactEmail || 'info@company.com' }}</p>
                <p class="contact-time">
                  24小时内回复
                </p>
              </div>
              <div class="contact-card">
                <div class="contact-icon">
                  <el-icon><Location /></el-icon>
                </div>
                <h3>公司地址</h3>
                <p>{{ themeConfig?.contactAddress || '北京市朝阳区xxx大厦' }}</p>
                <p class="contact-time">
                  欢迎来访
                </p>
              </div>
              <div class="contact-card">
                <div class="contact-icon">
                  <el-icon><Clock /></el-icon>
                </div>
                <h3>工作时间</h3>
                <p>周一至周五</p>
                <p class="contact-time">
                  9:00-18:00
                </p>
              </div>
            </div>
          </section>

          <!-- 留言表单 -->
          <section class="contact-form-section">
            <div class="section-header">
              <h2 class="section-title">
                在线留言
              </h2>
              <p class="section-subtitle">
                填写以下表单，我们会尽快与您联系
              </p>
            </div>
            <div class="form-container">
              <el-form
                ref="formRef"
                :model="formData"
                :rules="formRules"
                label-width="100px"
                class="contact-form"
              >
                <div class="form-row">
                  <el-form-item label="姓名" prop="name">
                    <el-input v-model="formData.name" placeholder="请输入您的姓名" />
                  </el-form-item>
                  <el-form-item label="联系电话" prop="phone">
                    <el-input v-model="formData.phone" placeholder="请输入您的联系电话" />
                  </el-form-item>
                </div>
                <div class="form-row">
                  <el-form-item label="邮箱地址" prop="email">
                    <el-input v-model="formData.email" placeholder="请输入您的邮箱地址" />
                  </el-form-item>
                  <el-form-item label="公司名称" prop="company">
                    <el-input v-model="formData.company" placeholder="请输入您的公司名称" />
                  </el-form-item>
                </div>
                <el-form-item label="咨询主题" prop="subject">
                  <el-select v-model="formData.subject" placeholder="请选择咨询主题" style="width: 100%">
                    <el-option label="产品咨询" value="product" />
                    <el-option label="服务咨询" value="service" />
                    <el-option label="合作洽谈" value="cooperation" />
                    <el-option label="技术支持" value="technical" />
                    <el-option label="其他" value="other" />
                  </el-select>
                </el-form-item>
                <el-form-item label="留言内容" prop="message">
                  <el-input
                    v-model="formData.message"
                    type="textarea"
                    :rows="6"
                    placeholder="请详细描述您的需求或问题..."
                  />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" size="large" @click="submitForm">
                    提交留言
                  </el-button>
                  <el-button size="large" @click="resetForm">
                    重置表单
                  </el-button>
                </el-form-item>
              </el-form>
            </div>
          </section>

          <!-- 地图区域 -->
          <section class="map-section">
            <div class="section-header">
              <h2 class="section-title">
                公司位置
              </h2>
              <p class="section-subtitle">
                欢迎您到访我们的办公地点
              </p>
            </div>
            <div class="map-container">
              <div class="map-placeholder">
                <el-icon><Location /></el-icon>
                <p>地图加载中...</p>
                <p class="map-address">
                  {{ themeConfig?.contactAddress || '北京市朝阳区xxx大厦' }}
                </p>
              </div>
            </div>
          </section>
        </div>
      </div>
    </div>
  </PortalLayout>
</template>

<style scoped>
.contact-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

/* 页面头部 */
.page-header {
  background: linear-gradient(
    135deg,
    var(--portal-primary-color, #409eff) 0%,
    var(--portal-secondary-color, #67c23a) 100%
  );
  color: white;
  padding: 3rem 2rem;
  border-radius: 12px;
  margin-bottom: 2rem;
  text-align: center;
}

.header-content {
  max-width: 800px;
  margin: 0 auto;
}

.page-title {
  font-size: 2.5rem;
  font-weight: bold;
  margin: 0 0 1rem;
}

.page-subtitle {
  font-size: 1.1rem;
  opacity: 0.9;
  margin: 0;
  line-height: 1.6;
}

/* 主要内容区域 */
.contact-main {
  margin-bottom: 3rem;
}

.section-container {
  max-width: 1200px;
  margin: 0 auto;
}

.section-header {
  text-align: center;
  margin-bottom: 3rem;
}

.section-title {
  font-size: 2rem;
  font-weight: bold;
  color: var(--portal-text-color, #303133);
  margin: 0 0 1rem;
}

.section-subtitle {
  font-size: 1.1rem;
  color: #606266;
  margin: 0;
  line-height: 1.6;
}

/* 联系信息 */
.contact-info-section {
  margin-bottom: 4rem;
}

.contact-info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
}

.contact-card {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  text-align: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease;
}

.contact-card:hover {
  transform: translateY(-4px);
}

.contact-icon {
  width: 80px;
  height: 80px;
  background: linear-gradient(
    135deg,
    var(--portal-primary-color, #409eff) 0%,
    var(--portal-secondary-color, #67c23a) 100%
  );
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1.5rem;
  color: white;
  font-size: 2rem;
}

.contact-card h3 {
  font-size: 1.3rem;
  font-weight: bold;
  color: var(--portal-text-color, #303133);
  margin: 0 0 1rem;
}

.contact-card p {
  color: #606266;
  line-height: 1.6;
  margin: 0 0 0.5rem;
}

.contact-time {
  font-size: 0.9rem;
  color: #909399;
}

/* 留言表单 */
.contact-form-section {
  margin-bottom: 4rem;
}

.form-container {
  max-width: 800px;
  margin: 0 auto;
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.contact-form {
  width: 100%;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
}

/* 地图区域 */
.map-section {
  margin-bottom: 2rem;
}

.map-container {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.map-placeholder {
  height: 400px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #909399;
}

.map-placeholder .el-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
}

.map-placeholder p {
  margin: 0.5rem 0;
  font-size: 1.1rem;
}

.map-address {
  font-size: 1rem;
  color: #606266;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .contact-content {
    padding: 0 1rem;
  }

  .page-header {
    padding: 2rem 1rem;
  }

  .page-title {
    font-size: 2rem;
  }

  .contact-info-grid {
    grid-template-columns: 1fr;
  }

  .form-row {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .form-container {
    padding: 1.5rem;
  }
}
</style>
