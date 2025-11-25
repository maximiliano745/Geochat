// src/router/index.js

import { createRouter, createWebHistory } from 'vue-router';

// 1. Importa tus componentes (asumiendo que están en src/views o src/pages)
import LoginPage from '../pages/LoginPage.vue'; 
import HomePage from '../pages/HomePage.vue';
import NotFoundPage from '../pages/NotFoundPage.vue';

// 2. Define las rutas
const routes = [
  // Ruta de Login (Pública)
  {
    path: '/login',
    name: 'Login',
    component: LoginPage,
    meta: { requiresAuth: false } // No requiere autenticación
  },
  
  // Ruta Principal (Protegida)
  {
    path: '/',
    name: 'Home',
    component: HomePage,
    meta: { requiresAuth: true } // REQUIERE autenticación
  },

  // Ruta 404 (Catch-all)
  {
    path: '/:catchAll(.*)',
    name: 'NotFound',
    component: NotFoundPage,
  }
];

// 3. Crea la instancia del router
const router = createRouter({
  history: createWebHistory(), // Usa el historial de navegación HTML5 (URLs limpias)
  routes,
});

// 4. Implementa el "Guarda de Navegación" (Navigation Guard)
router.beforeEach((to, from, next) => {
  // Obtiene el estado de autenticación (ej: token guardado)
  const isAuthenticated = localStorage.getItem('authToken'); 
  
  // Si la ruta requiere autenticación Y el usuario NO está autenticado
  if (to.meta.requiresAuth && !isAuthenticated) {
    // Redirige al usuario a la página de login
    next({ name: 'Login' });
  } else {
    // Permite el acceso a la ruta
    next();
  }
});

export default router;