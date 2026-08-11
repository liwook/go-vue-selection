import { defineStore } from 'pinia'
import { client } from '@/api/client'
import type { components, paths } from '@/api/schema'

// 用户信息结构（对应 GET /api/v1/acl/user/info 的 data 字段）
type ResponseUserInfo = components['schemas']['types.ResponseUserInfo']

// 登录入参类型（对应 openapi 里的 POST /api/v1/acl/index/login）
type LoginFormData =
  paths['/api/v1/acl/index/login']['post']['requestBody']['content']['application/json']

export const useUserStore = defineStore('User', {
  state: () => {
    return {
      token: localStorage.getItem('token') || '',
      menuRoutes: [] as string[], // 菜单路由（store 计算好，给 layout 用）
      username: '',
      avatar: '',
      buttons: [] as string[], // 按钮权限
    }
  },
  actions: {
    // 登录：拿到 token 后存本地（持久化），后续靠 token 换用户信息
    async userLogin(data: LoginFormData) {
      try {
        const { data: res } = await client.POST('/api/v1/acl/index/login', { body: data })
        if (res?.code === 200) {
          const token = res.data as string
          this.token = token
          localStorage.setItem('token', token)
        } else {
          // 业务失败也 throw，让调用方 catch 到
          throw new Error(res?.message || '登录失败')
        }
      } catch (error) {
        // client 已 throw：HTTP 异常 / 401 / 业务码非200
        // 这里统一转成 rejected，交给登录页提示
        return Promise.reject(error)
      }
    },
    // 拿用户信息：靠 token 换 routes/buttons/name/avatar
    async getUserInfo() {
      try {
        const { data: res } = await client.GET('/api/v1/acl/user/info')
        if (res?.code === 200) {
          const data = res.data as ResponseUserInfo | undefined
          if (data) {
            this.username = data.name ?? ''
            this.avatar = data.avatar ?? ''
            this.buttons = data.buttons ?? []
            this.menuRoutes = data.routes ?? [] // 原始 routes（menu.code 数组）
          }
        } else {
          throw new Error(res?.message || '获取用户信息失败')
        }
      } catch (error) {
        return Promise.reject(error)
      }
    },
    // 退出登录
    async userLogout() {
      try {
        const { data: res } = await client.POST('/api/v1/acl/index/logout')
        if (res?.code === 200) {
          this.token = ''
          this.username = ''
          this.avatar = ''
          this.buttons = []
          this.menuRoutes = []
          localStorage.removeItem('token')
        } else {
          throw new Error(res?.message || '退出失败')
        }
      } catch (error) {
        return Promise.reject(error)
      }
    },
  },
  getters: {},
})
