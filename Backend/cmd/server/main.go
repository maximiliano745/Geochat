package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors" // Para manejar CORS de manera sencilla con Gin
	"github.com/jackc/pgx/v5"

	// 🚨 IMPORTACIÓN CLAVE: Ajusta estas rutas a tu estructura de módulo
	"Geochat/Backend/internal/location"
	"Geochat/Backend/internal/handlers" // Asumiendo que los handlers de auth están aquí
)

// Global DB connection pool variable
var db *pgx.Conn

func main() {
	// --- 1. Inicialización de la Base de Datos ---
	// La función initDB intenta conectarse a la DB. Si falla, termina el programa aquí.
	if err := initDB(); err != nil {
		log.Fatalf("FATAL: No se pudo inicializar la base de datos: %v", err)
	}
	defer db.Close(context.Background())
	log.Println("INFO: Conexión a Post establecida exitosamente.")

	// 🟢 1.1. Inicialización del Sistema de Caché y Broadcaster (requerido por location handlers)
	// Como initDB fue exitoso, 'db' no es nil y puede ser usado por InitCache.
	if err := location.InitCache(db); err != nil {
		log.Fatalf("FATAL: No se pudo inicializar la caché y Broadcaster de ubicación: %v", err)
	}
	log.Println("INFO: Caché de ubicación y Broadcaster inicializados.")

	// --- 2. Configuración del Router (Gin) ---
	r := gin.Default() // Gin incluye Logging y Recovery middleware por defecto

	// 🟢 2.1. Configuración de CORS
	// Permitir solicitudes desde cualquier origen para desarrollo.
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // Permite cualquier origen (Frontend Vue, Mobile, etc.)
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Rutas principales
	r.GET("/", homeHandler)
	r.GET("/health", healthCheckHandler)

	// --- Rutas de API v1 ---
	v1 := r.Group("/v1")
	{
		// Rutas de Autenticación (Web3 SIWE)
		authGroup := v1.Group("/auth")
		{
			// GET /v1/auth/nonce?address=...
			authGroup.GET("/nonce", handlers.RequestNonce)
			// POST /v1/auth/verify
			authGroup.POST("/verify", handlers.VerifyAndAuthenticate)
		}

		// 🟢 REGISTRO DE HANDLERS DE UBICACIÓN
		// POST /v1/location (Actualización de ubicación)
		v1.POST("/location", location.LocationPostHandler(db)) 
		// GET /v1/location (Obtener todas las ubicaciones desde caché)
		v1.GET("/location", location.LocationGetHandler(db))
		
		// 🟢 REGISTRO DEL HANDLER DE WEBSOCKET
		// GET /v1/ws/location (Conexión WebSocket para actualizaciones en tiempo real)
		v1.GET("/ws/location", location.LocationWsHandler(db)) 
	}
	
	// --- 3. Inicio del Servidor ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("INFO: Servidor backend Go iniciado en https://localhost%s", addr)

	// Gin Run() inicia el servidor
	err := r.Run(addr) 
	if err != nil {
		log.Fatalf("FATAL: Error al iniciar el servidor: %v", err)
	}
}

// initDB se encarga de configurar, verificar y CREAR LA TABLA en PostgreSQL.
func initDB() error {
	// Usamos el nombre del servicio Docker 'db' como el host de la conexión.
	connStr := "postgresql://postgres:postgres@db:5432/geochat"
	maxAttempts := 5

	for i := 0; i < maxAttempts; i++ {
		conn, err := pgx.Connect(context.Background(), connStr)
		if err == nil {
			// Si la conexión fue exitosa, la asignamos a la variable global.
			db = conn

			// 🟢 1. Verificación y Creación de la Tabla 'user_locations'
			log.Println("INFO: Verificando/Creando la tabla 'user_locations'...")
			createTableSQL := `
			CREATE TABLE IF NOT EXISTS user_locations (
				id SERIAL PRIMARY KEY,
				user_id TEXT NOT NULL UNIQUE, 
				latitude DOUBLE PRECISION NOT NULL,
				longitude DOUBLE PRECISION NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
			);
			`
			// Ejecutamos la creación de la tabla
			if _, err := conn.Exec(context.Background(), createTableSQL); err != nil {
				// Si falla la creación de la tabla, cerramos la conexión y reportamos el error.
				conn.Close(context.Background())
				return fmt.Errorf("error al crear/verificar la tabla user_locations: %w", err)
			}
			log.Println("INFO: La tabla 'user_locations' está lista.")

			return nil // La conexión y la tabla están listas
		}
		log.Printf("ADVERTENCIA: Falló el intento %d de conexión a DB: %v. Reintentando en 3 segundos...", i+1, err)
		time.Sleep(3 * time.Second)
	}
	return fmt.Errorf("falló la conexión a PostgreSQL después de %d intentos", maxAttempts)
}

// --- Handlers Adaptados para Gin ---

// homeHandler maneja la solicitud GET /.
func homeHandler(c *gin.Context) {
	// Gin usa c.String para enviar una respuesta de texto
	c.String(http.StatusOK, "Bienvenido a la API de GeoChat. El backend Go (Gin) está corriendo.")
}

// healthCheckHandler maneja la solicitud GET /health.
func healthCheckHandler(c *gin.Context) {
	if db == nil {
		// Gin usa c.String en lugar de http.Error
		c.String(http.StatusInternalServerError, "ERROR: La conexión a la base de datos no está inicializada.")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		log.Printf("ERROR: Falló el ping a la base de datos: %v", err)
		c.String(http.StatusServiceUnavailable, fmt.Sprintf("Estado: DB Desconectada. Error: %v", err))
		return
	}

	c.String(http.StatusOK, "Estado: OK. El servidor Go (Gin) y la DB están operativos.")
}

// Nota: El middleware de logging de gorilla/mux fue eliminado ya que gin.Default() lo incluye.