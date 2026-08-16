// ECharts 统一封装入口。
// 注意：echarts-liquidfill 内部依赖 echarts/lib/echarts（全量包），会强制把整库打进 bundle，
// 因此此处使用全量 echarts（按需引入 core 模式反而会因双实例叠加体积更大）。
// 体积优化改由 vite 拆 vendor chunk + preload 解决（见 vite.config.ts 的 manualChunks / build.rollupOptions）。
import * as echarts from 'echarts'
import 'echarts-liquidfill'
import type { EChartsOption } from 'echarts'
import type { TopLevelFormatterParams } from 'echarts/types/dist/shared'

export default echarts
export type { EChartsOption, TopLevelFormatterParams }
