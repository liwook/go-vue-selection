import vue from '@vitejs/plugin-vue'
import { fileURLToPath } from 'url'
import { defineConfig } from 'vite'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import Icons from 'unplugin-icons/vite'
import IconsResolver from 'unplugin-icons/resolver'
// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [
        ElementPlusResolver({ importStyle: 'sass' }),
        IconsResolver({ prefix: 'Icon', enabledCollections: ['ep', 'mdi'] }),
      ],
    }),
    Components({
      resolvers: [
        ElementPlusResolver({ importStyle: 'sass' }),
        IconsResolver({ prefix: 'Icon', enabledCollections: ['ep', 'mdi'] }),
      ],
    }),
    Icons(),
  ],

  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `
          @use "@/styles/element/index.scss" as *;
          @use "@/styles/variables.scss" as *;
        `,
      },
    },
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
})
