package location

import "time"

// LocationData representa la estructura de los datos de ubicación recibidos.
// Las etiquetas `json` son cruciales para el parsing automático.
type LocationData struct {
	UserID    string    `json:"user_id"` // Identificador del usuario que envía la ubicación
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"` // Opcional, si la app móvil incluye la marca de tiempo
}

// Puedes añadir más modelos aquí, como estructuras de respuesta HTTP.
