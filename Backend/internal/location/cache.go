package location

import (
	"context"
	"fmt" // Agregado para fmt.Errorf en InitCache
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

// UserLocation representa la estructura de ubicación del usuario.
// El tag JSON 'userId' es crucial para la comunicación con el frontend.
type UserLocation struct {
	UserID    string    `json:"userId"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt time.Time `json:"updatedAt"` // Marca de tiempo del último reporte
}

// LocationCache maneja el almacenamiento en memoria (caché) de las ubicaciones.
type LocationCache struct {
	data map[string]UserLocation // Mapa: UserID -> UserLocation
	mu   sync.RWMutex            // Mutex para proteger el acceso concurrente
}

// Broadcaster maneja el sistema Pub/Sub para distribuir las actualizaciones 
// de ubicación a todos los clientes de WebSocket conectados.
type Broadcaster struct {
	// Broadcast es el canal por donde el POST handler inyecta nuevas ubicaciones.
	Broadcast chan UserLocation 
	
	// Clients almacena los canales de salida de cada cliente WS activo.
	Clients map[chan UserLocation]bool
	
	mu sync.RWMutex // Mutex para proteger el mapa de clientes
}

// Instancias globales del sistema de gestión de ubicación.
var (
	locationCache = &LocationCache{
		data: make(map[string]UserLocation),
	}
	broadcaster = &Broadcaster{
		Broadcast: make(chan UserLocation),
		Clients:   make(map[chan UserLocation]bool),
	}
)

// InitCache inicializa la caché en memoria:
// 1. Carga los datos existentes desde la base de datos (PostgreSQL).
// 2. Lanza el loop principal del Broadcaster en una goroutine.
func InitCache(db *pgx.Conn) error {
	log.Println("INFO: Iniciando la carga de caché de ubicaciones y Broadcaster...")
	
	// 1. Lanza el Broadcaster para que escuche el canal Broadcast.
	go broadcaster.Start()
	
	// 2. Consulta todas las ubicaciones existentes en la DB.
	query := `SELECT user_id, latitude, longitude, created_at FROM user_locations`
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		return fmt.Errorf("error al consultar ubicaciones para caché: %w", err)
	}
	defer rows.Close()

	// 3. Carga los datos en la caché en memoria.
	locationCache.mu.Lock()
	defer locationCache.mu.Unlock()

	for rows.Next() {
		var loc UserLocation
		var createdAt time.Time
		
		if err := rows.Scan(&loc.UserID, &loc.Latitude, &loc.Longitude, &createdAt); err != nil {
			log.Printf("ADVERTENCIA: Error al escanear fila de ubicación: %v", err)
			continue
		}
		// Usamos el timestamp de la DB como la marca de la última actualización.
		loc.UpdatedAt = createdAt
		locationCache.data[loc.UserID] = loc
	}

	log.Printf("INFO: Caché de ubicaciones inicializada con %d entradas.", len(locationCache.data))
	return nil
}

// Start es el loop principal del Broadcaster. 
// Distribuye cualquier mensaje recibido en 'Broadcast' a todos los clientes suscritos.
func (b *Broadcaster) Start() {
	for {
		// Bloquea hasta que una nueva ubicación esté lista para ser distribuida.
		newLocation := <-b.Broadcast

		b.mu.RLock()
		
		// Distribuye la ubicación a cada cliente conectado.
		for clientChan := range b.Clients {
			select {
			case clientChan <- newLocation:
				// Envío asíncrono exitoso.
			default:
				// El canal del cliente está bloqueado o cerrado (el cliente se ha ido).
				// Necesitamos liberarlo y remover al cliente.
				b.mu.RUnlock() 
				b.RemoveClient(clientChan)
				b.mu.RLock() // Re-adquirir el bloqueo de lectura para continuar el loop de manera segura.
			}
		}
		b.mu.RUnlock()
	}
}

// AddClient registra un nuevo cliente WebSocket al sistema Broadcaster.
func (b *Broadcaster) AddClient(clientChan chan UserLocation) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Clients[clientChan] = true
	log.Printf("DEBUG: Nuevo cliente WS conectado. Total: %d", len(b.Clients))
}

// RemoveClient elimina un cliente del sistema Broadcaster y cierra su canal de comunicación.
func (b *Broadcaster) RemoveClient(clientChan chan UserLocation) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.Clients[clientChan]; ok {
		delete(b.Clients, clientChan)
		// Cierra el canal para liberar los recursos.
		close(clientChan)
		log.Printf("DEBUG: Cliente WS desconectado. Total: %d", len(b.Clients))
	}
}

// UpdateCache actualiza la ubicación en la caché y notifica a todos los clientes 
// a través del Broadcaster.
func UpdateCache(loc UserLocation) {
	// 1. Aseguramos una marca de tiempo actual para la caché.
	loc.UpdatedAt = time.Now()
	
	// 2. Actualizar la caché de forma segura.
	locationCache.mu.Lock()
	locationCache.data[loc.UserID] = loc
	locationCache.mu.Unlock()
	
	// 3. Publicar el cambio. Esto no bloquea porque el Broadcaster.Start() 
	// se ejecuta en su propia goroutine.
	broadcaster.Broadcast <- loc
}

// GetAllLocations devuelve todas las ubicaciones almacenadas actualmente en la caché.
func GetAllLocations() []UserLocation {
	locationCache.mu.RLock()
	defer locationCache.mu.RUnlock()

	// Copia el mapa a una slice para devolverlo sin exponer la estructura interna.
	locations := make([]UserLocation, 0, len(locationCache.data))
	for _, loc := range locationCache.data {
		locations = append(locations, loc)
	}
	return locations
}