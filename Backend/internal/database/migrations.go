package database

import (
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/jackc/pgx/v5/stdlib" // Importa el driver PostgreSQL
)

// Las credenciales provienen de tu docker-compose.yml
const (
	dbHost = "db"
	dbPort = 5432
	dbUser = "postgres"
	dbPassword = "postgres"
	dbName = "geochat"
)

// ConnectDB construye la cadena de conexión y la abre.
func ConnectDB() *sql.DB {
	// Cadena de conexión DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Abrir la conexión a la base de datos
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Error al abrir la conexión a la base de datos: %v", err)
	}

	// Probar la conexión (ping)
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos (Host: %s): %v", dbHost, err)
	}

	fmt.Println("Conexión a PostgreSQL establecida con éxito.")
	return db
}

// RunMigrations es una función placeholder donde iría la lógica para crear tablas.
func RunMigrations(db *sql.DB) {
	fmt.Println("Ejecutando migraciones...")
	
	// --- Ejemplo: Creación de una tabla simple ---
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		email VARCHAR(100) NOT NULL UNIQUE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error al ejecutar la migración (creación de tabla 'users'): %v", err)
	}
	fmt.Println("Tabla 'users' creada o ya existente.")
	// ---------------------------------------------
	
	fmt.Println("Migraciones completadas.")
}