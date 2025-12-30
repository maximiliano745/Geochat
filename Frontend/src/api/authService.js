// URL base de tu backend (Asegúrate de que tu backend esté corriendo en esta dirección)
// Nota: Si tu backend está en un dominio diferente, necesitarás manejar CORS.
const API_BASE_URL = 'http://localhost:8081'; // CAMBIADO a 8081 para coincidir con el backend Go

/**
 * Función para obtener un Nonce (número de un solo uso) para una dirección Ethereum.
 * Esto inicia el proceso de Sign-In with Ethereum (SIWE).
 * @param {string} address - La dirección de la cartera del usuario.
 * @returns {Promise<string>} El nonce devuelto por el servidor.
 */
export async function getNonce(address) {
    if (!address) {
        throw new Error("Se requiere una dirección para obtener el nonce.");
    }
    
    console.log(`Solicitando nonce para la dirección: ${address}`);
    
    try {
        const response = await fetch(`${API_BASE_URL}/api/auth/nonce?address=${address}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || `Error ${response.status} al solicitar el nonce.`);
        }

        const data = await response.json();
        // El backend debe devolver el nonce en el campo 'nonce'
        return data.nonce; 

    } catch (error) {
        console.error("Fallo al obtener el nonce:", error);
        throw new Error(`Fallo en la comunicación con el servidor: ${error.message}`);
    }
}


/**
 * Función para autenticar al usuario enviando la firma y el mensaje Nonce.
 * @param {string} address - La dirección de la cartera.
 * @param {string} signature - La firma generada por el usuario.
 * @param {string} message - El mensaje Nonce original.
 * @returns {Promise<{success: boolean, token?: string, message?: string}>} Resultado de la autenticación.
 */
export async function authenticateWeb3(address, signature, message) {
    if (!address || !signature || !message) {
        return { success: false, message: "Faltan datos de autenticación." };
    }

    console.log("Enviando firma para verificación...");

    try {
        const response = await fetch(`${API_BASE_URL}/api/auth/verify`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ address, signature, message })
        });

        const data = await response.json();

        if (response.ok) {
            // Si la verificación es exitosa, el backend debería devolver un JWT
            return {
                success: true,
                token: data.token,
                message: "Verificación exitosa."
            };
        } else {
            // Si la verificación falla, el backend debería devolver un mensaje de error
            return {
                success: false,
                message: data.message || "Fallo en la verificación de la firma."
            };
        }

    } catch (error) {
        console.error("Error durante la autenticación:", error);
        return { success: false, message: `Error de red: ${error.message}` };
    }
}