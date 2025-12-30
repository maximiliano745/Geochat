package location

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket" // Importación de Gorila/WS
	"github.com/jackc/pgx/v5"
)

// NOTA: Las estructuras UserLocation, LocationCache y Broadcaster 
// y las funciones AddClient, RemoveClient, UpdateCache, GetAllLocations 
// y InitCache DEBEN estar definidas en otro archivo (ej. cache.go) 
// dentro del mismo paquete 'location'.

// Upgrader para WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Permite cualquier origen para facilitar el desarrollo/pruebas.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// LocationPostHandler maneja las solicitudes POST para actualizar la ubicación de un usuario.
// Debe ser inyectado con la conexión a la base de datos.
// Ruta: POST /v1/location
func LocationPostHandler(db *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loc UserLocation
		
		// 1. Validar y bindear el JSON entrante
		if err := c.ShouldBindJSON(&loc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Petición JSON inválida", "error": err.Error()})
			return
		}

		// 🚨 TODO: Aquí se debe implementar la verificación del token JWT
		// Por ahora, asumimos que el userId en el cuerpo es válido.
		if loc.UserID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "El ID de usuario es obligatorio."})
			return
		}

		// 2. Insertar/Actualizar la ubicación en la base de datos (PostgreSQL)
		// Usamos UPSERT para que inserte si no existe o actualice si ya existe.
		query := `
			INSERT INTO user_locations (user_id, latitude, longitude)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id) DO UPDATE
			SET latitude = $2, longitude = $3, created_at = CURRENT_TIMESTAMP
		`
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		_, err := db.Exec(ctx, query, loc.UserID, loc.Latitude, loc.Longitude)
		if err != nil {
			log.Printf("ERROR DB: Fallo al actualizar ubicación para %s: %v", loc.UserID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error al guardar en la DB."})
			return
		}

		// 3. Actualizar la caché en memoria y notificar al Broadcaster
		UpdateCache(loc) 
		
		log.Printf("INFO: Ubicación actualizada para el usuario: %s (Lat: %.4f, Lon: %.4f)", loc.UserID, loc.Latitude, loc.Longitude)
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Ubicación reportada exitosamente."})
	}
}

// LocationGetHandler devuelve todas las ubicaciones almacenadas en la caché.
// Ruta: GET /v1/location
func LocationGetHandler(db *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 🚨 TODO: Implementar la verificación de autenticación
		
		// Obtener todas las ubicaciones desde la caché
		locations := GetAllLocations() 
		
		c.JSON(http.StatusOK, gin.H{"success": true, "locations": locations, "count": len(locations)})
	}
}

// LocationWsHandler maneja la conexión WebSocket para actualizaciones en tiempo real.
// Ruta: GET /v1/ws/location
func LocationWsHandler(db *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 🚨 TODO: Implementar verificación de autenticación antes de la actualización
		
		// 1. Actualizar la conexión HTTP a WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("ERROR WS: Falló la actualización de WebSocket: %v", err)
			return
		}
		defer conn.Close()

		// 2. Crear un canal para este cliente WS
		clientChan := make(chan UserLocation)
		broadcaster.AddClient(clientChan) // Registrar el canal en el Broadcaster
		defer broadcaster.RemoveClient(clientChan) // Asegurar la limpieza al salir

		// 3. Primer Envío: Enviar el estado actual completo
		initialLocations := GetAllLocations()
		if len(initialLocations) > 0 {
			initialData, err := json.Marshal(gin.H{"type": "init", "locations": initialLocations})
			if err != nil {
				log.Printf("ERROR WS: Fallo al serializar datos iniciales: %v", err)
			} else {
				if err := conn.WriteMessage(websocket.TextMessage, initialData); err != nil {
					log.Printf("ERROR WS: Fallo al enviar datos iniciales: %v", err)
					return // Si el envío falla, terminamos la conexión.
				}
			}
		}

		// Loop de escucha: lee el canal clientChan y envía datos al cliente WS
		for loc := range clientChan {
			// Serializar el mensaje de actualización.
			// Incluye el tipo 'update' para que el frontend sepa cómo manejarlo.
			updateData, err := json.Marshal(gin.H{"type": "update", "location": loc})
			if err != nil {
				log.Printf("ERROR WS: Fallo al serializar actualización: %v", err)
				continue
			}

			// Enviar el mensaje al cliente
			if err := conn.WriteMessage(websocket.TextMessage, updateData); err != nil {
				log.Printf("ERROR WS: Fallo al escribir mensaje: %v", err)
				// Si el envío falla, es probable que el cliente esté desconectado. 
				// Rompemos el loop, y el defer RemoveClient limpiará el Broadcaster.
				return 
			}
		}
	}
}