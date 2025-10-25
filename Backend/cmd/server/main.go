package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"

	// 🚨 IMPORTACIÓN CLAVE: Debes ajustar "GeoChat/Backend" si tu módulo Go tiene otro nombre
	"Geochat/Backend/internal/location"
)

// Global DB connection pool variable
var db *pgx.Conn

func main() {
	// --- 1. Inicialización de la Base de Datos ---
	if err := initDB(); err != nil {
		log.Fatalf("FATAL: No se pudo inicializar la base de datos: %v", err)
	}
	defer db.Close(context.Background())
	log.Println("INFO: Conexión a Post establecida exitosamente.")

	// --- 2. Configuración del Router ---
	router := mux.NewRouter()

	// Rutas principales
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")

	// 🟢 REGISTRO DEL HANDLER DE UBICACIÓN
	// Se pasa la conexión 'db' al handler para que pueda ejecutar las consultas SQL.
	router.HandleFunc("/v1/location", location.LocationPostHandler(db)).Methods("POST")

	// Middleware de registro simple (opcional pero útil)
	loggedRouter := loggingMiddleware(router)

	// --- 3. Inicio del Servidor ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto configurado en docker-compose.yml
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("INFO: Servidor backend Go iniciado en http://localhost%s", addr)

	err := http.ListenAndServe(addr, loggedRouter)
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
				user_id TEXT NOT NULL,
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

// --- Handlers --- (Se mantienen)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Bienvenido a la API de GeoChat. El backend Go está corriendo.")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "ERROR: La conexión a la base de datos no está inicializada.", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		log.Printf("ERROR: Falló el ping a la base de datos: %v", err)
		http.Error(w, fmt.Sprintf("Estado: DB Desconectada. Error: %v", err), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Estado: OK. El servidor Go y la DB están operativos.")
}

// --- Middleware --- (Se mantiene)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Solicitud recibida: %s %s desde %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}