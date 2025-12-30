package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"Geochat/Backend/internal/services" // CORREGIDO: Usando ruta completa del módulo para servicios
)

// RequestNonce maneja la solicitud GET para obtener el mensaje a firmar.
// Ruta: GET /v1/auth/nonce?address=...
func RequestNonce(c *gin.Context) {
	address := c.Query("address") 
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Falta la dirección de Ethereum."})
		return
	}

	nonceMessage := services.GenerateNonce(address)
	
	// Devuelve el mensaje al cliente.
	c.JSON(http.StatusOK, gin.H{"success": true, "message": nonceMessage})
}

// VerifyAndAuthenticate maneja la solicitud POST para verificar la firma.
// Ruta: POST /v1/auth/verify
func VerifyAndAuthenticate(c *gin.Context) {
	var req struct {
		Address   string `json:"address"`
		Signature string `json:"signature"`
		Message   string `json:"message"` 
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Petición JSON inválida."})
		return
	}
	
	// 1. Verificar la firma
	if !services.VerifySignature(req.Address, req.Signature, req.Message) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Firma inválida o no coincide con la dirección."})
		return
	}
	
	// 2. Generar el token JWT
	token, err := services.GenerateAuthToken(req.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al generar el token de autenticación."})
		return
	}

	// 3. Devolver el token
	c.JSON(http.StatusOK, gin.H{"success": true, "token": token, "message": "Autenticación exitosa."})
}