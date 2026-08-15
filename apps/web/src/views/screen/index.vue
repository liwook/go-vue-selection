<template>
  <div class="screen">
    <!-- 顶部标题栏 -->
    <header class="screen-header">
      <span class="deco-line" />
      <h1>智游云枢 · 智慧旅游可视化大数据展示平台</h1>
      <span class="deco-line" />
    </header>

    <!-- 主体：左 28% / 中 44% / 右 28% -->
    <main class="screen-body">
      <section class="col col-left">
        <div class="chart-card">
          <div class="big-number">
            <div class="big-number__title">实时游客总数</div>
            <div class="big-number__value">{{ data.realTimeStats.total.toLocaleString() }}</div>
            <div class="big-number__unit">人</div>
          </div>
        </div>
        <div class="chart-card">
          <div class="chart-card__title">实时游客占比</div>
          <EChart :option="liquidOption" />
        </div>
        <div class="chart-card">
          <div class="chart-card__title">男女比例</div>
          <EChart :option="genderOption" />
        </div>
        <div class="chart-card">
          <div class="chart-card__title">年龄分布</div>
          <EChart :option="ageOption" />
        </div>
      </section>
      <section class="col col-center">
        <div class="chart-card">占位：中列图表</div>
      </section>
      <section class="col col-right">
        <div class="chart-card">占位：右列图表</div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import * as echarts from 'echarts'
import 'echarts-liquidfill'
import type { EChartsOption } from 'echarts'
import { computed, onMounted } from 'vue'
import EChart from '@/components/EChart.vue'
import chinaJson from './china.json'
import { useScreenData } from './useScreenData'

const { data } = useScreenData()

onMounted(() => {
  // 注册中国地图（中列地图使用）
  echarts.registerMap('china', chinaJson as Parameters<typeof echarts.registerMap>[1])
})

// 左列：水球图（实时游客占比）
const liquidOption = computed(() => ({
  series: [
    {
      type: 'liquidFill',
      radius: '70%',
      data: [data.value.realTimeStats.trend],
      color: ['#4fc3f7'],
      backgroundStyle: { color: 'transparent', borderColor: 'transparent' },
      outline: { show: false },
      label: {
        formatter: () => `${Math.round(data.value.realTimeStats.trend * 100)}%`,
        fontSize: 28,
        color: '#ffffff',
      },
    },
  ],
}) as EChartsOption)

// 左列：男女比例
const genderOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'item' },
  legend: {
    bottom: 0,
    textStyle: { color: '#cfe4ff' },
    data: ['男', '女'],
  },
  series: [
    {
      type: 'pie',
      radius: ['40%', '65%'],
      center: ['50%', '45%'],
      label: { color: '#cfe4ff' },
      data: [
        { name: '男', value: data.value.genderRatio.male },
        { name: '女', value: data.value.genderRatio.female },
      ],
    },
  ],
}))

// 左列：年龄分布
const ageOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'item' },
  legend: {
    bottom: 0,
    textStyle: { color: '#cfe4ff' },
  },
  series: [
    {
      type: 'pie',
      radius: ['40%', '65%'],
      center: ['50%', '45%'],
      label: { color: '#cfe4ff' },
      data: data.value.ageRatio,
    },
  ],
}))
</script>

<style scoped lang="scss">
.screen {
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background: radial-gradient(ellipse at center, #0a1a3a 0%, #030b1f 100%);
  color: #cfe4ff;
  display: flex;
  flex-direction: column;
}

.screen-header {
  height: 72px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;

  h1 {
    font-size: 26px;
    letter-spacing: 4px;
    margin: 0;
    background: linear-gradient(90deg, #4fc3f7, #ffffff);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
  }

  .deco-line {
    width: 120px;
    height: 2px;
    background: linear-gradient(90deg, transparent, #4fc3f7);
  }
}

.screen-body {
  flex: 1;
  display: grid;
  grid-template-columns: 28% 44% 28%;
  gap: 16px;
  padding: 16px;
  box-sizing: border-box;
}

.col {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 0;
}

.chart-card {
  flex: 1;
  border: 1px solid rgba(79, 195, 247, 0.25);
  border-radius: 8px;
  background: rgba(10, 30, 60, 0.4);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 12px;
  box-sizing: border-box;
  min-height: 0;
}

.chart-card__title {
  font-size: 15px;
  color: #9fd3ff;
  margin-bottom: 8px;
  align-self: flex-start;
}

.chart-card :deep(.e-chart) {
  flex: 1;
  min-height: 0;
}

.big-number {
  text-align: center;

  &__title {
    font-size: 15px;
    color: #9fd3ff;
    margin-bottom: 12px;
  }

  &__value {
    font-size: 42px;
    font-weight: 700;
    color: #ffd166;
    text-shadow: 0 0 12px rgba(255, 209, 102, 0.6);
    line-height: 1.1;
  }

  &__unit {
    font-size: 14px;
    color: #cfe4ff;
    margin-top: 8px;
  }
}
</style>
