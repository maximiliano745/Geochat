import { createApp } from 'vue';
import App from './App.vue'; 
import router from './router/index.js'; 

// 1. Importamos la lógica del laboratorio (main.ts)
// Al importarlo así, el código dentro de main.ts se ejecuta automáticamente
import './main.ts'; 

// 2. Creamos y montamos la aplicación UNA sola vez
const app = createApp(App);

app.use(router);
app.mount('#app');

console.log("🚀 GeoChat Vue App: Inicialización completa y Laboratorio activo.");