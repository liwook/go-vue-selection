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
// Vue Router 5 推荐使用 `return` 式返回导航结果，取代已废弃的 `next()` 回调式
router.beforeEach(async (to) => {
  const userStore = useUserStore() // pinia 已在 main.ts 中注册（顺序在 router 之前），可直接使用
  const token = userStore.token

  // ① 已登录：放行
  if (token) {
    // 目标本身就是白名单（登录页 / 404），直接放行，不触发任何用户信息请求
    // 避免「localStorage 残留旧 token + 后端 500 → 跳首页失败 → 跳回登录页」的循环
    if (whiteList.includes(to.path)) {
      return true
    }

    try {
      // 同一个 session 内只拉一次用户信息并挂载一次动态路由
      if (!userStore.isRoutesAdded) {
        await userStore.getUserInfo()

        // 按后端 routes 过滤异步路由
        const userAsyncRoute = filterAsyncRoute(asyncRoutes, userStore.menuRoutes as string[])

        // 逐个 addRoute（动态挂载）
        userAsyncRoute.forEach((route: RouteRecordRaw) => {
          router.addRoute(route)
        })
        // 最后挂 anyRoute 做 404 兜底（必须在动态路由之后）
        router.addRoute(anyRoute)

        userStore.isRoutesAdded = true

        // 重新进入目标路由（确保 addRoute 生效后再放行）
        return { ...to, replace: true }
      }

      return true
    } catch (_error) {
      // 获取用户信息失败（token 失效等）→ 清登录态跳登录
      // 注意：若拦截器已对鉴权错误兜底清过登录态，这里 clearAuth 幂等无副作用。
      // 对于非鉴权错误（如 500），也清登录态回登录页，避免卡在受限页面。
      userStore.clearAuth()
      return { path: '/login' }
    }
  }
  // ② 未登录：白名单放行，否则去登录页
  if (whiteList.includes(to.path)) {
    return true
  }

  return { path: '/login' }
})

export default router
