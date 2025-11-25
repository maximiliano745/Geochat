<template>
  <div class="login-container">
    <h2>Iniciar Sesión en GeoChat</h2>
    
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="email">Email:</label>
        <input 
          type="email" 
          id="email" 
          v-model="email" 
          required
        />
      </div>
      
      <div class="form-group">
        <label for="password">Contraseña:</label>
        <input 
          type="password" 
          id="password" 
          v-model="password" 
          required
        />
      </div>
      
      <button type="submit" :disabled="isLoading">
        {{ isLoading ? 'Cargando...' : 'Ingresar' }}
      </button>
      
      <p v-if="error" class="error-message">{{ error }}</p>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
// Importa la función de tu servicio de API (¡debes crear este archivo!)
import { loginUser } from '../api/authService'; 

// 1. Estado Local para los Inputs y la Interfaz
const email = ref('');
const password = ref('');
const error = ref(null);
const isLoading = ref(false);

const router = useRouter();

// 2. Función que maneja el envío del formulario
const handleLogin = async () => {
  // Limpiar estados anteriores
  error.value = null;
  isLoading.value = true;
  
  try {
    // 3. Llamar a la API
    const response = await loginUser(email.value, password.value);
    
    if (response.success && response.token) {
      // 4. Éxito: Guardar el token de sesión y redirigir
      localStorage.setItem('authToken', response.token);
      
      // Limpia los inputs al terminar
      email.value = '';
      password.value = '';
      
      // Redirige a la página principal (que es una ruta protegida)
      router.push({ name: 'Home' }); 
      
    } else {
      // Manejo de errores de credenciales inválidas (ej: 401 del backend)
      error.value = response.message || 'Credenciales inválidas. Intenta de nuevo.';
    }
    
  } catch (err) {
    // Manejo de errores de red o servidor
    error.value = 'Error al conectar con el servidor. Revisa tu conexión.';
    console.error(err);
    
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  max-width: 400px;
  margin: 50px auto;
  padding: 20px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  border-radius: 8px;
}
.form-group {
  margin-bottom: 15px;
}
label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}
input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
}
button {
  width: 100%;
  padding: 10px;
  background-color: #42b883; /* Color verde de Vue */
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}
button:disabled {
  background-color: #888;
  cursor: not-allowed;
}
.error-message {
  color: #e53935;
  margin-top: 10px;
  text-align: center;
}
</style>
