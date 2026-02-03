// src/logica/EnergyRoaming.ts
import { User } from '../models/User';
import { Payment } from '../services/LogicaTransacciones';
import { IA_Friend } from '../services/Notificacion-IA_Friend';

export class EnergyRoaming { // <--- Solo exportas la clase
  constructor(
    private paymentService: typeof Payment,
    private aiService: typeof IA_Friend
  ) {}

  handleChargeRequest(visitor: User, nodeOwner: User, amountWh: number) {
     // ... tu lógica de PE y Pago
  }
}