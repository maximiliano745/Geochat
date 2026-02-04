import { User } from '../models/User';
import { Vault } from './Vault'; // Importamos tu motor de blockchain

export class Payment {
  static async process(from: User, to: User, amount: number): Promise<boolean> {
    console.log(`[Blockchain] Iniciando transacción de ${amount} TKG...`);

    try {
      // 1. Verificamos identidad soberana antes de mover fondos
      const isVerified = Vault.verifyProofOfWork(from.did);
      if (!isVerified) return false;

      // 2. Aplicamos tu lógica de sociedad 60/40
      const sharePartner = amount * 0.40; // Capitalista
      const shareOwner = amount * 0.60;   // Tu parte (Idea/Lab)

      // 3. Ejecución directa: GeoChat no es intermediario, es el puente.
      // Se envían los fondos desde el Vault del usuario hacia los destinos.
      Vault.transferTKG(from.did, to.did, shareOwner); 
      Vault.transferTKG(from.did, "did:geochat:partner-capital", sharePartner);

      console.log(`✅ Pago procesado: ${shareOwner} TKG al nodo y ${sharePartner} TKG al inversor.`);
      return true;

    } catch (error) {
      console.error("Error en la transacción:", error);
      return false;
    }
  }
}