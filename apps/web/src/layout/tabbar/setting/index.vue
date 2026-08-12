<template>
  <div class="layout-tabbar-setting">
    <!-- 全屏：原生 Fullscreen API，不引入 screenfull -->
    <el-tooltip content="全屏" placement="bottom">
      <el-icon class="layout-setting-icon" @click="toggleFullscreen">
        <FullScreen />
      </el-icon>
    </el-tooltip>

    <!-- 暗黑模式：切 html.dark 类 -->
    <el-tooltip :content="isDark ? '亮色模式' : '暗黑模式'" placement="bottom">
      <el-icon class="layout-setting-icon" @click="toggleDark">
        <component :is="isDark ? Sunny : Moon" />
      </el-icon>
    </el-tooltip>

    <!-- 主题色：选色后写入 --el-color-primary 变体 -->
    <el-color-picker
      :model-value="primaryColor"
      :predefine="predefineColors"
      @change="changePrimary"
    />
  </div>
</template>

<script setup lang="ts">
import { FullScreen, Moon, Sunny } from '@element-plus/icons-vue'
import { ref } from 'vue'
import { setPrimaryColor } from '@/utils/theme'

// 用当前 DOM 上的 dark 类初始化，避免刷新后开关状态与页面不同步（可选接 localStorage 持久化）
const isDark = ref(document.documentElement.classList.contains('dark'))
const primaryColor = ref('#409EFF')

// 主题色预设（点开选择器时快速选）
const predefineColors = ['#409EFF', '#1677ff', '#7265e6', '#28c76f', '#f5222d', '#fa541c']

const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    document.exitFullscreen()
  }
}

const toggleDark = () => {
  isDark.value = !isDark.value
  // 第二个参数强制设置类是否存在（而不是单纯 toggle），状态更可控
  document.documentElement.classList.toggle('dark', isDark.value)
}

const changePrimary = (color: string | null) => {
  if (!color) return
  primaryColor.value = color
  setPrimaryColor(color)
}
</script>

<style scoped lang="scss">
.layout-tabbar-setting {
  display: flex;
  align-items: center;
  gap: 18px;
}
.layout-setting-icon {
  font-size: 20px;
  cursor: pointer;
  color: #606266;
  transition: color 0.2s;
  &:hover {
    color: var(--el-color-primary);
  }
}
</style>
