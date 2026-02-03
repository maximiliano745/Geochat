// src/types/User.ts

export interface EnergyVault {
  accumulatedPE: number;
  locationDID: string;
}

export interface User {
  id: string;           // ID interno o público
  did: string;          // Decentralized Identifier (Sovereign Identity)
  name: string;
  energyVault: EnergyVault; // Aquí conectamos con tu interfaz anterior
  // Nota: No incluimos la clave privada aquí por seguridad

  isTeslaModeActive?: boolean; // El usuario lo activa manualmente
}
