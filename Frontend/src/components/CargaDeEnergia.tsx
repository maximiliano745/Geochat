
// src/components/BotonCarga.tsx del EnergiRoadming.ts

import { EnergyRoaming } from '../logica/EnergyRoaming'; // <--- Importas la receta
import { Payment } from '../services/LogicaTransacciones';
import { IA_Friend } from '../services/Notificacion-IA_Friend';

// Instancias FUERA de la lógica para romper el ciclo
const roamingEngine = new EnergyRoaming(Payment, IA_Friend);

export const BotonCarga = ({ visitor, owner, amount }: any) => {
  return (
    <button onClick={() => roamingEngine.handleChargeRequest(visitor, owner, amount)}>
      Cargar Energía
    </button>
  );
};

