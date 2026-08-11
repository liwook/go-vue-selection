<template>
  <div class="layout_container">
    <!-- 左侧菜单 -->
    <div class="layout_slider">
      <Logo />
      <Menu :menu-list="menuList" />
    </div>

    <!-- 右侧区域：顶部导航 + 内容区 -->
    <div class="layout_right">
      <div class="layout_tabbar">
        顶部导航（占位）
      </div>
      <div class="layout_main">
        <Main />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import { asyncRoutes, constantRoutes } from '@/router/routes'
import { useUserStore } from '@/store/modules/user'
import { filterAsyncRoute } from '@/utils/routeFilter'
import Logo from './logo/index.vue'
import Main from './main/index.vue'
import Menu from './menu/index.vue'

const userStore = useUserStore()

// 菜单取：常量路由 + 按当前用户权限过滤后的异步路由
const menuList = computed<RouteRecordRaw[]>(() => [
  ...constantRoutes,
  ...filterAsyncRoute(asyncRoutes, userStore.menuRoutes as string[]),
])
</script>

<style scoped lang="scss">
.layout_container {
  width: 100%;
  height: 100%;
  display: flex;

  .layout_slider {
    width: $base-menu-width;
    height: 100%;
    background: #001529;
  }

  .layout_right {
    flex: 1;
    display: flex;
    flex-direction: column;

    .layout_tabbar {
      height: $base-tabbar-height;
      background: skyblue;
    }

    .layout_main {
      flex: 1;
      padding: 20px;
      overflow: auto;
    }
  }
}
</style>
