// 主题色工具：把选中的主色写入根元素的 CSS 变量，驱动 Element Plus 整体换肤
// 用法：setPrimaryColor('#1677ff')

const PRIMARY_COLOR_KEY = 'go-vue-selection-primary-color'

// 统一解析任意颜色为 rgb 元组；支持 #fff / #ffffff / rgb(...) / rgba(...)
export function parseColor(color: string): [number, number, number] {
  const trimmed = color.trim()

  // hex
  if (trimmed.startsWith('#')) {
    let normalized = trimmed.slice(1)
    if (normalized.length === 3) {
      normalized = normalized
        .split('')
        .map((c) => c + c)
        .join('')
    }
    const num = parseInt(normalized, 16)
    return [(num >> 16) & 255, (num >> 8) & 255, num & 255]
  }

  // rgb / rgba
  const rgbMatch = trimmed.match(/rgba?\s*\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)/)
  if (rgbMatch) {
    return [parseInt(rgbMatch[1], 10), parseInt(rgbMatch[2], 10), parseInt(rgbMatch[3], 10)]
  }

  throw new Error(`Unsupported color format: ${color}`)
}

// 把 rgb 三元组转回 hex 字符串
export function rgbToHex(r: number, g: number, b: number): string {
  const toHex = (n: number) => n.toString(16).padStart(2, '0')
  return `#${toHex(r)}${toHex(g)}${toHex(b)}`
}

// 颜色混合：把 color1 与 color2 按比例 weight 混合（weight 越大越接近 color2）
export function mix(color1: string, color2: string, weight: number): string {
  const [r1, g1, b1] = parseColor(color1)
  const [r2, g2, b2] = parseColor(color2)
  const r = Math.round(r1 * (1 - weight) + r2 * weight)
  const g = Math.round(g1 * (1 - weight) + g2 * weight)
  const b = Math.round(b1 * (1 - weight) + b2 * weight)
  return rgbToHex(r, g, b)
}

// 把主色写入根元素：Element Plus 的浅色/深色变体全部由主色推导
// 只需更新主色及变体，EP 内部派生变量（如 --el-menu-active-color）已自动引用主色变量
export function setPrimaryColor(primary: string, el: HTMLElement = document.documentElement): void {
  let normalized: string
  try {
    normalized = rgbToHex(...parseColor(primary))
  } catch {
    normalized = '#409eff'
  }

  el.style.setProperty('--el-color-primary', normalized)
  // EP 组件内部用 rgba(var(--el-color-primary-rgb), x) 做半透明效果，必须是 r, g, b 数字串
  const [r, g, b] = parseColor(normalized)
  el.style.setProperty('--el-color-primary-rgb', `${r}, ${g}, ${b}`)

  // 浅色变体：与主色混合白色
  const lightWeights: Record<string, number> = {
    'light-3': 0.3,
    'light-5': 0.5,
    'light-7': 0.7,
    'light-8': 0.8,
    'light-9': 0.9,
  }
  for (const [key, weight] of Object.entries(lightWeights)) {
    el.style.setProperty(`--el-color-primary-${key}`, mix(normalized, '#ffffff', weight))
  }
  // 深色变体：与主色混合黑色
  el.style.setProperty('--el-color-primary-dark-2', mix(normalized, '#000000', 0.2))

  // Sass 按需加载不会注入 dist/index.css 中 :root 的派生变量映射，
  // 浏览器回退到默认 #409eff，因此运行时切换主题色必须显式覆盖这些派生变量。
  const light9 = mix(normalized, '#ffffff', 0.9)
  const dark2 = mix(normalized, '#000000', 0.2)

  const derivedVariables: Record<string, string> = {
    // 菜单
    '--el-menu-active-color': normalized,
    '--el-menu-hover-text-color': normalized,
    '--el-menu-hover-bg-color': light9,
    // 按钮
    '--el-button-primary-bg-color': normalized,
    '--el-button-primary-border-color': normalized,
    '--el-button-primary-hover-bg-color': dark2,
    '--el-button-primary-hover-border-color': dark2,
    '--el-button-primary-active-bg-color': dark2,
    '--el-button-primary-active-border-color': dark2,
    // 输入框
    '--el-input-focus-border-color': normalized,
    '--el-input-hover-border-color': normalized,
    // 链接
    '--el-link-text-color': normalized,
    '--el-link-hover-text-color': dark2,
    // 单选按钮
    '--el-radio-button-checked-bg-color': normalized,
    '--el-radio-button-checked-border-color': normalized,
    // 标签页
    '--el-tabs-active-color': normalized,
    '--el-tabs-hover-color': normalized,
    // 开关/滑块/分页/复选框
    '--el-switch-on-color': normalized,
    '--el-slider-main-bg-color': normalized,
    '--el-pagination-hover-color': normalized,
    '--el-checkbox-checked-input-border-color': normalized,
    '--el-checkbox-checked-bg-color': normalized,
    // 选择器/树/标签
    '--el-select-input-focus-border-color': normalized,
    '--el-tree-node-hover-bg-color': light9,
    '--el-tag-bg-color': light9,
    '--el-tag-border-color': normalized,
    '--el-tag-hover-color': dark2,
  }
  for (const [variable, value] of Object.entries(derivedVariables)) {
    el.style.setProperty(variable, value)
  }

  // 持久化
  localStorage.setItem(PRIMARY_COLOR_KEY, normalized)
}

// 读取持久化的主题色；没有则返回默认 EP 蓝
export function getStoredPrimaryColor(): string {
  return localStorage.getItem(PRIMARY_COLOR_KEY) || '#409eff'
}

// 应用持久化主题色（一般在 main.ts 应用启动时调用）
export function applyStoredPrimaryColor(el: HTMLElement = document.documentElement): string {
  const color = getStoredPrimaryColor()
  setPrimaryColor(color, el)
  return color
}
