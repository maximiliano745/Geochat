import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0', 
    port: 5174, // <-- CAMBIADO A 5174
    allowedHosts: [
      'opulent-chainsaw-xpprp6gww7h6jg6-5174.app.github.dev', // <-- Nueva URL para 5174 (Aunque Codespaces la añade automáticamente)
      '135.237.130.231'            
    ]
  }
})