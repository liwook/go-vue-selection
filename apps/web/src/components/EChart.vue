<template>
  <div ref="el" class="e-chart"></div>
</template>

<script setup lang="ts">
import type { EChartsOption } from 'echarts'
import * as echarts from 'echarts'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps<{
  option: EChartsOption
}>()

const el = ref<HTMLDivElement>()
let chart: echarts.ECharts | null = null
let observer: ResizeObserver | null = null

onMounted(() => {
  if (!el.value) return
  chart = echarts.init(el.value)
  chart.setOption(props.option)

  // 用 ResizeObserver 观察图表容器自身尺寸变化（列宽因 flex 重排时也能即时自适应）
  observer = new ResizeObserver(() => {
    chart?.resize()
  })
  observer.observe(el.value)
})

// option 变化（含响应式数据更新）时重绘
watch(
  () => props.option,
  (opt) => {
    chart?.setOption(opt)
  },
  { deep: true },
)

onBeforeUnmount(() => {
  observer?.disconnect()
  observer = null
  chart?.dispose()
  chart = null
})
</script>

<style scoped>
.e-chart {
  width: 100%;
  height: 100%;
}
</style>
