
import { User } from '../models/User';

export class Payment {
  static async process(from: User, to: User, amount: number): Promise<boolean> {
    console.log(`Procesando pago de ${amount} TKG...`);
    // Aquí iría la lógica para interactuar con el Smart Contract
    // Recuerda que GeoChat no toca las llaves privadas directamente.
    return true; 
  }
}
