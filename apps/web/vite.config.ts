import { fileURLToPath } from 'node:url'
import vue from '@vitejs/plugin-vue'
import AutoImport from 'unplugin-auto-import/vite'
import IconsResolver from 'unplugin-icons/resolver'
import Icons from 'unplugin-icons/vite'
import UnpluginSvgComponent from 'unplugin-svg-component/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
import { defineConfig, loadEnv } from 'vite'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
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

    UnpluginSvgComponent({
      iconDir: 'src/icons', // 本地 svg 目录
      prefix: '', // 组件前缀，空=直接用文件名
      componentName: 'SvgIcon',
      treeShaking: true,
      dts: true,
      dtsDir: 'src',
    }),
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



    // 仅 pnpm run dev 时生效，pnpm run build:prod 打包时此配置被忽略。
    // rewrite: 去掉 /api 前缀，与 Nginx 的 proxy_pass 行为一致（proxy_pass 末尾有 / 会剥离 location 路径前缀）
    server: {
      proxy: {
        '/api': {
          target: env.API_PROXY_TARGET, // 复用 .env 里的地址，不重复维护 IP
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ''),
        },
      },
    },
  }
})
