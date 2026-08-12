<template>
  <div class="layout_container">
    <!-- 左侧菜单 -->
    <div class="layout_slider">
      <Logo />
      <Menu :menu-list="menuList" />
    </div>

    <!-- 右侧区域：顶部导航 + 内容区 -->
    <div class="layout_right">
      <!-- 顶栏：原来是占位文本，现在换成 Tabbar 组件 -->
      <Tabbar class="layout_tabbar" />
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
import Tabbar from './tabbar/index.vue' // ← 新增

const userStore = useUserStore()

// 菜单取：常量路由 + 按当前用户权限过滤后的异步路由（menuRoutes 是 string[]，需先过滤）
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
      background: #fff;
      box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08); // 底部细分割线，更精致
    }

    .layout_main {
      flex: 1;
      padding: 20px;
      overflow: auto;
    }
  }
}

// 暗黑模式：Element Plus 只管自身组件，自定义背景需手动覆盖
html.dark {
  .layout_slider {
    background: #141414;
  }
  .layout_tabbar {
    background: #141414;
    border-bottom: 1px solid #303030;
  }
}
</style>
