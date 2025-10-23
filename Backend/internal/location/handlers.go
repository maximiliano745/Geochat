package location

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	//"time" // 👈 Se mantiene el 'time' para mayor claridad y posible uso futuro.

	"github.com/jackc/pgx/v5"
)

// LocationPostHandler es una CLOUSURE que retorna un http.HandlerFunc,
// permitiendo inyectar la dependencia de la base de datos (db).
func LocationPostHandler(db *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Método HTTP
		if r.Method != "POST" {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		// 2. Decodificación del JSON
		var locData LocationData
		if err := json.NewDecoder(r.Body).Decode(&locData); err != nil {
			log.Printf("ERROR: Falló la decodificación de JSON: %v", err)
			http.Error(w, "Cuerpo de solicitud JSON inválido.", http.StatusBadRequest)
			return
		}

		// 3. Validación de datos
		if locData.UserID == "" || locData.Latitude == 0.0 || locData.Longitude == 0.0 {
			http.Error(w, "Datos de ubicación incompletos (requiere user_id, latitud, longitud).", http.StatusBadRequest)
			return
		}

		// 4. Guardar en la Base de Datos
		// Se usa r.Context() para manejar posibles timeouts/cancelaciones de la solicitud HTTP.
		if err := insertLocation(r.Context(), db, locData); err != nil {
			log.Printf("ERROR: Falló la inserción en DB: %v", err)
			http.Error(w, "Error interno del servidor al guardar la ubicación.", http.StatusInternalServerError)
			return
		}

		// 5. Respuesta exitosa
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // 201 Created indica que se creó un nuevo recurso
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Ubicación recibida y guardada.",
			"user_id": locData.UserID,
		})
	}
}

// insertLocation ejecuta la consulta SQL para guardar la ubicación.
func insertLocation(ctx context.Context, db *pgx.Conn, loc LocationData) error {
	insertSQL := `
        INSERT INTO user_locations (user_id, latitude, longitude) 
        VALUES ($1, $2, $3);
    `
	// El campo created_at se asigna automáticamente por el DEFAULT de la tabla.
	_, err := db.Exec(ctx, insertSQL, loc.UserID, loc.Latitude, loc.Longitude)

	if err != nil {
		log.Printf("DB ERROR: Error al ejecutar insert: %v", err)
	}
	return err
}