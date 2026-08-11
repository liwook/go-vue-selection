import { client } from './client'

export { client }

// ============================================================================
// 使用示例（在组件 / store / composable 中调用）：
//
// import { client } from '@/api'
//
// // 登录：参数类型、返回类型均由 schema.d.ts 自动推导，写错字段名会编译报错
// const { data, error } = await client.POST('/api/v1/acl/index/login', {
//   body: { username: 'admin', password: '123456' },
// })
//
// 注意：本项目约定「拦截器失败即 throw」，业务层（store / 组件）用 try/catch 感知，
// 不再依赖 openapi-fetch 的 error 返回值分支。登录相关逻辑统一在 useUserStore 中。
// ============================================================================
