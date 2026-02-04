import { User } from '../models/User';
import { Payment } from '../services/LogicaTransacciones'; // Asumo que este conecta con tu Vault
import { IA_Friend } from '../services/IA_Friend';

export class EnergyRoaming {
  private static readonly CURRENT_ENERGY_PRICE = 0.0001; // Precio base en TKG/PAXG

  constructor(
    private paymentService: typeof Payment,
    private aiService: typeof IA_Friend
  ) {}

  handleChargeRequest(visitor: User, nodeOwner: User, amountWh: number) {
    console.log(`[EnergyRoaming] ${visitor.did} solicita ${amountWh}Wh de ${nodeOwner.did}`);

    // 1. Verificación de Puntos de Energía (PE) generados en casa
    if (visitor.energyVault.accumulatedPE >= amountWh) {
      
      // Transferencia directa de valor energético (Trueque Digital)
      visitor.energyVault.accumulatedPE -= amountWh;
      nodeOwner.energyVault.accumulatedPE += amountWh;
      
      this.aiService.notify(visitor.did, 
        `🔋 ¡Roaming Exitoso! Usaste ${amountWh} PE de tu reserva soberana.`);
      
    } else {
      // 2. Si no hay PE, se activa el flujo económico (TKG/PAXG)
      const cost = amountWh * EnergyRoaming.CURRENT_ENERGY_PRICE;
      
      // Aquí el paymentService dispara la lógica que respeta el 60/40
      this.paymentService.process(visitor, nodeOwner, cost);
      
      this.aiService.notify(visitor.did, 
        `⚠️ PE insuficientes. Se han procesado ${cost} TKG por la carga.`);
    }
  }
}