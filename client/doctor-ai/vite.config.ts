// vite.config.ts
import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig(({ mode }) => {
  // Load env file based on `mode` in the current working directory.
  const env = loadEnv(mode, process.cwd(), '')
  
  // List of env variables to expose
  const envPrefix = ['VITE_']
  
  return {
    plugins: [react(), tailwindcss()],
    envPrefix,
    server: {
      proxy: {
        '/api': {
          target: process.env.VITE_BACKEND_URL || 'http://localhost:8080',
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, '')
        }
      }
    },
    define: {
      'process.env.VITE_GROQAPIKEY': JSON.stringify(env.VITE_GROQAPIKEY)
    }
  }
})