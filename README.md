# Proyecto Capstone: Trazabilidad de Suministros con Hyperledger Fabric

**Laboratorio Capstone — Semanas 4 + 5 | E-Business / Blockchain Empresarial**  
**Estudiantes:** Daniela Díaz Caro · Estefanía Villalba Paternina  
**Docente:** Jaider Reyes Herazo

Red de blockchain permisionada para la gestión de trazabilidad en cadenas de suministro, con un Smart Contract (Chaincode) desarrollado en Go y una API REST en Node.js.

---

## Estructura del repositorio

```
proyecto-ebusiness/
├── chaincode-go/
│   ├── go.mod                  # Dependencias Go (fabric-contract-api-go v1.2.1)
│   └── chaincode/
│       └── supplycc.go         # Smart Contract — 10 funciones de Supply Chain
├── api/
│   ├── app.js                  # API REST Express + fabric-gateway
│   └── package.json
├── scripts/
│   ├── deploy-supplycc.sh      # Script de despliegue completo (11 pasos)
│   └── flujo-completo.sh       # Script de prueba del flujo Supply Chain
└── README.md
```

---

## Requisitos del sistema

### A. Docker y Docker Compose

```bash
sudo apt update
sudo apt install -y docker.io docker-compose
sudo usermod -aG docker $USER
# Cierra y vuelve a abrir la terminal después de este paso
```

### B. Go 1.21

```bash
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version   # debe mostrar go1.21.x
```

### C. Node.js y npm

```bash
sudo apt install -y nodejs npm
node --version   # debe ser v16 o superior
```

### D. Git y curl

```bash
sudo apt install -y git curl
```

---

## Instalación y puesta en marcha

### Paso 1 — Clonar el repositorio

```bash
git clone https://github.com/DannyDC29/proyecto-ebusiness.git
cd proyecto-ebusiness
```

### Paso 2 — Descargar los binarios e imágenes de Hyperledger Fabric

Esto descarga los ejecutables (`peer`, `orderer`, etc.) y las imágenes Docker de Fabric 2.5.4:

```bash
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.5.4 1.5.7
```

> El comando crea la carpeta `fabric-samples/` dentro de `proyecto-ebusiness/`.  
> Puede tardar varios minutos dependiendo de la conexión.

### Paso 3 — Levantar la red y el canal

```bash
cd fabric-samples/test-network
./network.sh down          # limpia cualquier estado anterior
./network.sh up createChannel -c supplychannel -ca -s couchdb
```

Deben aparecer **8 contenedores** corriendo:

```
peer0.org1.example.com
peer0.org2.example.com
orderer.example.com
ca_org1  ca_org2  ca_orderer
couchdb0  couchdb1
```

Verifica con: `docker ps`

### Paso 4 — Desplegar el Smart Contract (Chaincode)

```bash
cd ~/proyecto-ebusiness
chmod +x scripts/deploy-supplycc.sh
./scripts/deploy-supplycc.sh
```

El script ejecuta 11 pasos automáticos:

| Paso | Acción |
|------|--------|
| 0 | Verifica que la red esté activa |
| 1 | Descarga dependencias Go (`go mod vendor`) |
| 2 | Empaqueta el chaincode (`peer lifecycle chaincode package`) |
| 3 | Instala en `peer0.org1` |
| 4 | Instala en `peer0.org2` |
| 5 | Obtiene el Package ID |
| 6 | Aprueba como Org1MSP |
| 7 | Aprueba como Org2MSP |
| 8 | Verifica aprobaciones (`checkcommitreadiness`: debe mostrar `true/true`) |
| 9 | Hace commit al canal `supplychannel` |
| 10 | Inicializa el ledger con `InitLedger` (3 productos de ejemplo) |
| 11 | Consulta final con `ConsultarTodos` |

Al finalizar exitosamente verás los 3 productos del Caribe colombiano en el ledger.

### Paso 5 — Iniciar la API REST

En una **nueva terminal**:

```bash
cd ~/proyecto-ebusiness/api
npm install
node app.js
```

La API quedará disponible en `http://localhost:3000`.

---

## Guía de pruebas (Endpoints)

Abre otra terminal y usa los siguientes comandos:

### Verificar salud de la conexión

```bash
curl http://localhost:3000/health
```

Respuesta esperada:
```json
{"status":"OK","mspId":"Org1MSP","canal":"supplychannel","chaincode":"supplycc"}
```

### Consultar todos los productos

```bash
curl http://localhost:3000/productos | python3 -m json.tool
```

### Registrar un nuevo producto (POST)

```bash
curl -X POST http://localhost:3000/productos \
  -H "Content-Type: application/json" \
  -d '{"id":"PROD-100","nombre":"Cafe Especial","cantidad":500,"unidad":"kg","origen":"Huila","destino":"Puerto","temperatura":"18C"}'
```

### Consultar un producto por ID (GET)

```bash
curl http://localhost:3000/productos/PROD-100 | python3 -m json.tool
```

### Transferir a Org2 (PUT)

```bash
curl -X PUT http://localhost:3000/productos/PROD-100/transferir \
  -H "Content-Type: application/json" \
  -d '{"nuevoPropietario":"Org2MSP","detalle":"Envio terrestre Bogota-Barranquilla"}'
```

### Confirmar recepción (PUT)

```bash
curl -X PUT http://localhost:3000/productos/PROD-100/confirmar \
  -H "Content-Type: application/json" \
  -d '{"temperatura":"17C","detalle":"Lote recibido en buen estado"}'
```

### Actualizar cantidad (PUT)

```bash
curl -X PUT http://localhost:3000/productos/PROD-100/cantidad \
  -H "Content-Type: application/json" \
  -d '{"cantidad":480,"motivo":"Merma en transporte 4%"}'
```

### Marcar como entregado (PUT)

```bash
curl -X PUT http://localhost:3000/productos/PROD-100/entregar \
  -H "Content-Type: application/json" \
  -d '{"clienteFinal":"SuperAlimentos S.A.","evidencia":"Factura 2026-089"}'
```

### Consultar historial de trazabilidad (GET)

```bash
curl http://localhost:3000/productos/PROD-100/historial | python3 -m json.tool
```

### Filtrar por estado (GET)

```bash
curl http://localhost:3000/productos/estado/REGISTRADO | python3 -m json.tool
curl http://localhost:3000/productos/estado/ENTREGADO | python3 -m json.tool
```

---

## Flujo completo con el CLI del peer

Para ejecutar el flujo completo de Supply Chain usando directamente el peer (sin la API):

```bash
chmod +x scripts/flujo-completo.sh
./scripts/flujo-completo.sh
```

El script ejecuta 8 etapas interactivas:
1. Estado inicial del ledger
2. Org1 registra un nuevo lote (PROD-010)
3. Org1 actualiza la cantidad (merma en bodega)
4. Org1 transfiere el lote a Org2
5. Org2 confirma la recepción
6. Org2 marca el producto como entregado
7. Historial completo de PROD-010
8. Estado final de todos los productos

---

## CouchDB — World State visual

Mientras la red está activa, puedes ver el World State en el navegador:

```
http://localhost:5984/_utils
```

La base de datos del canal se llama `supplychannel_` (con guión bajo al final).

---

## Solución de problemas comunes

### Error: `No connection established` / TLS certificate error

Verifica que la ruta en `api/app.js` apunte a `~/proyecto-ebusiness/fabric-samples/test-network`. El archivo ya está configurado correctamente.

### Error: `failed to marshal response: invalid UTF-8`

El `supplycc.go` de este repositorio ya no tiene caracteres especiales (tildes, ñ, °). Si editaste el archivo manualmente, ejecuta:

```bash
grep -Pn "[^\x00-\x7F]" chaincode-go/chaincode/supplycc.go
```

Si muestra líneas, reemplaza el archivo con la versión del repositorio.

### Error: `go.mod` con versiones en conflicto (v1 y v2 simultáneas)

El `go.mod` de este repositorio usa exclusivamente `fabric-contract-api-go v1.2.1`, que es la versión compatible. Si el error persiste:

```bash
cd chaincode-go
rm -f go.sum
rm -rf vendor
go mod download
go mod tidy
go mod vendor
```

### Error: `ABORTED: failed to endorse transaction`

Verifica que el chaincode está desplegado correctamente:

```bash
export PATH=$HOME/proyecto-ebusiness/fabric-samples/bin:$PATH
export FABRIC_CFG_PATH=$HOME/proyecto-ebusiness/fabric-samples/config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_ADDRESS="localhost:7051"
export CORE_PEER_TLS_ROOTCERT_FILE=$HOME/proyecto-ebusiness/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
export CORE_PEER_MSPCONFIGPATH=$HOME/proyecto-ebusiness/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp

peer lifecycle chaincode querycommitted -C supplychannel
```

Debe mostrar `supplycc` con `Version: 1.0, Sequence: 1`.

### Reiniciar todo desde cero

```bash
cd ~/proyecto-ebusiness/fabric-samples/test-network
./network.sh down
./network.sh up createChannel -c supplychannel -ca -s couchdb
cd ~/proyecto-ebusiness
./scripts/deploy-supplycc.sh
```

---

## Smart Contract — Funciones disponibles

| Función | Descripción | Estado resultante |
|---------|-------------|-------------------|
| `InitLedger` | Inicializa 3 productos de ejemplo | REGISTRADO |
| `RegistrarProducto` | Registra un nuevo lote | REGISTRADO |
| `TransferirProducto` | Transfiere entre organizaciones | EN_TRANSITO |
| `ConfirmarRecepcion` | Confirma que Org2 recibió el lote | RECIBIDO |
| `ActualizarCantidad` | Actualiza cantidad (solo el propietario) | sin cambio |
| `MarcarEntregado` | Entrega al cliente final | ENTREGADO |
| `LeerProducto` | Consulta estado actual (World State) | — |
| `ObtenerHistorial` | Historial inmutable del ledger | — |
| `ConsultarTodos` | Rich Query — todos los productos | — |
| `ConsultarPorEstado` | Rich Query — filtrar por estado | — |

---

## Ciclo de vida de un producto

```
[REGISTRADO] --TransferirProducto--> [EN_TRANSITO] --ConfirmarRecepcion--> [RECIBIDO] --MarcarEntregado--> [ENTREGADO]
```

---

## Tecnologías

- **Hyperledger Fabric 2.5.4** — red permisionada con 2 organizaciones
- **Go 1.21** — Smart Contract (chaincode)
- **fabric-contract-api-go v1.2.1** — SDK de chaincode
- **Node.js + Express** — API REST
- **@hyperledger/fabric-gateway v1.4** — conexión a la red desde Node.js
- **CouchDB** — base de datos de estado (World State) con soporte Rich Query
- **Docker / Docker Compose** — contenedores de la red
