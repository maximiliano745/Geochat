// src/main.js

import { createApp } from 'vue';
import App from './App.vue';
import router from './router'; // Importa el router que acabas de crear

// Crea la aplicación e inyecta el router
createApp(App).use(router).mount('#app');
