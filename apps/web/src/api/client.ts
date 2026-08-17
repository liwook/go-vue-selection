import { ElMessage } from 'element-plus'
import createClient from 'openapi-fetch'
import { useUserStore } from '@/store/modules/user'
import type { paths } from './schema'
import 'element-plus/es/components/message/style/css'

// 鉴权类业务码：需要登录 / 无效 Token（与后端 pkg/result/code.go 保持一致）
const AUTH_ERROR_CODES = [110002, 110003]

// 跳转到登录页（避免在拦截器里直接依赖 router，用 location 兜底；hash 模式无需刷新）
function redirectToLogin() {
  if (location.hash !== '#/login' && location.pathname !== '/login') {
    location.hash = '#/login'
  }
}

// 当前是否在登录页（或正在跳往登录页）—— 用于鉴权类错误时抑制重复弹窗
function isAtLogin() {
  return location.hash.startsWith('#/login') || location.pathname === '/login'
}

// 走 vite 的 /api 代理（见 vite.config 的 server.proxy）
// baseUrl 留空：请求路径写完整后端路径 /api/v1/...，由 vite 代理原样转发，
// 同时和 schema.d.ts 里的 paths 键（/api/v1/...）完全一致，类型才能匹配。
export const client = createClient<paths>({
  baseUrl: '',
})

// 统一的请求/响应拦截（中间件）
client.use({
  // ① 请求前注入 token（Bearer）
  onRequest({ request }) {
    const token = localStorage.getItem('token')
    if (token) {
      request.headers.set('Authorization', `Bearer ${token}`)
    }
  },

  // ② 响应拦截：区分 HTTP 错误与业务码错误
  onResponse: async ({ response }) => {
    // HTTP 层错误（4xx/5xx）
    if (!response.ok) {
      if (response.status === 401) {
        ElMessage({ type: 'warning', message: '登录已过期，请重新登录' })
      } else {
        ElMessage({ type: 'error', message: `请求失败（${response.status}）` })
      }
      throw new Error(`请求失败（${response.status}）`)
    }

    // 业务码错误（HTTP 200 但 body.code !== 200）
    const contentType = response.headers.get('content-type') ?? ''
    if (!contentType.includes('application/json')) return

    const data = (await response.clone().json()) as { code?: number; message?: string }
    if (typeof data.code === 'number' && data.code !== 200) {
      if (AUTH_ERROR_CODES.includes(data.code)) {
        // 鉴权类错误：全局兜底——清登录态 + 跳登录页
        useUserStore().clearAuth()
        const wasAtLogin = isAtLogin()
        redirectToLogin()
        // 仅当用户「在业务页主动操作」触发失效时才提示一次；
        // 若是「跳往/已处于登录页」则静默，避免「打开登录页就弹窗」的糟糕体验。
        if (!wasAtLogin) {
          ElMessage({ type: 'warning', message: '登录已过期，请重新登录' })
        }
      } else {
        ElMessage({ type: 'error', message: data.message ?? '请求失败' })
      }
      // 抛错以便调用方在 catch 中感知（openapi-fetch 默认不会因业务 code 失败）
      throw new Error(data.message ?? '请求失败')
    }
  },
})
