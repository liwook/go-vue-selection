import { createPinia } from 'pinia'
import { createApp } from 'vue'
import SvgIcon from '~virtual/svg-component'
import App from './App.vue'
import router from './router'

// 样式：全局重置
import './style.css'
// Element Plus 暗黑模式变量（必须放在自定义 style.css 之前，否则会被业务样式覆盖）
import 'element-plus/theme-chalk/dark/css-vars.css'

// 应用启动时先恢复用户持久化的主题色，避免刷新后闪回默认蓝
import { applyStoredPrimaryColor } from '@/utils/theme'

applyStoredPrimaryColor()

const app = createApp(App)

// 全局注册图标组件（任意组件内可直接用 <SvgIcon name="xxx" />）
app.component('SvgIcon', SvgIcon)

// 注册 pinia（必须在 router 之前，守卫里 useUserStore() 才能拿到激活的实例）
app.use(createPinia())

// 注册路由
app.use(router)

app.mount('#app')
