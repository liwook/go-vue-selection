// ECharts 统一封装入口（按需引入）。
// 水球图已改用自绘 SVG/CSS 组件（components/LiquidBall.vue），不再依赖 echarts-liquidfill，
// 因此可以走 echarts/core 按需引入，仅打包大屏实际用到图表类型，显著降低 vendor 体积。

import type { EChartsOption } from 'echarts'
import { BarChart, LineChart, MapChart, PieChart } from 'echarts/charts'
import {
  GridComponent,
  LegendComponent,
  TooltipComponent,
  VisualMapComponent,
} from 'echarts/components'
import * as echarts from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import type { TopLevelFormatterParams } from 'echarts/types/dist/shared'

echarts.use([
  PieChart,
  LineChart,
  BarChart,
  MapChart,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  VisualMapComponent,
  CanvasRenderer,
])

export default echarts
export type { EChartsOption, TopLevelFormatterParams }
