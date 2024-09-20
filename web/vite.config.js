import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@src': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    proxy: {
      '/v1': {
        target: 'http://localhost:12345',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/v1/, '/v1'),
      },
      '/ws': {
        target: 'ws://localhost:12345', // WebSocket 伺服器的地址
        ws: true, // 啟用 WebSocket 代理
        changeOrigin: true,
      },
    },
  },
})
