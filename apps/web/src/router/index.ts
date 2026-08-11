import { createRouter, createWebHashHistory } from 'vue-router'
import { constantRoute } from './routes'

const router = createRouter({
  // 哈希模式（部署简单，不需要服务器配置）
  history: createWebHashHistory(),
  routes: constantRoute,
  // 路由切换时滚动到顶部
  scrollBehavior() {
    return { left: 0, top: 0 }
  },
})

export default router
