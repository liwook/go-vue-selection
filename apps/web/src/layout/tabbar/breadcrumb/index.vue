<template>
  <el-breadcrumb separator="/">
    <el-breadcrumb-item v-for="item in matched" :key="item.path">
      <!-- 纯位置指示：统一灰色、不可点击 -->
      <span class="layout-breadcrumb-text">
        <el-icon v-if="item.meta?.icon" class="breadcrumb-icon">
          <component :is="getRouteIcon(item.meta.icon)" />
        </el-icon>
        {{ item.meta.title }}
      </span>
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { getRouteIcon } from '@/utils/icons'

const $route = useRoute()

// route.matched 是当前路由匹配到的所有层级；过滤掉没 title / 隐藏的项
const matched = computed(() => $route.matched.filter((r) => r.meta?.title && !r.meta?.hidden))
</script>

<style scoped lang="scss">
.layout-breadcrumb-text {
  color: #606266;
}

.breadcrumb-icon {
  margin-right: 4px;
}
</style>
