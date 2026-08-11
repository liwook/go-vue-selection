import type { RouteRecordRaw } from 'vue-router'
import { createRouter, createWebHashHistory } from 'vue-router'
import { useUserStore } from '@/store/modules/user'
import { filterAsyncRoute } from '@/utils/routeFilter'
import { anyRoute, asyncRoutes, constantRoutes } from './routes'

const router = createRouter({
  // 哈希模式（部署简单，不需要服务器配置）
  history: createWebHashHistory(),
  routes: constantRoutes,
})

// 白名单：免登录也能访问的页面
const whiteList = ['/login', '/404']

// 路由前置守卫（全局）：负责鉴权 + 动态挂载权限路由
router.beforeEach(async (to, _from, next) => {
  const userStore = useUserStore() // pinia 已在 main.ts 中注册（顺序在 router 之前），可直接使用
  const token = userStore.token

  // ① 已登录：放行
  if (token) {
    if (to.path === '/login') {
      next({ path: '/' }) // 已登录还去登录页？直接打回首页
    } else {
      try {
        // 已拉过用户信息（menuRoutes 非空）则跳过，避免每次跳转重复请求
        if (userStore.menuRoutes.length === 0) {
          await userStore.getUserInfo()
        }

        // 按后端 routes 过滤异步路由
        const userAsyncRoute = filterAsyncRoute(asyncRoutes, userStore.menuRoutes as string[])

        // 逐个 addRoute（动态挂载）
        userAsyncRoute.forEach((route: RouteRecordRaw) => {
          router.addRoute(route)
        })
        // 最后挂 anyRoute 做 404 兜底（必须在动态路由之后）
        router.addRoute(anyRoute)

        // 重新进入目标路由（确保 addRoute 生效后再放行）
        next({ ...to, replace: true })
      } catch (error) {
        // 获取用户信息失败（token 失效等）→ 清登录态跳登录
        userStore.token = ''
        localStorage.removeItem('token')
        next({ path: '/login' })
      }
    }
  }
  // ② 未登录：白名单放行，否则去登录页
  else {
    if (whiteList.includes(to.path)) {
      next()
    } else {
      next({ path: '/login' })
    }
  }
})

export default router
