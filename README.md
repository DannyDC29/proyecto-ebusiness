# Proyecto Final: Trazabilidad de Suministros con Hyperledger Fabric

Este proyecto implementa una red de blockchain permisionada para la gestión de trazabilidad en cadenas de suministro, incluyendo un Smart Contract (Chaincode) desarrollado en Go y una API REST construida en Node.js.

## Requisitos Previos

Antes de comenzar, asegúrese de tener instalado en su sistema (preferiblemente Ubuntu/WSL2):
* **Docker y Docker Compose** (V2 recomendado)
* **Go** (Versión 1.20 o superior)
* **Node.js** (Versión 18 o superior) y **npm**
* **Git**
* **Curl**

## Instalación y Configuración

Siga estos pasos en orden para montar el entorno:

### 1. Clonar el repositorio
```bash
git clone [https://github.com/DannyDC29/proyecto-ebusiness.git](https://github.com/DannyDC29/proyecto-ebusiness.git)
cd proyecto-ebusiness
```

### 2. Descargar Binarios e Imágenes de Fabric
Si no tiene los binarios de Hyperledger Fabric en su equipo, ejecute el script oficial desde la raíz del proyecto:

```bash
curl -sSL [https://bit.ly/2ysbOFE](https://bit.ly/2ysbOFE) | bash -s -- 2.5.4 1.5.7
```

### 3. Levantar la Red y el Canal
Navegue a la carpeta de la red de prueba e inicie el canal:

```bash
cd fabric-samples/test-network
./network.sh down
./network.sh up createChannel -c supplychannel -ca -s couchdb
```

### 4. Desplegar el Smart Contract (Chaincode)
Instale el contrato inteligente encargado de la lógica de trazabilidad:

```bash
./network.sh deployCC -ccn supplycc -ccp ../asset-transfer-basic/chaincode-go/ -ccl go
```

### 5. Configurar e Iniciar la API REST
Navegue a la carpeta de la API, instale las dependencias y ejecute el servidor:

```bash
cd api
npm install
node app.js
```

La API estará disponible en **http://localhost:3000**.

### Uso de la API (Endpoints)
Puede probar la funcionalidad usando los siguientes comandos curl en una nueva terminal:

A. Verificar Salud de la Conexión

```bash
curl http://localhost:3000/health
```

B. Registrar un Nuevo Producto (POST)

```bash
curl -X POST http://localhost:3000/productos \
  -H "Content-Type: application/json" \
  -d '{"id":"PROD-100","nombre":"Cafe Especial","cantidad":500,"unidad":"kg","origen":"Huila","destino":"Puerto","temperatura":"18C"}'
```

C. Consultar Historial de Trazabilidad (GET)

```bash
curl http://localhost:3000/productos/PROD-100/historial | python3 -m json.tool
```

### Estructura del Repositorio
* **fabric-samples/test-network:** Scripts de despliegue de la red.
* **fabric-samples/asset-transfer-basic/chaincode-go:** Código fuente del Smart Contract en Go.
* **fabric-samples/test-network/api:** Servidor Express y Gateway de conexión.
* **README.md: Documentación del proyecto.

Curso: E-Business / Hyperledger Fabric Capstone
