<template>
  <div class="profile-wrapper">
    <el-card class="profile-card" shadow="hover">
      <div class="header">
        <el-avatar :size="96" :src="avatarUrl" />
        <div class="title">
          <h2>{{ userStore.name || '未登录用户' }}</h2>
          <p>个人信息</p>
        </div>
      </div>

      <el-descriptions :column="1" border class="profile-desc">
        <el-descriptions-item label="用户名">
          {{ userStore.name || '—' }}
        </el-descriptions-item>
        <el-descriptions-item label="头像">
          <el-avatar :size="40" :src="avatarUrl" />
        </el-descriptions-item>
        <el-descriptions-item label="菜单权限">
          共 {{ userStore.menuRoutes.length }} 项
        </el-descriptions-item>
        <el-descriptions-item label="按钮权限">
          共 {{ userStore.buttons.length }} 项
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()
const { avatar } = storeToRefs(userStore)

// 默认头像兜底（与首页保持一致）
const defaultAvatar = `${import.meta.env.BASE_URL}chopper.jpeg`
const avatarUrl = computed(() => avatar.value || defaultAvatar)
</script>

<style scoped lang="scss">
.profile-wrapper {
  padding: 16px;
}
.profile-card {
  max-width: 720px;
  margin: 0 auto;
}
.header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--el-border-color-lighter);
  margin-bottom: 16px;
}
.title {
  h2 {
    margin: 0 0 4px;
    font-size: 20px;
    color: var(--el-text-color-primary);
  }
  p {
    margin: 0;
    color: var(--el-text-color-secondary);
    font-size: 13px;
  }
}
.profile-desc {
  margin-top: 8px;
}
</style>
