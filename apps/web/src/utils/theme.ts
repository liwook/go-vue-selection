// 主题色工具：把选中的主色写入根元素的 CSS 变量，驱动 Element Plus 整体换肤
// 用法：setPrimaryColor('#1677ff')

// 把 hex 颜色转成 rgb 元组
export function hexToRgb(hex: string): [number, number, number] {
  // 支持 #fff 简写与 #ffffff 完整形式
  let normalized = hex.replace('#', '')
  if (normalized.length === 3) {
    normalized = normalized
      .split('')
      .map((c) => c + c)
      .join('')
  }
  const num = parseInt(normalized, 16)
  return [(num >> 16) & 255, (num >> 8) & 255, num & 255]
}

// 把 rgb 三元组转回 hex 字符串
export function rgbToHex(r: number, g: number, b: number): string {
  const toHex = (n: number) => n.toString(16).padStart(2, '0')
  return `#${toHex(r)}${toHex(g)}${toHex(b)}`
}

// 颜色混合：把 color1 与 color2 按比例 weight 混合（weight 越大越接近 color2）
export function mix(color1: string, color2: string, weight: number): string {
  const [r1, g1, b1] = hexToRgb(color1)
  const [r2, g2, b2] = hexToRgb(color2)
  const r = Math.round(r1 * (1 - weight) + r2 * weight)
  const g = Math.round(g1 * (1 - weight) + g2 * weight)
  const b = Math.round(b1 * (1 - weight) + b2 * weight)
  return rgbToHex(r, g, b)
}

// 把主色写入根元素：Element Plus 的浅色/深色变体全部由主色推导
export function setPrimaryColor(primary: string, el: HTMLElement = document.documentElement): void {
  el.style.setProperty('--el-color-primary', primary)

  // 浅色变体：与主色混合白色
  const lightWeights: Record<string, number> = {
    'light-3': 0.3,
    'light-5': 0.5,
    'light-7': 0.7,
    'light-8': 0.8,
    'light-9': 0.9,
  }
  for (const [key, weight] of Object.entries(lightWeights)) {
    el.style.setProperty(`--el-color-primary-${key}`, mix(primary, '#ffffff', weight))
  }
  // 深色变体：与主色混合黑色
  el.style.setProperty('--el-color-primary-dark-2', mix(primary, '#000000', 0.2))
}
