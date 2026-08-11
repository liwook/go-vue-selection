import { ElMessage } from 'element-plus'
import createClient from 'openapi-fetch'
import type { paths } from './schema'
import 'element-plus/es/components/message/style/css'

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
      ElMessage({ type: 'error', message: data.message ?? '请求失败' })
      // 抛错以便调用方在 catch 中感知（openapi-fetch 默认不会因业务 code 失败）
      throw new Error(data.message ?? '请求失败')
    }
  },
})
