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
        preserveColor: /(welcome|logo)/, // 保留彩色 svg 的原色（不替换为 currentColor）
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

    build: {
      // 拆 vendor chunk：把体积大且基本不变的 echarts 单独成包，
      // 浏览器可并行下载 + 长期缓存命中后秒开，避免全量塞进 screen chunk 拖慢首屏。
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (id.includes('node_modules/echarts') || id.includes('node_modules/zrender')) {
              return 'echarts'
            }
            if (id.includes('node_modules')) {
              return 'vendor'
            }
          },
        },
      },
    },

    // 仅 pnpm run dev 时生效，pnpm run build:prod 打包时此配置被忽略。
    // 后端接口前缀为 /api/v1/，前端 baseURL 用 /api + 请求路径 /v1/... 拼成 /api/v1/...，
    // 因此代理【不】重写（剥离）/api 前缀，直接原样转发，避免把 /api 误删导致 404。
    // 生产环境 Nginx 需用 proxy_pass http://backend（末尾不带 /）以保持同样行为。
    server: {
      proxy: {
        '/api': {
          target: env.API_PROXY_TARGET, // 复用 .env 里的地址，不重复维护 IP
          changeOrigin: true,
          // 不 rewrite：保留 /api 前缀，转发到后端即 /api/v1/...
        },
      },
    },
  }
})
