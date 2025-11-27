<template>
  <div class="home-container">
    <h2>🎉 ¡Bienvenido a GeoChat!</h2>
    <p class="message">
      Has accedido con éxito. Este es el contenido principal de tu aplicación.
      Aquí es donde cargarás los datos del usuario y las funcionalidades de tu chat.
    </p>

    <button @click="handleLogout" class="logout-button">
      Cerrar Sesión
    </button>
    
    </div>
</template>

<script setup>
import { useRouter } from 'vue-router';

const router = useRouter();

// Función para cerrar sesión
const handleLogout = () => {
  // 1. Eliminar el token de autenticación del almacenamiento local
  localStorage.removeItem('authToken');
  
  // 2. Redirigir al usuario a la página de login
  // Al intentar acceder a la ruta '/', el Vue Router lo redirigirá automáticamente a '/login'
  // por la 'Guarda de Navegación' que creamos. Pero forzamos la navegación a '/login' para mayor claridad.
  router.push({ name: 'Login' });
};

// Opcional: Hook para cargar datos al iniciar
// import { onMounted } from 'vue';
// onMounted(() => {
//   // Aquí podrías llamar a una API para obtener los datos del perfil del usuario logueado
//   console.log('HomePage cargada. Token de autenticación verificado.');
// });
</script>

<style scoped>
.home-container {
  max-width: 800px;
  margin: 50px auto;
  padding: 30px;
  text-align: center;
  border: 1px solid #e0e0e0;
  border-radius: 10px;
  background-color: #f9f9f9;
}

h2 {
  color: #42b883; /* Color verde de Vue */
  margin-bottom: 20px;
}

.message {
  margin-bottom: 30px;
  color: #555;
  line-height: 1.6;
}

.logout-button {
  padding: 10px 20px;
  background-color: #e53935; /* Color rojo para la acción de salir */
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.logout-button:hover {
  background-color: #c62828;
}
</style>
