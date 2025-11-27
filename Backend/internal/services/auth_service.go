package services

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go" 
	"github.com/google/uuid"      
)

// --- ALMACENAMIENTO DE NONCES ---
var nonceStore = make(map[string]string) 

// ¡CRÍTICO! Cambia esta clave en producción.
const JWT_SECRET = "tu_clave_secreta_aqui_cambiala_en_produccion" 

// GenerateNonce crea el mensaje único que el usuario firmará.
func GenerateNonce(address string) string {
	expirationTime := time.Now().Add(5 * time.Minute).Format(time.RFC3339)
	nonce := uuid.New().String()
	
	// Mensaje claro que el usuario verá en MetaMask
	message := fmt.Sprintf(
		"Bienvenido a GeoChatNativo. Firma para autenticarte. Dirección: %s. Nonce: %s. Expira: %s",
		address, 
		nonce, 
		expirationTime,
	)

	// Almacenar el nonce para validación
	nonceStore[address] = nonce
	
	return message
}

// VerifySignature realiza la verificación criptográfica.
// ESTO ES UN MOCK. Debe ser reemplazado por la lógica de ecrecover real.
func VerifySignature(address, signature, message string) bool {
	// Aquí va la lógica real de ecrecover de Ethereum.
	if address != "" && signature != "" && message != "" {
		fmt.Printf("Verificación SIWE MOCK exitosa para la dirección: %s\n", address)
		return true 
	}
	
	return false
}

// GenerateAuthToken crea un token JWT para la sesión.
func GenerateAuthToken(address string) (string, error) {
	// Definir las claims (información) del token
	claims := jwt.MapClaims{
		"authorized": true,
		"address":    address,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // 24 horas de expiración
	}

	// Crear y firmar el token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
