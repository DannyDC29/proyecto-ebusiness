# Proyecto Final: Trazabilidad de Suministros con Hyperledger Fabric

Este proyecto implementa una red de blockchain permisionada para la gestión de trazabilidad en cadenas de suministro, incluyendo un Smart Contract (Chaincode) desarrollado en Go y una API REST construida en Node.js.

## 1. Requisitos del Sistema e Instalación

Si no tiene instaladas las herramientas necesarias en su entorno Ubuntu/WSL2, ejecute los siguientes comandos:

### A. Instalar Docker y Docker Compose
```bash
sudo apt update
sudo apt install -y docker.io docker-compose
sudo usermod -aG docker $USER
# Reinicie la sesión de terminal después de esto

```

### B. Instalar Go (Golang)

```bash
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### C. Instalar Node.js y npm

```bash
sudo apt install -y nodejs npm

```

### D. Instalar Git y Curl

```bash
sudo apt install -y git curl

```

---

## 2. Instalación del Proyecto

### Paso 1: Clonar el repositorio

```bash
git clone https://github.com/DannyDC29/proyecto-ebusiness.git
cd proyecto-ebusiness

```

### Paso 2: Descargar Binarios e Imágenes de Fabric

Esto descargará los ejecutables de Hyperledger Fabric (peer, orderer, etc.) y las imágenes de Docker necesarias:

```bash
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.5.4 1.5.7

```

---

## 3. Puesta en Marcha de la Red

### Paso 3: Levantar la Red y el Canal

Navegue a la carpeta de la red de prueba e inicie el canal con base de datos CouchDB:

```bash
cd fabric-samples/test-network
./network.sh down
./network.sh up createChannel -c supplychannel -ca -s couchdb

```

### Paso 4: Desplegar el Smart Contract (Chaincode)

Instale el contrato inteligente encargado de la lógica de trazabilidad:

```bash
./network.sh deployCC -ccn supplycc -ccp ../asset-transfer-basic/chaincode-go/ -ccl go

```

---

## 4. Configuración de la API REST

### Paso 5: Iniciar el servidor de la API

Navegue a la carpeta de la API, instale las dependencias de Node.js y ejecute el servidor:

```bash
cd api
npm install
node app.js

```

La API estará disponible en **http://localhost:3000**.

---

## 5. Guía de Pruebas (Endpoints)

Abra una **nueva terminal** y use los siguientes comandos para interactuar con la Blockchain:

### A. Verificar Salud de la Conexión

```bash
curl http://localhost:3000/health

```

### B. Registrar un Nuevo Producto (POST)

```bash
curl -X POST http://localhost:3000/productos \
  -H "Content-Type: application/json" \
  -d '{"id":"PROD-100","nombre":"Cafe Especial","cantidad":500,"unidad":"kg","origen":"Huila","destino":"Puerto","temperatura":"18C"}'

```

### C. Consultar Historial de Trazabilidad (GET)

```bash
curl http://localhost:3000/productos/PROD-100/historial | python3 -m json.tool

```

---

## Estructura del Repositorio

* **fabric-samples/test-network:** Scripts de gestión de la red.
* **fabric-samples/asset-transfer-basic/chaincode-go:** Lógica del Smart Contract en Go.
* **fabric-samples/test-network/api:** Servidor Express y Gateway de conexión.

**Curso:** E-Business / Hyperledger Fabric Capstone
