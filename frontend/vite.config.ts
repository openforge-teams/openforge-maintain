import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api/v2/core': {
        target: 'http://localhost:9999',
        changeOrigin: true,
      },
      '/api/v2': {
        target: 'http://localhost:10000',
        changeOrigin: true,
      },
      '/ws': {
        target: 'ws://localhost:10000',
        ws: true,
        changeOrigin: true,
      },
    },
  },
})
