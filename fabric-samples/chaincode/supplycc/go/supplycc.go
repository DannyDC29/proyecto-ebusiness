package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// =====================================================================
//  supplycc — Chaincode de Supply Chain Management
//  Laboratorio Blockchain Empresarial | Jaider Reyes Herazo
// =====================================================================

// Producto representa un lote en la cadena de suministro
type Producto struct {
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	Cantidad    int    `json:"cantidad"`
	Unidad      string `json:"unidad"`
	Propietario string `json:"propietario"`
	Estado      string `json:"estado"`   // REGISTRADO | EN_TRANSITO | RECIBIDO | ENTREGADO
	Origen      string `json:"origen"`
	Destino     string `json:"destino"`
	Timestamp   string `json:"timestamp"`
	TxID        string `json:"txId"`
}

// EventoHistorial registra cada cambio de estado
type EventoHistorial struct {
	TxID      string `json:"txId"`
	Timestamp string `json:"timestamp"`
	Estado    string `json:"estado"`
	Actor     string `json:"actor"`
	Detalle   string `json:"detalle"`
}

// SupplyContract implementa la interfaz ContractInterface
type SupplyContract struct {
	contractapi.Contract
}

// =====================================================================
//  FUNCIÓN 1 — InitLedger
//  Inicializa el ledger con datos de ejemplo
// =====================================================================
func (sc *SupplyContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	productos := []Producto{
		{
			ID:          "PROD-001",
			Nombre:      "Cacao Orgánico",
			Cantidad:    500,
			Unidad:      "kg",
			Propietario: "Org1MSP",
			Estado:      "REGISTRADO",
			Origen:      "Tumaco, Colombia",
			Destino:     "Medellín Hub",
			Timestamp:   time.Now().Format(time.RFC3339),
			TxID:        ctx.GetStub().GetTxID(),
		},
		{
			ID:          "PROD-002",
			Nombre:      "Café Especial",
			Cantidad:    200,
			Unidad:      "kg",
			Propietario: "Org1MSP",
			Estado:      "REGISTRADO",
			Origen:      "Huila, Colombia",
			Destino:     "Bogotá DC",
			Timestamp:   time.Now().Format(time.RFC3339),
			TxID:        ctx.GetStub().GetTxID(),
		},
	}

	for _, p := range productos {
		pJSON, err := json.Marshal(p)
		if err != nil {
			return fmt.Errorf("error serializando producto %s: %v", p.ID, err)
		}
		if err := ctx.GetStub().PutState(p.ID, pJSON); err != nil {
			return fmt.Errorf("error guardando producto %s: %v", p.ID, err)
		}
	}
	return nil
}

// =====================================================================
//  FUNCIÓN 2 — RegistrarProducto
//  El Proveedor (Org1) registra un nuevo lote en la red
// =====================================================================
func (sc *SupplyContract) RegistrarProducto(
	ctx contractapi.TransactionContextInterface,
	id, nombre string,
	cantidad int,
	unidad, origen, destino string,
) error {
	// Verificar que no exista
	existe, err := sc.existeProducto(ctx, id)
	if err != nil {
		return err
	}
	if existe {
		return fmt.Errorf("el producto %s ya existe en el ledger", id)
	}

	// Obtener identidad del invocador
	clientID, _ := ctx.GetClientIdentity().GetMSPID()

	producto := Producto{
		ID:          id,
		Nombre:      nombre,
		Cantidad:    cantidad,
		Unidad:      unidad,
		Propietario: clientID,
		Estado:      "REGISTRADO",
		Origen:      origen,
		Destino:     destino,
		Timestamp:   time.Now().Format(time.RFC3339),
		TxID:        ctx.GetStub().GetTxID(),
	}

	pJSON, err := json.Marshal(producto)
	if err != nil {
		return err
	}

	// Emitir evento blockchain
	ctx.GetStub().SetEvent("ProductoRegistrado", pJSON)

	return ctx.GetStub().PutState(id, pJSON)
}

// =====================================================================
//  FUNCIÓN 3 — TransferirProducto
//  Cambia el propietario y estado cuando hay transferencia
// =====================================================================
func (sc *SupplyContract) TransferirProducto(
	ctx contractapi.TransactionContextInterface,
	id, nuevoPropietario, detalle string,
) error {
	producto, err := sc.LeerProducto(ctx, id)
	if err != nil {
		return err
	}

	// Registrar evento en historial (composite key)
	evento := EventoHistorial{
		TxID:      ctx.GetStub().GetTxID(),
		Timestamp: time.Now().Format(time.RFC3339),
		Estado:    "EN_TRANSITO",
		Actor:     producto.Propietario,
		Detalle:   detalle,
	}
	sc.guardarHistorial(ctx, id, evento)

	// Actualizar estado del producto
	producto.Propietario = nuevoPropietario
	producto.Estado = "EN_TRANSITO"
	producto.TxID = ctx.GetStub().GetTxID()
	producto.Timestamp = time.Now().Format(time.RFC3339)

	pJSON, err := json.Marshal(producto)
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("ProductoTransferido", pJSON)
	return ctx.GetStub().PutState(id, pJSON)
}

// =====================================================================
//  FUNCIÓN 4 — ConfirmarRecepcion
//  El Distribuidor (Org2) confirma que recibió el lote
// =====================================================================
func (sc *SupplyContract) ConfirmarRecepcion(
	ctx contractapi.TransactionContextInterface,
	id, detalle string,
) error {
	producto, err := sc.LeerProducto(ctx, id)
	if err != nil {
		return err
	}

	if producto.Estado != "EN_TRANSITO" {
		return fmt.Errorf("el producto %s no está EN_TRANSITO (estado actual: %s)", id, producto.Estado)
	}

	clientID, _ := ctx.GetClientIdentity().GetMSPID()

	evento := EventoHistorial{
		TxID:      ctx.GetStub().GetTxID(),
		Timestamp: time.Now().Format(time.RFC3339),
		Estado:    "RECIBIDO",
		Actor:     clientID,
		Detalle:   detalle,
	}
	sc.guardarHistorial(ctx, id, evento)

	producto.Estado = "RECIBIDO"
	producto.TxID = ctx.GetStub().GetTxID()
	producto.Timestamp = time.Now().Format(time.RFC3339)

	pJSON, _ := json.Marshal(producto)
	ctx.GetStub().SetEvent("RecepcionConfirmada", pJSON)
	return ctx.GetStub().PutState(id, pJSON)
}

// =====================================================================
//  FUNCIÓN 5 — LeerProducto
//  Consulta el estado actual de un producto (World State)
// =====================================================================
func (sc *SupplyContract) LeerProducto(
	ctx contractapi.TransactionContextInterface,
	id string,
) (*Producto, error) {
	pJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("error leyendo producto %s: %v", id, err)
	}
	if pJSON == nil {
		return nil, fmt.Errorf("producto %s no encontrado en el ledger", id)
	}

	var producto Producto
	if err := json.Unmarshal(pJSON, &producto); err != nil {
		return nil, err
	}
	return &producto, nil
}

// =====================================================================
//  FUNCIÓN 6 — ObtenerHistorial
//  Retorna el historial completo de un producto (Ledger history)
// =====================================================================
func (sc *SupplyContract) ObtenerHistorial(
	ctx contractapi.TransactionContextInterface,
	id string,
) ([]*EventoHistorial, error) {
	iterator, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var historial []*EventoHistorial
	for iterator.HasNext() {
		respuesta, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		var producto Producto
		if err := json.Unmarshal(respuesta.Value, &producto); err != nil {
			continue
		}

		evento := &EventoHistorial{
			TxID:      respuesta.TxId,
			Timestamp: time.Unix(respuesta.Timestamp.Seconds, 0).Format(time.RFC3339),
			Estado:    producto.Estado,
			Actor:     producto.Propietario,
			Detalle:   "Registro histórico del ledger",
		}
		historial = append(historial, evento)
	}
	return historial, nil
}

// =====================================================================
//  FUNCIÓN 7 — ConsultarTodos
//  Rich query: retorna todos los productos (requiere CouchDB)
// =====================================================================
func (sc *SupplyContract) ConsultarTodos(
	ctx contractapi.TransactionContextInterface,
) ([]*Producto, error) {
	query := `{"selector":{"id":{"$gt":""}}}`
	iterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var productos []*Producto
	for iterator.HasNext() {
		resp, err := iterator.Next()
		if err != nil {
			return nil, err
		}
		var p Producto
		if err := json.Unmarshal(resp.Value, &p); err != nil {
			continue
		}
		productos = append(productos, &p)
	}
	return productos, nil
}

// =====================================================================
//  FUNCIÓN 8 — ConsultarPorEstado
//  Rich query: filtra por estado (REGISTRADO | EN_TRANSITO | RECIBIDO)
// =====================================================================
func (sc *SupplyContract) ConsultarPorEstado(
	ctx contractapi.TransactionContextInterface,
	estado string,
) ([]*Producto, error) {
	query := fmt.Sprintf(`{"selector":{"estado":"%s"}}`, estado)
	iterator, err := ctx.GetStub().GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var productos []*Producto
	for iterator.HasNext() {
		resp, err := iterator.Next()
		if err != nil {
			return nil, err
		}
		var p Producto
		if err := json.Unmarshal(resp.Value, &p); err != nil {
			continue
		}
		productos = append(productos, &p)
	}
	return productos, nil
}

// =====================================================================
//  HELPERS INTERNOS
// =====================================================================
func (sc *SupplyContract) existeProducto(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

func (sc *SupplyContract) guardarHistorial(
	ctx contractapi.TransactionContextInterface,
	productoID string,
	evento EventoHistorial,
) {
	key, _ := ctx.GetStub().CreateCompositeKey("historial", []string{productoID, evento.TxID})
	eJSON, _ := json.Marshal(evento)
	ctx.GetStub().PutState(key, eJSON)
}

// =====================================================================
//  MAIN
// =====================================================================
func main() {
	chaincode, err := contractapi.NewChaincode(&SupplyContract{})
	if err != nil {
		fmt.Printf("Error creando chaincode: %v\n", err)
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error iniciando chaincode: %v\n", err)
	}
}
