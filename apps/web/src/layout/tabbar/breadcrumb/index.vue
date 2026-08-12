<template>
  <el-breadcrumb separator="/">
    <el-breadcrumb-item v-for="item in matched" :key="item.path">
      <!-- 首页项可点，其余只展示文字；每一级都带图标 -->
      <span v-if="item.meta?.title === '首页'" class="layout-breadcrumb-link" @click="goHome">
        <el-icon v-if="item.meta?.icon" class="breadcrumb-icon">
          <component :is="getRouteIcon(item.meta.icon)" />
        </el-icon>
        {{ item.meta.title }}
      </span>
      <span v-else class="layout-breadcrumb-text">
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
import { useRoute, useRouter } from 'vue-router'
import { getRouteIcon } from '@/utils/icons'

const $route = useRoute()
const $router = useRouter()

// route.matched 是当前路由匹配到的所有层级；过滤掉没 title / 隐藏的项
const matched = computed(() => $route.matched.filter((r) => r.meta?.title && !r.meta?.hidden))

const goHome = () => $router.push('/home')
</script>

<style scoped lang="scss">
.layout-breadcrumb-link {
  color: #409eff;
  cursor: pointer;
  font-weight: 500;
  &:hover {
    text-decoration: underline;
  }
}
.layout-breadcrumb-text {
  color: #606266;
}

.breadcrumb-icon {
  margin-right: 4px;
}
</style>
