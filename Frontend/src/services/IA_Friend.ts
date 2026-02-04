// src/services/IA_Friend.ts

export class IA_Friend {
  static notify(did: string, message: string): void {
    console.log(`\x1b[36m[IA Friend]\x1b[0m Para ${did}: ${message}`);
  }

  static generateImpactReport(did: string, stats: { 
    totalTKG: number, 
    connections: number, 
    energyShared: number,
    teslaMode: boolean 
  }) {
    console.log(`\n\x1b[35m--- 📊 REPORTE DE IMPACTO GEOCHAT --- \x1b[0m`);
    console.log(`Hola ${did}, este es el resumen de tu nodo hoy:`);
    
    const miParte = stats.totalTKG * 0.60;
    
    console.log(`* 💰 Ganancia neta (60%): ${miParte.toFixed(4)} TKG`);
    console.log(`* ⚡ Energía distribuida: ${stats.energyShared} Wh`);
    console.log(`* 🤝 Conexiones gestionadas: ${stats.connections}`);

    if (stats.teslaMode) {
        console.log(`\x1b[32m* 🌟 NOTA FILANTRÓPICA: Gracias a tu Modo Tesla, la comunidad creció un 15% más rápido hoy. ¡Eres un pilar de la red!\x1b[0m`);
    }
    console.log(`\x1b[35m------------------------------------\x1b[0m\n`);
  }
}