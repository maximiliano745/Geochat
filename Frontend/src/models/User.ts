// src/types/User.ts

export interface EnergyVault {
  accumulatedPE: number;
  locationDID: string;
}

export interface User {
    id: string;
    name: string;
    did: string;
    energyVault: {
        accumulatedPE: number;
        locationDID: string;
    };
    isTeslaModeActive?: boolean; // El usuario lo activa manualmente
}
