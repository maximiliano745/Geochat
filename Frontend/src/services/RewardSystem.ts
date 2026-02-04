// src/services/RewardSystem.ts
import { Vault } from './Vault';
import { IA_Friend } from './IA_Friend';

export class RewardSystem {
  // Definimos que al llegar a 1000 Wh (1kWh) compartidos, se emite un premio
  private static readonly MILESTONE_TESLA_WH = 1000; 

  static checkAndIssueReward(ownerDid: string, totalSharedWh: number) {
    if (totalSharedWh >= this.MILESTONE_TESLA_WH) {
      
      // 1. Ejecución en Blockchain: Emitir el Voucher para un nuevo kit
      Vault.issueVoucher(ownerDid, "KIT-EXPANSION-2026");
      
      // 2. Notificación empática
      IA_Friend.notify(ownerDid, 
        `🎉 ¡Hito de Filantropía alcanzado! Has compartido ${totalSharedWh}Wh. ` +
        `Se ha emitido un Voucher en tu Vault para la expansión de tu nodo.`
      );
      
      return true;
    }
    return false;
  }
}