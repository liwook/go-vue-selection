import type { Component } from 'vue'
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

// 图标映射表：手动 import 再当对象，自动导入扫不到“动态字符串”
export const routeIconMap: Record<string, Component> = {
  HomeFilled,
  Histogram,
  Lock,
  UserFilled,
  Menu: MenuIcon,
  Goods,
  Stamp,
  Files,
}

export const getRouteIcon = (icon?: string): Component | undefined =>
  icon ? routeIconMap[icon] : undefined
