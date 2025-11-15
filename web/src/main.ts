import type { UserModule } from './types'
// 统一导入el-icon图标
import * as ElIconModules from '@element-plus/icons-vue'
import { ViteSSG } from 'vite-ssg'
import { routes } from 'vue-router/auto-routes'
import App from './App.vue'
import 'element-plus/es/components/message/style/css'

import 'element-plus/es/components/loading/style/css'

// import "~/styles/element/index.scss";

// import ElementPlus from "element-plus";
// import all element css, uncommented next line
// import "element-plus/dist/index.css";

// or use cdn, uncomment cdn link in `index.html`

import 'element-plus/es/components/notification/style/css'
// 引入全局类型声明

import 'element-plus/es/components/message-box/style/css'

import '~/styles/index.scss'
import 'uno.css'
// If you want to use ElMessage, import it.
// import 'element-plus/theme-chalk/el-message.css'
// import 'element-plus/theme-chalk/el-message-box.css'
// if you do not need ssg:
// import { createApp } from "vue";

// const app = createApp(App);
// app.use(createRouter({
//   history: createWebHistory(),
//   routes,
// }))
// // app.use(ElementPlus);
// app.mount("#app");

// https://github.com/antfu/vite-ssg
export const createApp = ViteSSG(
  App,
  {
    routes,
    base: import.meta.env.BASE_URL,
  },
  (ctx) => {
    // 注册所有 Element Plus 图标
    for (const [key, component] of Object.entries(ElIconModules)) {
      ctx.app.component(key, component)
    }

    // install all modules under `modules/`
    Object.values(import.meta.glob<{ install: UserModule }>('./modules/*.ts', { eager: true }))
      .forEach(i => i.install?.(ctx))
    // ctx.app.use(Previewer)
  },
)
