<template>
  <div class="layout-tabbar">
    <Breadcrumb class="layout-tabbar-left" />
    <div class="layout-tabbar-right">
      <Setting />
      <Userinfo />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElNotification } from 'element-plus'
import { computed, onMounted } from 'vue'
import { useUserStore } from '@/store/modules/user'
import Breadcrumb from './breadcrumb/index.vue'
import Setting from './setting/index.vue'
import Userinfo from './userinfo/index.vue'

const userStore = useUserStore()

const displayName = computed(() => userStore.name ?? 'admin')

// 按当前小时计算问候语（与下拉里的逻辑一致）
const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '凌晨好'
  if (hour < 12) return '早上好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  return '晚上好'
})

// 模块级标志：整页生命周期内只弹一次（避免 HMR / 重复挂载重复弹）
let notified = false
onMounted(() => {
  if (notified) return
  notified = true
  ElNotification({
    title: `Hi, ${greeting.value}!`,
    message: `欢迎回来，${displayName.value}`,
    type: 'success',
    duration: 3000,
    offset: 60,
  })
})
</script>

<style scoped lang="scss">
.layout-tabbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding: 0 20px;
}
.layout-tabbar-right {
  display: flex;
  align-items: center;
  gap: 20px;
}
</style>
