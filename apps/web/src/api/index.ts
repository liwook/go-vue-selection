import { client } from './client'

export { client }

// ============================================================================
// 使用示例（在组件 / store / composable 中调用）：
//
// import { api } from '@/api'
//
// // 登录：参数类型、返回类型均由 schema.d.ts 自动推导，写错字段名会编译报错
// const { data, error } = await api.POST('/api/v1/acl/index/login', {
//   body: { username: 'admin', password: '123456' },
// })
// if (error) {
//   // error 已被 client.use 的 onResponse 统一提示，这里可做额外处理
//   return
// }
// // data 类型为 result.ResponseData，可直接取 data.data.token 等字段
// console.log(data?.data)
//
// // 其他写法：
// await api.GET('/api/v1/acl/index/info')
// await api.DELETE('/api/v1/acl/index/logout')
// ============================================================================

// 这里给出登录接口的显式封装，方便直接 import 使用
export async function login(username: string, password: string) {
  const { data, error } = await client.POST('/api/v1/acl/index/login', {
    body: { username, password },
  })
  if (error) return null
  return data
}
