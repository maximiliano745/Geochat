import { Vault } from '../services/Vault';
import { IA_Friend } from '../services/IA_Friend';

interface TrustProfile {
  did: string;
  reputationScore: number;
  uptimeHistory: number[];
  successfulBarters: number;
}

export class TrustEngine {
  constructor(
    private vaultService: typeof Vault,
    private aiService: typeof IA_Friend // <--- Estás usando inyección de dependencias
  ) { }

  evaluateForHardwareUpgrade(profile: TrustProfile): boolean {
    const isReliable = profile.reputationScore > 85;
    const isActive = profile.uptimeHistory.every(up => up > 0.8);
    return isReliable && isActive;
  }

  grantMicroLoan(merchantId: string) {
    this.vaultService.issueVoucher(merchantId, "HARDWARE_KIT_ELITE");

    // SOLUCIÓN: Usa 'this.aiService' en lugar de 'IA_Friend' directamente
    this.aiService.notify(merchantId, "¡Felicidades! Tu reputación ha construido tu infraestructura");
  }

  applyTeslaReputationBonus(profile: TrustProfile, isTeslaModeActive: boolean) {
    if (isTeslaModeActive) {
      profile.reputationScore += 5;
      // Aquí ya lo tenías bien usando 'this.aiService'
      this.aiService.notify(profile.did, "Tu generosidad (Modo Tesla) ha fortalecido tu perfil de confianza.");
    }
  }
}