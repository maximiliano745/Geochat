import { createApp } from 'vue';
// Importa el componente raíz (App.vue)
import App from './App.vue'; 
// Importa el sistema de rutas configurado (asumiendo src/router/index.js)
import router from './router/index.js'; 

// NOTA: Si necesitas estilos globales que no sean Tailwind, puedes descomentar la siguiente línea:
// import './style.css'; 

// --- Inicialización de la Aplicación ---
createApp(App)
    // Conecta el sistema de rutas
    .use(router)
    // Monta la aplicación en el elemento con id="app"
    .mount('#app');

console.log("GeoChat Vue App: Inicialización completa.");