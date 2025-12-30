<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 p-4">
    <div class="w-full max-w-md bg-white p-8 rounded-xl shadow-2xl border border-gray-100">
      
      <!-- Encabezado -->
      <h2 class="text-3xl font-extrabold text-gray-900 text-center mb-2">GeoChat</h2>
      <p class="text-center text-sm text-gray-500 mb-8">Autenticación Web3 (SIWE)</p>

      <!-- Estado de la Cartera -->
      <div class="mb-6 p-4 border rounded-lg" :class="[isWalletConnected ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200']">
        <div class="flex items-center space-x-3">
          <svg v-if="isWalletConnected" class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944c-3.141 0-6.104 1.135-8.318 3.04a12.022 12.022 0 00-.655 1.543l.63.63a10.02 10.02 0 0113.84-1.259zM19 12a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
          <svg v-else class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
          <span class="font-medium" :class="[isWalletConnected ? 'text-green-800' : 'text-red-800']">
            {{ isWalletConnected ? 'Conectado' : 'Desconectado' }}
          </span>
        </div>
        <p v-if="address" class="mt-2 text-xs text-gray-600 truncate">
          Dirección: <span class="font-mono">{{ address }}</span>
        </p>
      </div>

      <!-- Botones de Acción -->
      <div class="space-y-4">
        <!-- 1. Conectar Cartera -->
        <button 
          @click="connectWallet" 
          :disabled="isWalletConnected || isSigning"
          class="w-full flex justify-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 disabled:bg-indigo-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition duration-150 ease-in-out">
          Conectar Cartera Web3
        </button>

        <!-- 2. Solicitar Nonce y Firmar -->
        <button 
          @click="handleSignAndVerify" 
          :disabled="!isWalletConnected || isSigning"
          class="w-full flex justify-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-indigo-700 bg-indigo-100 hover:bg-indigo-200 disabled:bg-gray-100 disabled:text-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition duration-150 ease-in-out">
          <span v-if="isSigning" class="flex items-center">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-indigo-700" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
            Firmando...
          </span>
          <span v-else>
            Firmar Mensaje y Entrar
          </span>
        </button>
      </div>

      <!-- Área de Mensajes -->
      <div v-if="statusMessage" :class="{'bg-yellow-100 border-yellow-400': statusType === 'info', 'bg-red-100 border-red-400': statusType === 'error', 'bg-green-100 border-green-400': statusType === 'success'}" class="mt-6 p-3 border rounded-lg text-sm text-gray-700">
        <p class="font-medium" :class="{'text-yellow-800': statusType === 'info', 'text-red-800': statusType === 'error', 'text-green-800': statusType === 'success'}">
          {{ statusMessage }}
        </p>
      </div>
      
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
// Importamos las funciones de servicio que tienen la URL correcta del backend
import { getNonce, authenticateWeb3 } from '@/api/authService.js'; 

// Usamos el router de Vue para la navegación
const router = useRouter();
const window = globalThis; // Acceso al proveedor Web3

// --- ESTADO ---
const address = ref(''); 
const nonceMessage = ref('');
const signature = ref('');
const isSigning = ref(false);
const statusMessage = ref('');
const statusType = ref('info'); // 'info', 'success', 'error'

// --- COMPUTADAS ---
const isWalletConnected = computed(() => address.value.length > 0);

/**
 * Conecta la cartera Web3 usando window.ethereum
 */
const connectWallet = async () => {
    statusMessage.value = 'Conectando a la cartera...';
    statusType.value = 'info';

    if (!window.ethereum) {
        statusMessage.value = 'MetaMask o proveedor Web3 no detectado.';
        statusType.value = 'error';
        return;
    }

    try {
        // Solicita acceso a las cuentas
        const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
        
        if (accounts.length > 0) {
            address.value = accounts[0]; 
            statusMessage.value = '¡Conexión exitosa! Ahora firma para autenticarte.';
            statusType.value = 'success';
        } else {
            statusMessage.value = 'Conexión rechazada por el usuario.';
            statusType.value = 'error';
        }
    } catch (error) {
        statusMessage.value = `Error al conectar: ${error.message || 'Error desconocido'}`;
        statusType.value = 'error';
        console.error("Error connecting wallet:", error);
    }
}

/**
 * Firma el mensaje nonce usando personal_sign.
 * @param {string} message - El nonce a firmar.
 * @returns {Promise<string>} La firma.
 */
const signMessage = async (message) => {
    if (!window.ethereum) {
        throw new Error("Proveedor Web3 no disponible.");
    }
    
    statusMessage.value = 'Esperando la firma en tu cartera...';
    
    // El método personal_sign es el estándar para la autenticación SIWE
    const sig = await window.ethereum.request({
        method: 'personal_sign',
        params: [message, address.value], // Orden: mensaje, cuenta
    });
    
    return sig;
}

// --- Flujo Principal SIWE ---
const handleSignAndVerify = async () => {
    if (!isWalletConnected.value) {
        statusMessage.value = 'Debes conectar tu cartera primero.';
        statusType.value = 'error';
        return;
    }

    isSigning.value = true;
    statusType.value = 'info';

    try {
        // 1. Obtener el Nonce (Usando la función del servicio con la URL correcta)
        const nonce = await getNonce(address.value);
        nonceMessage.value = nonce;
        
        // 2. Firmar el Mensaje
        const sig = await signMessage(nonce);
        signature.value = sig;
        
        // 3. Enviar Firma para Verificación y obtener JWT (Usando la función del servicio)
        const authResult = await authenticateWeb3(address.value, signature.value, nonceMessage.value);

        if (authResult.success) {
            // Guardar el token
            localStorage.setItem('authToken', authResult.token);
            statusMessage.value = '¡Autenticación exitosa! Redirigiendo...';
            statusType.value = 'success';
            
            // Redirigir a la página principal (Home)
            setTimeout(() => {
                router.push({ name: 'Home' }); 
            }, 1000); // Pequeño retraso para que el usuario vea el éxito

        } else {
             // Error manejado por authenticateWeb3
            statusMessage.value = authResult.message || 'Error de autenticación: Firma inválida o token no recibido.';
            statusType.value = 'error';
        }

    } catch (error) {
        // Este catch maneja errores de getNonce o signMessage (ej. usuario cancela)
        statusMessage.value = `Fallo en el proceso de autenticación: ${error.message || 'Firma de mensaje cancelada o fallida.'}`;
        statusType.value = 'error';
        console.error('Authentication flow error:', error);
    } finally {
        isSigning.value = false;
    }
}
</script>