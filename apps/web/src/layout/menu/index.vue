<template>
  <!-- 根层用 el-menu 容器；递归子层只渲染菜单项，不重复套 el-menu -->
  <el-menu
    v-if="isRoot"
    :default-active="$route.path"
    background-color="#001529"
    text-color="#fff"
    active-text-color="#409eff"
  >
    <SideMenu :menu-list="menuList" :parent-path="parentPath" :is-root="false" />
  </el-menu>

  <template v-else>
    <template v-for="item in menuList" :key="item.path">
      <!-- 叶子节点（无 children 或只有一个 child）：直接渲染菜单项 -->
      <el-menu-item
        v-if="isLeaf(item) && !leaf(item).meta?.hidden"
        :index="resolvePath(leaf(item).path)"
        @click="goRoute(resolvePath(leaf(item).path))"
      >
        <el-icon><component :is="getIcon(leaf(item).meta?.icon)" /></el-icon>
        <span>{{ leaf(item).meta?.title }}</span>
      </el-menu-item>

      <!-- 多子节点：渲染可展开子菜单 -->
      <el-sub-menu
        v-else-if="item.children && item.children.length > 1 && item.meta?.title"
        :index="resolvePath(item.path)"
      >
        <template #title>
          <el-icon><component :is="getIcon(item.meta?.icon)" /></el-icon>
          <span>{{ item.meta?.title }}</span>
        </template>
        <SideMenu
          :menu-list="item.children"
          :parent-path="resolvePath(item.path)"
          :is-root="false"
        />
      </el-sub-menu>
    </template>
  </template>
</template>

<script setup lang="ts">
import {
  Files,
  Goods,
  Histogram,
  HomeFilled,
  Lock,
  Menu as MenuIcon,
  Stamp,
  UserFilled,
} from '@element-plus/icons-vue'
import type { Component } from 'vue'
import type { RouteRecordRaw } from 'vue-router'
import { useRouter } from 'vue-router'

// 给递归组件起名字，才能在自己模板里用 <SideMenu> 调用自己
defineOptions({ name: 'SideMenu' })

// 根层默认 true：layout 调用 <Menu :menu-list="..."/> 不传 isRoot 时，自动套 el-menu 容器
const props = withDefaults(
  defineProps<{
    menuList: RouteRecordRaw[]
    parentPath?: string
    isRoot?: boolean
  }>(),
  { isRoot: true, parentPath: '' },
)

const $router = useRouter()

// 叶子判断：没有子节点，或只有一个子节点（这种会“折叠”成菜单项展示）
const isLeaf = (item: RouteRecordRaw) => !item.children || item.children.length <= 1

// 取出真正要展示的那个节点：单子节点时返回它的 child，否则返回自身
const leaf = (item: RouteRecordRaw) => (item.children?.length === 1 ? item.children[0] : item)

// 拼完整路径：绝对路径直接用，相对路径拼上父级路径
const resolvePath = (path: string) =>
  path.startsWith('/') ? path : props.parentPath ? `${props.parentPath}/${path}` : `/${path}`

// 图标映射表：手动 import 再当对象，自动导入扫不到“动态字符串”
const iconMap: Record<string, Component> = {
  HomeFilled,
  Histogram,
  Lock,
  UserFilled,
  Menu: MenuIcon,
  Goods,
  Stamp,
  Files,
}
const getIcon = (icon?: string): Component | undefined => (icon ? iconMap[icon] : undefined)

// 点击菜单项跳转
const goRoute = (path: string) => $router.push(path)
</script>
