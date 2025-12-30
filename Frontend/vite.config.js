import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path' // Importamos el módulo 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  // 🛑 FIX CLAVE: Configurar la base de la URL pública a la raíz.
  base: '/', 
  
  // 🟢 FIX: Configuración de alias para resolver "@/" a "src/"
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },

  // --- Configuración de Servidor para Entornos Virtualizados ---
  server: {
    // Escuchar en todas las interfaces para ser accesible en Codespaces
    host: '0.0.0.0', 
    // Usar el puerto 5174 como MANDATORIO para asegurar la exposición del puerto.
    strictPort: true, 
    port: 5175 
  }
})