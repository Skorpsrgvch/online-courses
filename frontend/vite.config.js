import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: true, // Разрешаем доступ извне контейнера (нужно для Docker)
    port: 3000, // Порт для разработки
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // Адрес вашего Go бэкенда
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
