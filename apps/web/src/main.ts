import { createApp } from 'vue'
import './style.css'
import SvgIcon from '~virtual/svg-component' // ← 插件生成的虚拟模块
import App from './App.vue'

const app = createApp(App)
app.component('SvgIcon', SvgIcon) // ← 注册为全局组件
app.mount('#app')