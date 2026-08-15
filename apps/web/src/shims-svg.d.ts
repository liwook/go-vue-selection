// 手写 shim：为 unplugin-svg-component 生成的虚拟模块提供类型兜底。
// 该插件会在 `vite build` 启动阶段生成 src/svg-component.d.ts，但该文件
// 属于自动生成产物（见 .gitignore），不入库。为了让 `vue-tsc -b` 在
// `vite build` 之前能通过类型检查，这里提供一份等价的模块声明。
declare module '~virtual/svg-component' {
  import type { DefineComponent } from 'vue'

  const SvgIcon: DefineComponent<{
    name: string
    [key: string]: unknown
  }>

  export default SvgIcon
}
