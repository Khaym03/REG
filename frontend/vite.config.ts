import { defineConfig } from 'vite'
import react, { reactCompilerPreset } from '@vitejs/plugin-react'
import babel from '@rolldown/plugin-babel'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'
import { tanstackRouter } from '@tanstack/router-plugin/vite'
import wails from '@wailsio/runtime/plugins/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    tanstackRouter(),
    react(),
    babel({ presets: [reactCompilerPreset()] }),
    tailwindcss(),
    wails('./bindings')
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      bindings: path.resolve(__dirname, './bindings'),
      wails: path.resolve(__dirname, './wailsjs')
    }
  },
  server: {
    host: '127.0.0.1',
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true
  }
})
