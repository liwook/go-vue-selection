<template>
  <el-breadcrumb separator="/">
    <el-breadcrumb-item v-for="(item, index) in matched" :key="item.path">
      <!-- 最后一项为当前页：灰色、不可点击；其余层级可点击回退 -->
      <span
        v-if="index < matched.length - 1"
        class="layout-breadcrumb-link"
        @click="goTo(item)"
      >
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

// 点击前序层级跳转：优先用 redirect，否则跳首个子路由（拼接父 path），再否则用当前 path
const goTo = (route: (typeof matched.value)[number]) => {
  if (typeof route.redirect === 'string') {
    $router.push(route.redirect)
    return
  }
  const firstChild = route.children?.[0]?.path
  if (firstChild) {
    // 子路由 path 多为相对路径，需拼上父级 path，绝对路径则原样使用
    const target = firstChild.startsWith('/') ? firstChild : `${route.path}/${firstChild}`.replace('//', '/')
    $router.push(target)
    return
  }
  $router.push(route.path)
}
</script>

<style scoped lang="scss">
.layout-breadcrumb-link {
  color: var(--el-color-primary);
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
