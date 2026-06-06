const express = require('express');
const { connect, hash, signers } = require('@hyperledger/fabric-gateway');
const grpc = require('@grpc/grpc-js');
const crypto = require('crypto');
const fs = require('fs').promises;
const path = require('path');

const app = express();
app.use(express.json());

const mspId = process.env.mspId || 'Org1MSP';
const port = process.env.PORT || 3000;

const rootPath = '/home/dani/fabric-work/fabric-samples/test-network';
const orgName = mspId === 'Org1MSP' ? 'org1.example.com' : 'org2.example.com';
const mspPath = `${rootPath}/organizations/peerOrganizations/${orgName}/users/Admin@${orgName}/msp`;

async function getGatewayConnection() {
    const tlsCertPath = `${rootPath}/organizations/peerOrganizations/${orgName}/peers/peer0.${orgName}/tls/ca.crt`;
    const tlsCert = await fs.readFile(tlsCertPath);
    const credentials = grpc.credentials.createSsl(tlsCert);
    const client = new grpc.Client(mspId === 'Org1MSP' ? 'localhost:7051' : 'localhost:9051', credentials);

    // Buscar la llave privada
    const keyPath = path.join(mspPath, 'keystore');
    const keyFiles = await fs.readdir(keyPath);
    const keyBuffer = await fs.readFile(path.join(keyPath, keyFiles[0]));
    const privateKey = crypto.createPrivateKey(keyBuffer);
    const signer = signers.newPrivateKeySigner(privateKey);

    // Buscar el certificado (BUSCA CUALQUIER .PEM)
    const certDir = path.join(mspPath, 'signcerts');
    const certFiles = await fs.readdir(certDir);
    const certFile = certFiles.find(f => f.endsWith('.pem'));
    const certBuffer = await fs.readFile(path.join(certDir, certFile));

    return connect({ client, identity: { mspId, credentials: certBuffer }, signer, hash: hash.sha256 });
}

app.get('/health', (req, res) => res.json({ status: 'OK', mspId }));

app.get('/productos', async (req, res) => {
    try {
        const gateway = await getGatewayConnection();
        const network = gateway.getNetwork('supplychannel');
        const contract = network.getContract('supplycc');
        const result = await contract.evaluateTransaction('ConsultarTodos');
        res.json(JSON.parse(new TextDecoder().decode(result)));
        gateway.close();
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

app.get('/productos/:id', async (req, res) => {
    try {
        const gateway = await getGatewayConnection();
        const network = gateway.getNetwork('supplychannel');
        const contract = network.getContract('supplycc');
        const result = await contract.evaluateTransaction('LeerProducto', req.params.id);
        res.json(JSON.parse(new TextDecoder().decode(result)));
        gateway.close();
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

app.get('/productos/:id/historial', async (req, res) => {
    try {
        const gateway = await getGatewayConnection();
        const network = gateway.getNetwork('supplychannel');
        const contract = network.getContract('supplycc');
        const result = await contract.evaluateTransaction('ObtenerHistorial', req.params.id);
        res.json(JSON.parse(new TextDecoder().decode(result)));
        gateway.close();
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

app.post('/productos', async (req, res) => {
    try {
        const { id, nombre, cantidad, unidad, origen, destino, temperatura } = req.body;
        const gateway = await getGatewayConnection();
        const network = gateway.getNetwork('supplychannel');
        const contract = network.getContract('supplycc');
        await contract.submitTransaction('RegistrarProducto', id, nombre, cantidad.toString(), unidad, origen, destino, temperatura);
        res.json({ message: `Producto ${id} registrado` });
        gateway.close();
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

app.put('/productos/:id/transferir', async (req, res) => {
    try {
        const { nuevoPropietario, detalle } = req.body;
        const gateway = await getGatewayConnection();
        const network = gateway.getNetwork('supplychannel');
        const contract = network.getContract('supplycc');
        await contract.submitTransaction('TransferirProducto', req.params.id, nuevoPropietario, detalle);
        res.json({ message: `Producto ${req.params.id} transferido` });
        gateway.close();
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
});

app.listen(port, () => console.log(`🚀 API lista en http://localhost:${port}`));
