import vue from '@vitejs/plugin-vue'
import { fileURLToPath } from 'url'
import { defineConfig } from 'vite'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [
        ElementPlusResolver({ importStyle: 'sass' }),
      ],
    }),
    Components({
      resolvers: [
        ElementPlusResolver({ importStyle: 'sass' }),
      ],
    }),
    // 图标按需导入（unplugin-icons）与 UnpluginSvgComponent 见后续步骤
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
