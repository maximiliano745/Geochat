import { IA_Friend } from '../services/IA_Friend';

// Definimos la estructura de los datos de energía
interface PowerStats {
  level: number;
  isCharging: boolean;
  activeConnections: number;
}

export class EnergyManager {
  // El ID del dueño del nodo para que la IA sepa a quién hablarle
  private ownerId: string = "owner-123";

  // Definimos el método donde vive la lógica
  manageTraffic(stats: PowerStats) {
    
    // 1. Lógica de ahorro de batería
    if (stats.level < 30 && !stats.isCharging) {
      this.throttleConnections(0.5); // Ahora 'this' funciona porque está dentro de la clase
      this.limitToGeoChatOnly(true);
      
      IA_Friend.notify(this.ownerId, "Modo Ahorro: Solo procesando pagos y red interna.");
    }
    
    // 2. Lógica de incentivos por alto tráfico
    if (stats.activeConnections > 10) {
      this.increaseIncentiveRequest(); 
    }
  }

  // Definimos los métodos que antes faltaban para que 'this' pueda encontrarlos
  private throttleConnections(factor: number) {
    console.log(`[Hardware] Reduciendo velocidad al ${factor * 100}%`);
  }

  private limitToGeoChatOnly(status: boolean) {
    console.log(`[Red] Modo exclusivo GeoChat: ${status}`);
  }

  private increaseIncentiveRequest() {
    console.log(`[Economía] Solicitando más TKG por alto desgaste de hardware.`);
  }
}