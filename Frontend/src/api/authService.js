// --- ¡CORRECCIÓN DE PUERTO A 8080! ---
// ⚠️ DEBES pegar tu URL de Codespaces REAL aquí. 
// Ejemplo ajustado a tu puerto 8080: https://opulent-chainsaw-xpprp6gww7h6jg6-8080.app.github.dev
const API_ROOT = "https://opulent-chainsaw-xpprp6gww7h6jg6-8080.app.github.dev/"; 

const BASE_API_URL = `${API_ROOT}/v1/auth`; 
const LOCATION_API_URL = `${API_ROOT}/v1/location`; 

/**
 * 1. Solicita al servidor un nonce (mensaje único) para la dirección dada.
 * @param {string} address - Dirección de Ethereum del usuario.
 * @returns {Promise<string>} El mensaje (nonce) a ser firmado.
 */
export async function getNonce(address) {
    try {
        const response = await fetch(`${BASE_API_URL}/nonce?address=${address}`);
        const data = await response.json();

        if (!response.ok || !data.success) {
            throw new Error(data.message || "Fallo al obtener el nonce del servidor.");
        }
        
        // El backend Go devuelve el mensaje en el campo 'message'
        return data.message; 
    } catch (error) {
        console.error("Error en getNonce:", error);
        throw new Error("Error de red al solicitar el nonce.");
    }
}

/**
 * 2. Envía la firma y la información al servidor para verificar y obtener el JWT.
 * @param {string} address - Dirección de Ethereum.
 * @param {string} signature - Firma generada por el usuario.
 * @param {string} message - El mensaje (nonce) firmado.
 * @returns {Promise<object>} El token JWT o un mensaje de error.
 */
export async function authenticateWeb3(address, signature, message) {
    try {
        const response = await fetch(`${BASE_API_URL}/verify`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                address, 
                signature, 
                message 
            }),
        });

        const data = await response.json();
        
        if (!response.ok || !data.success) {
             // El backend Go ya maneja los errores 401 (no autorizado)
            throw new Error(data.message || "Fallo en la verificación de la firma.");
        }

        return { success: true, token: data.token };

    } catch (error) {
        console.error("Error en authenticateWeb3:", error);
        return { success: false, message: error.message || "Error de red o firma inválida." };
    }
}

/**
 * 3. Reporta la ubicación actual del usuario al servidor.
 * @param {number} latitude - Latitud.
 * @param {number} longitude - Longitud.
 * @returns {Promise<object>} Respuesta del servidor.
 */
export async function reportLocation(latitude, longitude) {
    const token = localStorage.getItem('authToken');
    if (!token) {
        throw new Error("No hay token de autenticación. Redirigiendo a Login.");
    }

    try {
        // En una aplicación real, aquí también se extraería el address del JWT para el user_id, 
        // pero por ahora, solo enviamos los datos y confiamos en el middleware para el user_id.
        const response = await fetch(LOCATION_API_URL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                // ¡IMPORTANTE! Envía el token JWT en el header Authorization
                'Authorization': `Bearer ${token}`, 
            },
            body: JSON.stringify({ 
                latitude, 
                longitude 
            }),
        });

        const data = await response.json();
        
        if (!response.ok || !data.success) {
            throw new Error(data.message || `Fallo al reportar ubicación: ${response.status}`);
        }

        return { success: true, message: data.message };

    } catch (error) {
        console.error("Error en reportLocation:", error);
        throw new Error(error.message || "Error de red al reportar ubicación.");
    }
}