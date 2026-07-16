import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    // 本地开发保持浏览器同源请求，由 Vite 转发到 Gin，避免硬编码后端地址。
    proxy: {
      '/api': 'http://127.0.0.1:8080',
    },
  },
})
