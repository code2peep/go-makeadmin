import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import { createLocalSvgIconsPlugin } from './build/plugins/svg-icons'
// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), '')
    const apiProxyTarget = env.VITE_API_PROXY_TARGET || 'http://127.0.0.1:8000'

    return {
        // base: '/admin/',
        server: {
            host: '0.0.0.0',
            proxy: {
                '/api': {
                    target: apiProxyTarget,
                    changeOrigin: true
                }
            }
        },
        plugins: [
            vue(),
            vueJsx(),
            AutoImport({
                imports: ['vue', 'vue-router'],
                resolvers: [ElementPlusResolver({ importStyle: 'css' })],
                eslintrc: {
                    enabled: true
                }
            }),
            Components({
                directoryAsNamespace: true,
                resolvers: [ElementPlusResolver({ importStyle: 'css' })]
            }),
            createLocalSvgIconsPlugin({
                iconDirs: [fileURLToPath(new URL('./src/assets/icons', import.meta.url))],
                symbolId: 'local-icon-[dir]-[name]'
            })
        ],
        resolve: {
            alias: {
                '@': fileURLToPath(new URL('./src', import.meta.url))
            }
        },
        build: {
            // Rich text and chart chunks are intentionally split as lazy feature dependencies.
            chunkSizeWarningLimit: 900,
            rollupOptions: {
                output: {
                    manualChunks(id) {
                        if (id.includes('node_modules')) {
                            return id.toString().split('node_modules/')[1].split('/')[0].toString()
                        }
                    }
                }
            }
        }
    }
})
