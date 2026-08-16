<template>
  <el-dropdown>
    <!-- 触发区：头像 + 用户名 -->
    <span class="layout-userinfo-trigger">
      <el-avatar :size="32" class="layout-userinfo-avatar">
        <el-icon><User /></el-icon>
      </el-avatar>
      <span class="layout-userinfo-name">{{ displayName }}</span>
    </span>

    <template #dropdown>
      <el-dropdown-menu>
        <!-- 问候语（置灰不可点） -->
        <el-dropdown-item disabled class="layout-userinfo-greeting">
          {{ greeting }}{{ displayName }}
        </el-dropdown-item>
        <el-dropdown-item :icon="User" @click="goProfile">
          个人信息
        </el-dropdown-item>
        <el-dropdown-item :icon="SwitchButton" divided @click="logout">
          退出登录
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup lang="ts">
import { SwitchButton, User } from '@element-plus/icons-vue'
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()
const router = useRouter()

// 显示名：取用户昵称 name，未登录/接口未通时兜底 admin
const displayName = computed(() => userStore.name ?? 'admin')

// 时间段问候语（5 档，与首页保持一致）
const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '凌晨好'
  if (hour < 12) return '早上好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  return '晚上好'
})

// 个人信息：跳到 layout 下 /profile（hidden: true，不出现在左侧菜单）
const goProfile = () => router.push('/profile')

// 退出登录：清 token/状态后主动跳回登录页
const logout = async () => {
  await userStore.userLogout()
  await router.push('/login')
}
</script>

<style scoped lang="scss">
.layout-userinfo-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  outline: none;
}
.layout-userinfo-avatar {
  background: var(--el-color-primary);
  color: #fff;
}
.layout-userinfo-name {
  font-size: 14px;
  color: #303133;
}
</style>
