<script setup lang="ts">
import type { SysTheme, ThemeConfig } from '~/types/portal'
import { Aim, OfficeBuilding, Star, User, View } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onMounted, ref } from 'vue'
import PortalLayout from '~/components/layouts/PortalLayout.vue'
import request from '~/utils/request'

// 响应式数据
const themes = ref<SysTheme[]>([])
const currentTheme = ref<SysTheme | null>(null)
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
</script>

<template>
  <PortalLayout>
    <div class="about-content">
      <!-- 页面头部 -->
      <div class="page-header">
        <div class="header-content">
          <h1 class="page-title">
            关于我们
          </h1>
          <p class="page-subtitle">
            {{ themeConfig?.siteDescription || '专业的服务，卓越的品质' }}
          </p>
        </div>
      </div>

      <!-- 主要内容区域 -->
      <div class="about-main">
        <div class="section-container">
          <!-- 公司介绍 -->
          <section class="company-section">
            <div class="section-header">
              <h2 class="section-title">
                公司简介
              </h2>
              <p class="section-subtitle">
                我们是一家专注于创新和发展的现代化企业，致力于为客户提供最优质的产品和服务。
              </p>
            </div>
            <div class="company-content">
              <div class="company-text">
                <p>
                  自成立以来，我们始终坚持"客户至上、质量第一"的经营理念，通过不断的技术创新和服务优化，
                  为客户创造更大的价值。我们的团队由经验丰富的专业人士组成，具备深厚的行业知识和专业技能。
                </p>
                <p>
                  我们相信，只有通过持续的学习和改进，才能在竞争激烈的市场中保持领先地位。
                  因此，我们不断投资于员工培训和技术研发，确保我们始终能够提供最前沿的解决方案。
                </p>
              </div>
              <div class="company-image">
                <div class="image-placeholder">
                  <el-icon><OfficeBuilding /></el-icon>
                  <p>公司大楼</p>
                </div>
              </div>
            </div>
          </section>

          <!-- 核心价值观 -->
          <section class="values-section">
            <div class="section-header">
              <h2 class="section-title">
                核心价值观
              </h2>
              <p class="section-subtitle">
                我们的价值观指导着我们的每一个决策和行动
              </p>
            </div>
            <div class="values-grid">
              <div class="value-card">
                <div class="value-icon">
                  <el-icon><Star /></el-icon>
                </div>
                <h3 class="value-title">
                  卓越品质
                </h3>
                <p class="value-description">
                  我们追求卓越，在每一个细节上都力求完美，为客户提供最优质的产品和服务。
                </p>
              </div>
              <div class="value-card">
                <div class="value-icon">
                  <el-icon><User /></el-icon>
                </div>
                <h3 class="value-title">
                  客户至上
                </h3>
                <p class="value-description">
                  客户的需求是我们的首要考虑，我们始终以客户满意为目标，提供贴心的服务。
                </p>
              </div>
              <div class="value-card">
                <div class="value-icon">
                  <el-icon><Aim /></el-icon>
                </div>
                <h3 class="value-title">
                  创新驱动
                </h3>
                <p class="value-description">
                  我们鼓励创新思维，不断探索新的可能性，推动行业的发展和进步。
                </p>
              </div>
              <div class="value-card">
                <div class="value-icon">
                  <el-icon><View /></el-icon>
                </div>
                <h3 class="value-title">
                  诚信透明
                </h3>
                <p class="value-description">
                  我们坚持诚信经营，保持透明度，与客户、员工和合作伙伴建立长期信任关系。
                </p>
              </div>
            </div>
          </section>

          <!-- 团队介绍 -->
          <section class="team-section">
            <div class="section-header">
              <h2 class="section-title">
                我们的团队
              </h2>
              <p class="section-subtitle">
                由经验丰富的专业人士组成的精英团队
              </p>
            </div>
            <div class="team-content">
              <div class="team-stats">
                <div class="stat-item">
                  <div class="stat-number">
                    50+
                  </div>
                  <div class="stat-label">
                    专业员工
                  </div>
                </div>
                <div class="stat-item">
                  <div class="stat-number">
                    10+
                  </div>
                  <div class="stat-label">
                    年行业经验
                  </div>
                </div>
                <div class="stat-item">
                  <div class="stat-number">
                    1000+
                  </div>
                  <div class="stat-label">
                    成功项目
                  </div>
                </div>
                <div class="stat-item">
                  <div class="stat-number">
                    98%
                  </div>
                  <div class="stat-label">
                    客户满意度
                  </div>
                </div>
              </div>
              <div class="team-description">
                <p>
                  我们的团队由来自不同领域的专业人士组成，包括技术专家、行业顾问、项目经理等。
                  每位团队成员都具备深厚的专业知识和丰富的实践经验，能够为客户提供全方位的解决方案。
                </p>
                <p>
                  我们注重团队协作，鼓励知识分享和经验交流，确保每个项目都能得到最佳的执行效果。
                  通过定期的培训和技能提升，我们确保团队始终保持行业领先水平。
                </p>
              </div>
            </div>
          </section>

          <!-- 发展历程 -->
          <section class="history-section">
            <div class="section-header">
              <h2 class="section-title">
                发展历程
              </h2>
              <p class="section-subtitle">
                见证我们的成长和进步
              </p>
            </div>
            <div class="timeline">
              <div class="timeline-item">
                <div class="timeline-year">
                  2014
                </div>
                <div class="timeline-content">
                  <h3>公司成立</h3>
                  <p>公司正式成立，开始专注于行业解决方案的开发和实施。</p>
                </div>
              </div>
              <div class="timeline-item">
                <div class="timeline-year">
                  2016
                </div>
                <div class="timeline-content">
                  <h3>业务扩展</h3>
                  <p>成功拓展业务范围，服务客户数量突破100家。</p>
                </div>
              </div>
              <div class="timeline-item">
                <div class="timeline-year">
                  2018
                </div>
                <div class="timeline-content">
                  <h3>技术突破</h3>
                  <p>获得多项技术专利，成为行业技术领导者。</p>
                </div>
              </div>
              <div class="timeline-item">
                <div class="timeline-year">
                  2020
                </div>
                <div class="timeline-content">
                  <h3>数字化转型</h3>
                  <p>全面启动数字化转型，为客户提供更智能的解决方案。</p>
                </div>
              </div>
              <div class="timeline-item">
                <div class="timeline-year">
                  2024
                </div>
                <div class="timeline-content">
                  <h3>持续创新</h3>
                  <p>继续引领行业发展，为客户创造更大的价值。</p>
                </div>
              </div>
            </div>
          </section>
        </div>
      </div>
    </div>
  </PortalLayout>
</template>

<style scoped>
.about-content {
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
.about-main {
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

/* 公司介绍 */
.company-section {
  margin-bottom: 4rem;
}

.company-content {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 3rem;
  align-items: center;
}

.company-text p {
  margin: 0 0 1.5rem;
  line-height: 1.8;
  color: #606266;
  font-size: 1.1rem;
}

.company-image {
  text-align: center;
}

.image-placeholder {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  border-radius: 12px;
  padding: 3rem 2rem;
  color: #909399;
}

.image-placeholder .el-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
}

.image-placeholder p {
  margin: 0;
  font-size: 1.1rem;
}

/* 核心价值观 */
.values-section {
  margin-bottom: 4rem;
}

.values-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
}

.value-card {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  text-align: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease;
}

.value-card:hover {
  transform: translateY(-4px);
}

.value-icon {
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

.value-title {
  font-size: 1.3rem;
  font-weight: bold;
  color: var(--portal-text-color, #303133);
  margin: 0 0 1rem;
}

.value-description {
  color: #606266;
  line-height: 1.6;
  margin: 0;
}

/* 团队介绍 */
.team-section {
  margin-bottom: 4rem;
}

.team-content {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 3rem;
  align-items: center;
}

.team-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 2rem;
}

.stat-item {
  text-align: center;
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-number {
  font-size: 2.5rem;
  font-weight: bold;
  color: var(--portal-primary-color, #409eff);
  margin-bottom: 0.5rem;
}

.stat-label {
  color: #606266;
  font-size: 1rem;
}

.team-description p {
  margin: 0 0 1.5rem;
  line-height: 1.8;
  color: #606266;
  font-size: 1.1rem;
}

/* 发展历程 */
.history-section {
  margin-bottom: 2rem;
}

.timeline {
  position: relative;
  padding-left: 2rem;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  background: linear-gradient(to bottom, var(--portal-primary-color, #409eff), var(--portal-secondary-color, #67c23a));
}

.timeline-item {
  position: relative;
  margin-bottom: 2rem;
  padding-left: 2rem;
}

.timeline-item::before {
  content: '';
  position: absolute;
  left: -1rem;
  top: 0.5rem;
  width: 12px;
  height: 12px;
  background: var(--portal-primary-color, #409eff);
  border-radius: 50%;
  border: 3px solid white;
  box-shadow: 0 0 0 2px var(--portal-primary-color, #409eff);
}

.timeline-year {
  font-size: 1.2rem;
  font-weight: bold;
  color: var(--portal-primary-color, #409eff);
  margin-bottom: 0.5rem;
}

.timeline-content h3 {
  font-size: 1.1rem;
  font-weight: bold;
  color: var(--portal-text-color, #303133);
  margin: 0 0 0.5rem;
}

.timeline-content p {
  color: #606266;
  line-height: 1.6;
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .about-content {
    padding: 0 1rem;
  }

  .page-header {
    padding: 2rem 1rem;
  }

  .page-title {
    font-size: 2rem;
  }

  .company-content {
    grid-template-columns: 1fr;
    gap: 2rem;
  }

  .team-content {
    grid-template-columns: 1fr;
    gap: 2rem;
  }

  .team-stats {
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
  }

  .stat-item {
    padding: 1.5rem;
  }

  .stat-number {
    font-size: 2rem;
  }

  .values-grid {
    grid-template-columns: 1fr;
  }
}
</style>
