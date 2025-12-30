import { createRouter, createWebHistory } from 'vue-router';
import LoginPage from '../pages/LoginPage.vue'; // Importa el componente de Login
import HomePage from '../pages/HomePage.vue'; // Asegúrate de que este componente exista

const routes = [
  {
    path: '/',
    name: 'Login',
    component: LoginPage,
  },
  {
    path: '/home',
    name: 'Home',
    component: HomePage, // Usamos HomePage para el componente principal
    meta: { requiresAuth: true } // REQUIERE AUTENTICACIÓN
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Guardia de navegación global para verificar la autenticación
router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const isAuthenticated = localStorage.getItem('authToken');

  if (requiresAuth && !isAuthenticated) {
    // Si se requiere auth y no está autenticado, redirigir a login
    next({ name: 'Login' });
  } else if (!requiresAuth && isAuthenticated) {
    // Si está en login/registro y ya está autenticado, redirigir a home
    next({ name: 'Home' });
  } else {
    // Continuar normalmente
    next();
  }
});

export default router;