import { IA_Friend } from '../services/IA_Friend';

interface PowerStats {
  level: number;
  isCharging: boolean;
  activeConnections: number;
}

export class BateryGuards {
  private ownerDid: string = "did:geochat:owner-99";

  manageTraffic(stats: PowerStats): void {
    // 1. Protección por batería baja
    if (stats.level < 30 && !stats.isCharging) {
      this.throttleConnections(0.5); 
      this.limitToGeoChatOnly(true); 
      IA_Friend.notify(this.ownerDid, "Modo Ahorro: Solo procesando pagos y red interna.");
    }
    
    // 2. Protección por sobreesfuerzo
    if (stats.activeConnections > 10) {
      // Pasamos el número de conexiones para calcular el incentivo
      this.increaseIncentiveRequest(stats.activeConnections); 
    }
  }

  private throttleConnections(factor: number): void {
    console.log(`[BateryGuards] Reduciendo carga de red al ${factor * 100}%`);
  }

  private limitToGeoChatOnly(enabled: boolean): void {
    console.log(`[BateryGuards] Firewall restringido a red interna: ${enabled}`);
  }

  // MÉTODO ACTUALIZADO:
  private increaseIncentiveRequest(connections: number): void {
    const BASE_TARIFF = 0.001; 
    const SURGE_THRESHOLD = 10;
    
    // Lógica: Por cada conexión extra sobre el umbral, sumamos un pequeño plus por desgaste
    const extraConnections = connections - SURGE_THRESHOLD;
    const finalTariff = BASE_TARIFF + (extraConnections * 0.0001);

    console.log(`[Economía GeoChat] Tráfico alto (${connections} hilos). Tarifa ajustada a: ${finalTariff} PAXG por desgaste.`);
    
    // Aquí es donde el Vault entraría en acción en el futuro:
    // Vault.updateSurgePricing(this.ownerDid, finalTariff);
  }
}