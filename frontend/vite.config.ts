import { defineConfig } from 'vite'
import Vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [Vue()],
  resolve: {
    alias: {
      '@': 'C:/laragon/www/dail-it',
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:6366',
        changeOrigin: true,
      },
    },
  },
})
