# Proyecto Supply Chain - Hyperledger Fabric

Este repositorio contiene la implementación de una red de trazabilidad para suministros usando Hyperledger Fabric y una API REST en Node.js.

## Estructura del Proyecto
- `fabric-samples/test-network`: Configuración de la red.
- `fabric-samples/asset-transfer-basic/chaincode-go`: Smart Contract (Chaincode) en Go.
- `fabric-samples/test-network/api`: API REST en Node.js para interactuar con la red.

## Requisitos
- Docker y Docker Compose.
- Go 1.20+.
- Node.js v18+.
- WSL2 (si se ejecuta en Windows).

## Guía de Inicio Rápido
1. **Levantar la red:**
   ```bash
   cd fabric-samples/test-network
   ./network.sh down
   ./network.sh up createChannel -c supplychannel -ca -s couchdb

Desplegar el Chaincode:
   ./network.sh deployCC -ccn supplycc -ccp ../asset-transfer-basic/chaincode-go/ -ccl go

Levantar la API:
cd api
npm install
node app.js

Endpoints Principales
- GET /productos: Lista todos los productos.
- POST /productos: Registra un nuevo lote.
- GET /productos/:id/historial: Muestra la trazabilidad completa.
