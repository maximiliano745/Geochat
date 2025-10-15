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
)

// Global DB connection pool variable (Aunque pgx v5 prefiere usar la conexión directamente,
// usaremos esta variable para simplificar el ejemplo de check-up)
var db *pgx.Conn

func main() {
	// --- 1. Inicialización de la Base de Datos ---
	if err := initDB(); err != nil {
		log.Fatalf("FATAL: No se pudo inicializar la base de datos: %v", err)
	}
	defer db.Close(context.Background())
	log.Println("INFO: Conexión a PostgreSQL establecida exitosamente.")

	// --- 2. Configuración del Router ---
	router := mux.NewRouter()

	// Rutas principales
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")

	// Middleware de registro simple (opcional pero útil)
	loggedRouter := loggingMiddleware(router)

	// --- 3. Inicio del Servidor ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto configurado en docker-compose.yml
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("INFO: Servidor backend Go iniciado en http://localhost%s", addr)

	// El servidor se ejecuta en una goroutine separada para no bloquear
	err := http.ListenAndServe(addr, loggedRouter)
	if err != nil {
		log.Fatalf("FATAL: Error al iniciar el servidor: %v", err)
	}
}

// initDB se encarga de configurar y verificar la conexión a PostgreSQL.
func initDB() error {
	// Usamos el nombre del servicio Docker 'db' como el host de la conexión.
	// Las credenciales provienen del archivo docker-compose.yml.
	connStr := "postgresql://postgres:postgres@db:5432/geochat"
	
	// Intenta conectar con reintentos. Esto es crucial en Docker Compose,
	// ya que el servicio 'db' puede tardar un poco más en arrancar que 'app'.
	maxAttempts := 5
	for i := 0; i < maxAttempts; i++ {
		conn, err := pgx.Connect(context.Background(), connStr)
		if err == nil {
			// Si la conexión fue exitosa, la asignamos a la variable global.
			db = conn
			return nil
		}
		log.Printf("ADVERTENCIA: Falló el intento %d de conexión a DB: %v. Reintentando en 3 segundos...", i+1, err)
		time.Sleep(3 * time.Second)
	}
	return fmt.Errorf("falló la conexión a PostgreSQL después de %d intentos", maxAttempts)
}

// --- Handlers ---

// homeHandler responde con un mensaje simple.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Bienvenido a la API de GeoChat. El backend Go está corriendo.")
}

// healthCheckHandler verifica la conexión a la base de datos.
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "ERROR: La conexión a la base de datos no está inicializada.", http.StatusInternalServerError)
		return
	}

	// Ping a la base de datos para verificar que está viva
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

// --- Middleware ---

// loggingMiddleware registra cada solicitud.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Solicitud recibida: %s %s desde %s", r.Method, r.RequestURI, r.RemoteAddr)
		// Llama al siguiente handler
		next.ServeHTTP(w, r)
	})
}
