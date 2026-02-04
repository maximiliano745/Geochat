import { BateryGuards } from './logica/BateryGuards';
import { EnergyRoaming } from './logica/EnergyRoaming';
import { Payment } from './services/LogicaTransacciones';
import { IA_Friend } from './services/IA_Friend';
import { RewardSystem } from './services/RewardSystem';
import { User } from './models/User';

async function simularLaboratorioGeoChat() {
    console.log("🚀 --- INICIANDO SIMULACIÓN GEOCHAT 2026 --- 🚀\n");

    // 1. Configuración de Usuarios (DIDs Soberanos)
    const dueño: User = {
        id: "owner-99",
        name: "Dueño",
        did: "did:geochat:owner-99",
        energyVault: { accumulatedPE: 100, locationDID: "home-station-01" }
    };

    const visitante: User = {
        id: "visitor-77",
        name: "Visitante",
        did: "did:geochat:visitor-77",
        energyVault: { accumulatedPE: 0, locationDID: "mobile-unit-02" } // Sin energía previa
    };

    // 2. Inicialización de Servicios
    const guardias = new BateryGuards();
    const roaming = new EnergyRoaming(Payment, IA_Friend);

    // 3. Situación de Hardware
    console.log("📡 [Hardware] El nodo está al 20% de batería...");
    guardias.isTeslaModeActive = true; // El usuario elige ser filántropo
    
    guardias.manageTraffic({
        level: 20,
        isCharging: false,
        activeConnections: 15
    });

    console.log("\n🚗 [Roaming] El visitante solicita carga de 500Wh...");
    
    // 4. Ejecución del Roaming y Reparto 60/40
    await roaming.handleChargeRequest(visitante, dueño, 500);

    console.log("\n--- SIMULACIÓN FINALIZADA ---");
    // Añade esto al final de tu función simularLaboratorioGeoChat()

    // 5. Generación de Reporte de Impacto
    // Simulamos que el nodo ha tenido actividad acumulada
    IA_Friend.generateImpactReport(dueño.did, {
        totalTKG: 0.5,           // Total generado en el día
        connections: 45,         // Personas que pasaron por el nodo
        energyShared: 1200,      // Wh totales transferidos
        teslaMode: guardias.isTeslaModeActive 
    });

    // 6. Verificación de Recompensas (Cierre del ciclo)
    const energiaTotalHoy = 1200; // Wh acumulados en la sesión
    
    if (guardias.isTeslaModeActive) {
        RewardSystem.checkAndIssueReward(dueño.did, energiaTotalHoy);
    }

    console.log("\n🚀 --- LABORATORIO GEOCHAT COMPLETADO CON ÉXITO --- 🚀");
}

simularLaboratorioGeoChat();
