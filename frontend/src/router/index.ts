import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/Home.vue'
import PreloadVideo from '../views/PreloadVideo.vue'

// 定义路由表
const routes = [
  // 可以添加更多路由
  // 默认重定向到 main（可选）
  {
    path: '/',
    component: Home,
  },
  {
    path: '/#/preloadVideo',
    name: 'PreloadVideo',
    component: PreloadVideo,
  },
  {
    path: '/preloadVideo',
    name: 'PreloadVideo',
    component: PreloadVideo,
  },
]

const router = createRouter({
  // 使用 Hash 模式（Wails 推荐）
  history: createWebHashHistory(),
  routes,
})

export default router