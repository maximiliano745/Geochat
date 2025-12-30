package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors" 
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5" // Para la conexión a PostgreSQL
	"Geochat/Backend/internal/location" // Asegúrate de usar esta ruta completa
    "Geochat/Backend/internal/handlers" // Asegúrate de usar esta ruta completa
)

// initDBConnection simula la conexión a PostgreSQL.
// 🚨 NOTA IMPORTANTE: Esta es una función MOCK. 
// En una aplicación real, DEBES implementar la lógica de conexión real 
// utilizando una cadena de conexión y manejo de errores adecuado.
func initDBConnection() (*pgx.Conn, error) {
	log.Println("🚨 ADVERTENCIA: Usando mock de conexión DB. Reemplaza esta función con la lógica de conexión PostgreSQL real.")
	// Devolvemos nil, nil para simular una inicialización exitosa
	// sin una conexión activa. Los handlers deben ser capaces de tolerar esto
	// o el usuario debe proveer la conexión real.
	return nil, nil 
}

func main() {
	router := gin.Default()
	
	// --- CONFIGURACIÓN CORS ---
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// --------------------------

	// 1. Inicializar la Conexión a la Base de Datos (MOCK)
	db, err := initDBConnection() 
	if err != nil {
		log.Fatal("Error al inicializar la base de datos:", err)
	}
	
	// 2. Inicializar la Caché de Ubicaciones y el Broadcaster
	// Le pasamos la conexión DB (incluso si es mock) para que InitCache pueda ejecutarse.
	if err := location.InitCache(db); err != nil {
		log.Fatal("Error al inicializar la caché de ubicaciones:", err)
	}
	
	// Grupo de rutas API v1
	v1 := router.Group("/v1")
	{
		// --- RUTAS DE GEOLOCALIZACIÓN ---
		// Inyectamos la conexión DB a los handlers para que puedan acceder a PostgreSQL
		v1.POST("/location", location.LocationPostHandler(db))
		v1.GET("/location", location.LocationGetHandler(db))
		v1.GET("/ws/location", location.LocationWsHandler(db)) 
		
		// --- GRUPO DE RUTAS DE AUTENTICACIÓN WEB3 (SIWE) ---
		auth := v1.Group("/auth")
		{
			auth.GET("/nonce", handlers.RequestNonce) 
			auth.POST("/verify", handlers.VerifyAndAuthenticate) 
		}
	}
	
	// Iniciar el servidor en el puerto 8081
	err = router.Run(":8081") 
	if err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}