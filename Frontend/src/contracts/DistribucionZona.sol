// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract DistribucionZona {
    address public inversor;
    address public laboratorioGeoChat;
    uint256 public costoHardwareUSD;
    uint256 public acumuladoRecuperadoUSD;
    bool public hardwareAmortizado = false;

    // Al recibir un pago, especificamos si fue por Modo Tesla
    function repartirIngreso(uint256 montoPAXG, uint256 valorEnUSD, bool esModoTesla) public {
        if (esModoTesla) {
            // El ingreso filantrópico podría ir a un fondo de comunidad
            // o repartirse directo sin afectar la amortización técnica.
            enviarOro(laboratorioGeoChat, montoPAXG); 
            return; 
        }

        if (!hardwareAmortizado) {
            enviarOro(inversor, montoPAXG);
            acumuladoRecuperadoUSD += valorEnUSD;
            if (acumuladoRecuperadoUSD >= costoHardwareUSD) {
                hardwareAmortizado = true;
            }
        } else {
            enviarOro(inversor, montoPAXG / 2);
            enviarOro(laboratorioGeoChat, montoPAXG / 2);
        }
    }

    function enviarOro(address destinatario, uint256 monto) private {
        // Lógica para transferir PAXG (ERC20)
    }
}