
import { Vault } from '../services/Vault';

interface TKGTransaction {
  fromNode: string;
  toNode: string;
  amount: number;
  serviceType: 'INTERNET_DATA' | 'DELIVERY_FEE' | 'BARTER_MATCH';
}

export class TKGManager {
  // Inyectamos el Vault para que la lógica de TKG pueda usarlo
  constructor(private vaultService: typeof Vault) {}

  processService(tx: TKGTransaction) {
    // Usamos el servicio inyectado
    if (this.vaultService.verifyProofOfWork(tx.fromNode)) {
       this.vaultService.transferTKG(tx.fromNode, tx.toNode, tx.amount);
       console.log(`IA Friend: Pago de ${tx.serviceType} liquidado instantáneamente.`);
    } else {
       console.error("IA Friend: Error de seguridad. El nodo no pudo verificar su identidad.");
    }
  }
}