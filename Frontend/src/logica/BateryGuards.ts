// src/logica/BateryGuards.ts
import { IA_Friend } from '../services/IA_Friend';

interface PowerStats {
  level: number;
  isCharging: boolean;
  activeConnections: number;
}

export class BateryGuards {
  private ownerDid: string = "did:geochat:owner-99";
  public isTeslaModeActive: boolean = false; 
  private connectionsHandledInTeslaMode: number = 0;

  manageTraffic(stats: PowerStats): void {
    // Lógica del Modo Tesla (Filantropía manual)
    if (this.isTeslaModeActive) {
      IA_Friend.notify(this.ownerDid, "Modo Tesla activo: Priorizando la comunidad.");
      if (stats.activeConnections > 10) {
        this.increaseIncentiveRequest(stats.activeConnections, true);
      }
      return; 
    }

    // Protección de hardware (Acuerdo 60/40)
    if (stats.level < 30 && !stats.isCharging) {
      this.throttleConnections(0.5); 
      this.limitToGeoChatOnly(true); 
      IA_Friend.notify(this.ownerDid, "Modo Ahorro: Protegiendo hardware.");
    }
    
    if (stats.activeConnections > 10) {
      this.increaseIncentiveRequest(stats.activeConnections, false); 
    }
  }

  private throttleConnections(factor: number): void {
    console.log(`[BateryGuards] Reduciendo carga al ${factor * 100}%`);
  }

  private limitToGeoChatOnly(enabled: boolean): void {
    console.log(`[BateryGuards] Firewall restringido: ${enabled}`);
  }

  private increaseIncentiveRequest(connections: number, isTesla: boolean): void {
    const BASE_TARIFF = 0.001; 
    const SURGE_THRESHOLD = 10;
    const extraConnections = connections - SURGE_THRESHOLD;
    const finalTariff = BASE_TARIFF + (extraConnections * 0.0001);

    if (isTesla) {
      this.connectionsHandledInTeslaMode += extraConnections;
      if (this.connectionsHandledInTeslaMode >= 50) {
        IA_Friend.notify(this.ownerDid, `🌟 ¡Impacto Tesla! Has procesado ${this.connectionsHandledInTeslaMode} conexiones extra.`);
        this.connectionsHandledInTeslaMode = 0; 
      }
    }
    console.log(`[Economía GeoChat] Tarifa: ${finalTariff} PAXG. ¿Filantropía?: ${isTesla}`);
  }
}