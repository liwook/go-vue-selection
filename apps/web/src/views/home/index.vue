<template>
  <el-card class="home-card" shadow="hover">
    <div class="greeting">
      <el-avatar :size="96" :src="avatarUrl" />
      <div class="text">
        <h2>{{ greeting }}，{{ userStore.name || '游客' }}</h2>
        <p>欢迎使用选课管理系统</p>
      </div>
    </div>
  </el-card>
  <el-card class="welcome-card" shadow="hover" :body-style="{ padding: '0' }">
    <SvgIcon name="welcome" class="welcome-illustration" />
  </el-card>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()
// 用 storeToRefs 保持响应式（头像走 computed 兜底，用户名直接读 userStore.name）
const { avatar } = storeToRefs(userStore)

// 默认头像兜底（public/chopper.jpeg）
const defaultAvatar = `${import.meta.env.BASE_URL}chopper.jpeg`
const avatarUrl = computed(() => avatar.value || defaultAvatar)

// 按当前小时生成问候语
const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 6) return '凌晨好'
  if (h < 12) return '早上好'
  if (h < 14) return '中午好'
  if (h < 18) return '下午好'
  return '晚上好'
})
</script>

<style scoped lang="scss">
.home-card {
  max-width: 100%;
  margin: 24px;
}
.welcome-card {
  max-width: 720px;
  margin: 16px auto 24px;
}
.greeting {
  display: flex;
  align-items: center;
  gap: 24px;
}
.text h2 {
  margin: 0 0 8px;
  font-size: 26px;
}
.text p {
  margin: 0;
  font-size: 16px;
  color: var(--el-text-color-secondary);
}
.welcome-illustration {
  display: block;
  width: 100% !important;
  max-width: 720px;
  height: auto !important;
  aspect-ratio: 1024 / 626;
  margin: 0 auto;
}

</style>
