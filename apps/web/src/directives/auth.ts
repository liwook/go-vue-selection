import type { Directive, DirectiveBinding } from 'vue'
import { useUserStore } from '@/store/modules/user'

// 自定义指令 v-auth：按钮级权限控制
// 用法：<el-button v-auth="'btn.Trademark.add'">添加品牌</el-button>
// 原理：元素挂载后从 pinia 读取 buttons，
//       若不含指令值则直接 el.remove() 移除该 DOM（不渲染 = 无权限）
const auth: Directive<HTMLElement, string> = {
  mounted(el: HTMLElement, binding: DirectiveBinding<string>) {
    const userStore = useUserStore()
    // 保守策略：buttons 为空（未登录 / 信息未加载）时按「无权限」处理
    const buttons: string[] = userStore.buttons ?? []
    const required = binding.value
    if (!buttons.includes(required)) {
      el.remove()
    }
  },
}

// 说明：项目路由守卫（router/index.ts 的 beforeEach）会先 await getUserInfo()
// 把 buttons 写入 store，再 next() 放行渲染；因此组件 mounted 时 buttons 必然已就绪，
// 不会出现「信息未加载导致按钮被误删」的情况。上面的 ?? [] 仅作兜底防御。

export default auth
