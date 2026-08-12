import { defineStore } from 'pinia'
import { ref } from 'vue'
import { client } from '@/api/client'
import type { components, paths } from '@/api/schema'

// 用户信息结构（对应 GET /api/v1/acl/user/info 的 data 字段）
type ResponseUserInfo = components['schemas']['types.ResponseUserInfo']

// 登录入参类型（对应 openapi 里的 POST /api/v1/acl/index/login）
type LoginFormData =
  paths['/api/v1/acl/index/login']['post']['requestBody']['content']['application/json']

export const useUserStore = defineStore('User', () => {
  // state
  const token = ref(localStorage.getItem('token') || '')
  const menuRoutes = ref<string[]>([]) // 菜单路由（后端 menu.code 数组）
  const username = ref('')
  const avatar = ref('')
  const buttons = ref<string[]>([]) // 按钮权限
  const isRoutesAdded = ref(false) // 动态路由是否已挂载（防止重复 addRoute）

  // 重置权限态：登录成功 / 退出时调用，确保用新账号重新拉取并挂载动态路由
  function resetUserState() {
    token.value = ''
    username.value = ''
    avatar.value = ''
    buttons.value = []
    menuRoutes.value = []
    isRoutesAdded.value = false
  }

  // 登录：拿到 token 后存本地（持久化），后续靠 token 换用户信息
  // client 失败（HTTP 异常 / 401 / 业务码非200）直接 throw，调用方 await 时自然捕获
  async function userLogin(data: LoginFormData) {
    const { data: res } = await client.POST('/api/v1/acl/index/login', { body: data })
    if (res?.code !== 200) throw new Error(res?.message || '登录失败')
    const newToken = res.data as string
    resetUserState()
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  // 拿用户信息：靠 token 换 routes/buttons/name/avatar
  async function getUserInfo() {
    const { data: res } = await client.GET('/api/v1/acl/user/info')
    if (res?.code !== 200) throw new Error(res?.message || '获取用户信息失败')
    const data = res.data as ResponseUserInfo | undefined
    if (data) {
      username.value = data.name ?? ''
      avatar.value = data.avatar ?? ''
      buttons.value = data.buttons ?? []
      menuRoutes.value = data.routes ?? [] // 原始 routes（menu.code 数组）
    }
  }

  // 退出登录
  async function userLogout() {
    const { data: res } = await client.POST('/api/v1/acl/index/logout')
    if (res?.code !== 200) throw new Error(res?.message || '退出失败')
    resetUserState()
    localStorage.removeItem('token')
  }

  return { token, menuRoutes, username, avatar, buttons, isRoutesAdded, userLogin, getUserInfo, userLogout }
})
