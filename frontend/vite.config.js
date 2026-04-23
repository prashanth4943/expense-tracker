import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
  server: {
    proxy: {
      '/expenses': 'http://localhost:8080',
      '/categories': 'http://localhost:8080',
      '/health': 'http://localhost:8080',
    }
  }
})