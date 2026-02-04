import { Vault } from '../services/Vault';
import { IA_Friend } from '../services/IA_Friend';

export class RewardSystem {
  private static readonly MILESTONE_TESLA_WH = 1000; // Hito: 1kWh compartido

  static checkAndIssueReward(ownerDid: string, totalSharedWh: number) {
    if (totalSharedWh >= this.MILESTONE_TESLA_WH) {
      console.log(`\n[RewardSystem] ¡Hito alcanzado! Generando recompensa...`);
      
      // Usamos TU función del Vault para emitir el voucher
      Vault.issueVoucher(ownerDid, "KIT-EXPANSION-PRO");
      
      IA_Friend.notify(ownerDid, 
        "🎁 ¡Felicidades! Tu generosidad ha desbloqueado un Voucher de Expansión. " +
        "Ya puedes reclamar hardware adicional en la red GeoChat."
      );
    }
  }
}