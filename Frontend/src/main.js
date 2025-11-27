import { createApp } from 'vue'
import './style.css' // Asegúrate de que este archivo exista
import App from './App.vue'

// Importa y configura el router
// Asumimos que has creado este archivo en src/router/index.js
import router from './router' 

// --- Inicialización de Vue y Router ---
// Monta la aplicación e inyecta el router.
createApp(App).use(router).mount('#app')

// Si tienes inicialización de Firebase/Web3, debe ir aquí o en App.vue,
// pero es vital que no bloquee el montaje inicial.

// --- EJEMPLO SI ESTUVIERAS USANDO FIREBASE (NO USAR SI NO LO NECESITAS) ---
/*
import { initializeApp } from 'firebase/app';
const firebaseConfig = JSON.parse(__firebase_config);
const app = initializeApp(firebaseConfig);
createApp(App, { firebaseApp: app }).mount('#app');
*/