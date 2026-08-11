import { ElMessage } from 'element-plus'
import createClient from 'openapi-fetch'
import type { paths } from './schema'
import 'element-plus/es/components/message/style/css'

// 与 utils/request.ts 保持一致的 baseURL，走 vite 的 /api 代理
export const client = createClient<paths>({
  baseUrl: '/api',
})

// 统一的业务错误码拦截：与 utils/request.ts 的响应拦截逻辑保持一致。
// 后端在 HTTP 200 的 body 里用 code 字段表达业务结果（200=成功，其余为错误）。
client.use({
  onResponse: async ({ response }) => {
    // 仅处理 JSON 响应；非 2xx 由 openapi-fetch 抛错，交由调用方 catch。
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

export default client
