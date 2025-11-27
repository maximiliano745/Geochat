<template>
  <div class="login-container">
    <h2>Conectar a GeoChat (Web3)</h2>
    
    <button 
      @click="connectWallet" 
      :disabled="isLoading" 
      class="connect-button"
    >
      {{ isLoading ? 'Esperando firma...' : 'Conectar Cartera (MetaMask)' }}
    </button>
    
    <p v-if="error" class="error-message">{{ error }}</p>
    
    <p v-if="currentAccount" class="success-message">
      Conectado como: {{ currentAccount.substring(0, 6) }}...{{ currentAccount.substring(currentAccount.length - 4) }}
    </p>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { ethers } from 'ethers';
import { authenticateWeb3, getNonce } from '../api/authService'; 

const error = ref(null);
const isLoading = ref(false);
const currentAccount = ref(''); 
const router = useRouter();

// 1. Función principal para conectar y autenticar
const connectWallet = async () => {
  error.value = null;
  isLoading.value = true;
  
  try {
    // 2. Comprobar si MetaMask está disponible
    if (!window.ethereum) {
      error.value = 'MetaMask no está instalado. Por favor, instálalo para continuar.';
      isLoading.value = false;
      return;
    }
    
    // 3. Conectar la cartera y obtener la dirección
    const provider = new ethers.BrowserProvider(window.ethereum);
    const accounts = await provider.send('eth_requestAccounts', []);
    const address = accounts[0]; 
    currentAccount.value = address;
    
    // 4. Obtener un mensaje único (nonce) del servidor
    const message = await getNonce(address);
    
    // 5. Firma del Mensaje
    const signer = await provider.getSigner(address);
    const signature = await signer.signMessage(message); 

    // 6. Verificar la firma con el Backend
    const authResponse = await authenticateWeb3(address, signature, message);

    if (authResponse.success && authResponse.token) {
      // 7. Éxito: Guardar el token de sesión y redirigir
      localStorage.setItem('authToken', authResponse.token);
      router.push({ name: 'Home' }); 
    } else {
      error.value = authResponse.message;
    }
    
  } catch (err) {
    // Errores de usuario (cancelación de firma) o de red/servidor
    error.value = err.message || 'Ocurrió un error al conectar o firmar la transacción.';
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped>
/* Estilos similares a la versión Web2, pero enfocados en el botón de conexión */
.connect-button {
  /* Estilo diferente para destacar la conexión Web3 */
  background-color: #f6851b; /* Color de MetaMask/Ethereum */
  padding: 12px 25px;
  font-size: 1.1em;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.3s;
  margin-top: 20px;
}
.connect-button:hover {
  background-color: #e5740a;
}
.success-message {
    color: #42b883;
    font-weight: bold;
    margin-top: 15px;
}
</style>
