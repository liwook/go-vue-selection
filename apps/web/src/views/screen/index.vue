<template>
  <div class="screen">
    <!-- 顶部标题栏 -->
    <header class="screen-header">
      <span class="deco-line" />
      <h1>智游云枢 · 智慧旅游可视化大数据展示平台</h1>
      <el-button class="screen-header__back" type="primary" plain @click="goHome">
        返回首页
      </el-button>
      <span class="deco-line" />
    </header>

    <!-- 主体：左 28% / 中 44% / 右 28% -->
    <main class="screen-body">
      <section class="col col-left">
        <div class="chart-card">
          <div class="big-number">
            <div class="big-number__title">
              实时游客总数
            </div>
            <div class="big-number__value">
              {{ data.realTimeStats.total.toLocaleString() }}
            </div>
            <div class="big-number__unit">
              人
            </div>
          </div>
        </div>
        <div class="chart-card">
          <div class="chart-card__title">
            实时游客占比
          </div>
          <EChart :option="liquidOption" />
        </div>
        <div class="chart-card">
          <div class="chart-card__title">
            男女比例
          </div>
          <EChart :option="genderOption" />
        </div>
        <div class="chart-card">
          <div class="chart-card__title">
            年龄分布
          </div>
          <EChart :option="ageOption" />
        </div>
      </section>
      <section class="col col-center">
        <div class="chart-card map-card">
          <div class="chart-card__title">
            平台高峰预警信息
          </div>
          <EChart :option="mapOption" />
        </div>
        <div class="chart-card">
          <div class="chart-card__title">
            未来 30 天游客量趋势
          </div>
          <EChart :option="trendOption" />
        </div>
      </section>
      <section class="col col-right">
        <div class="chart-card">
          <div class="chart-card__title">
            热门景区游客量 TOP5
          </div>
          <EChart :option="spotOption" />
        </div>
        <div class="chart-card">
          <div class="chart-card__title">
            年度游客量对比
          </div>
          <EChart :option="yearOption" />
        </div>
        <div class="chart-card">
          <div class="chart-card__title">
            预约渠道分布
          </div>
          <EChart :option="channelOption" />
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import * as echarts from 'echarts'
import 'echarts-liquidfill'
import type { EChartsOption } from 'echarts'
import type { TopLevelFormatterParams } from 'echarts/types/dist/shared'
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import EChart from '@/components/EChart.vue'
import chinaJson from './china.json'
import { useScreenData } from './useScreenData'

// 注册中国地图：必须放在模块顶层（import 后立即执行），
// 确保在子组件 EChart 的 onMounted 调用 setOption 之前地图已注册，否则地图空白。
echarts.registerMap('china', chinaJson as Parameters<typeof echarts.registerMap>[1])

const { data } = useScreenData()
const $router = useRouter()

// 返回首页（走 router 导航，确保视图正确切换）
const goHome = () => $router.push('/home')

// 左列：水球图（实时游客占比）
const liquidOption = computed(
  () =>
    ({
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
    }) as EChartsOption,
)

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

// 中列：中国地图（平台高峰预警）
const mapOption = computed<EChartsOption>(() => ({
  tooltip: {
    trigger: 'item',
    // 自定义 formatter，避免未配置预警的省份悬停时显示 NaN
    formatter: (params: TopLevelFormatterParams) => {
      const item = Array.isArray(params) ? params[0] : params
      const name = item?.name ?? ''
      const raw = item?.value
      const v = typeof raw === 'number' && !Number.isNaN(raw) ? raw : null
      return v === null ? `${name}<br/>预警值：暂无` : `${name}<br/>预警值：${v}`
    },
  },
  visualMap: {
    min: 0,
    max: 100,
    right: 16,
    bottom: 16,
    orient: 'horizontal',
    text: ['高', '低'],
    textStyle: { color: '#cfe4ff' },
    calculable: true,
    inRange: { color: ['#0a2a5e', '#4fc3f7', '#ff8a65'] },
  },
  series: [
    {
      type: 'map',
      map: 'china',
      roam: false,
      // 放大并居中显示，避免地图区域显小
      layoutCenter: ['50%', '50%'],
      layoutSize: '120%',
      label: { show: false },
      itemStyle: {
        areaColor: '#0a1a3a',
        borderColor: 'rgba(79, 195, 247, 0.4)',
      },
      emphasis: {
        label: { color: '#fff' },
        itemStyle: { areaColor: '#4fc3f7' },
      },
      data: data.value.peakWarning,
    },
  ],
}))

// 中列：未来 30 天游客量趋势（渐变面积图）
const trendOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: 40, right: 20, top: 20, bottom: 30 },
  xAxis: {
    type: 'category',
    data: data.value.monthlyTrend.dates,
    axisLine: { lineStyle: { color: 'rgba(207,228,255,0.4)' } },
    axisLabel: { color: '#cfe4ff', interval: 4 },
  },
  yAxis: {
    type: 'value',
    axisLine: { lineStyle: { color: 'rgba(207,228,255,0.4)' } },
    axisLabel: { color: '#cfe4ff' },
    splitLine: { lineStyle: { color: 'rgba(79,195,247,0.1)' } },
  },
  series: [
    {
      type: 'line',
      smooth: true,
      symbol: 'none',
      data: data.value.monthlyTrend.values,
      lineStyle: { color: '#4fc3f7', width: 2 },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(79,195,247,0.5)' },
            { offset: 1, color: 'rgba(79,195,247,0.02)' },
          ],
        },
      },
    },
  ],
}))

// 右列：热门景区游客量 TOP5（横向柱状）
const spotOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  grid: { left: 80, right: 30, top: 20, bottom: 20 },
  xAxis: {
    type: 'value',
    axisLabel: { color: '#cfe4ff' },
    splitLine: { lineStyle: { color: 'rgba(79,195,247,0.1)' } },
  },
  yAxis: {
    type: 'category',
    data: data.value.hotScenicSpots.map((i) => i.name).reverse(),
    axisLine: { lineStyle: { color: 'rgba(207,228,255,0.4)' } },
    axisLabel: { color: '#cfe4ff' },
  },
  series: [
    {
      type: 'bar',
      data: data.value.hotScenicSpots.map((i) => i.value).reverse(),
      barWidth: '55%',
      itemStyle: {
        borderRadius: [0, 4, 4, 0],
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 1,
          y2: 0,
          colorStops: [
            { offset: 0, color: '#1e88e5' },
            { offset: 1, color: '#4fc3f7' },
          ],
        },
      },
      label: { show: true, position: 'right', color: '#cfe4ff' },
    },
  ],
}))

// 右列：年度游客量对比（多系列折线）
const yearOption = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'axis' },
  legend: {
    data: data.value.yearlyComparison.years,
    textStyle: { color: '#cfe4ff' },
    top: 0,
  },
  grid: { left: 40, right: 20, top: 30, bottom: 25 },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: data.value.yearlyComparison.months,
    axisLine: { lineStyle: { color: 'rgba(207,228,255,0.4)' } },
    axisLabel: { color: '#cfe4ff', interval: 1 },
  },
  yAxis: {
    type: 'value',
    axisLabel: { color: '#cfe4ff' },
    splitLine: { lineStyle: { color: 'rgba(79,195,247,0.1)' } },
  },
  series: data.value.yearlyComparison.years.map((name, idx) => ({
    name,
    type: 'line',
    smooth: true,
    data: data.value.yearlyComparison.data[idx],
    lineStyle: { width: 2 },
  })),
}))

// 右列：预约渠道分布（环形饼图）
const channelOption = computed<EChartsOption>(() => ({
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
      data: data.value.channelStats,
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
  position: relative;
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

    // 为左上角的"返回首页"按钮预留空间
    &:first-of-type {
      margin-left: 110px;
    }
  }

  &__back {
    position: absolute;
    left: 24px;
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

.map-card {
  flex: 1.4;
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
