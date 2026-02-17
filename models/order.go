// @name: Giuliana Moreta
// @date: 17/02/2026
// @Materia: Programación orientada a objetos
// @Curso: 3er Semestre
// @Carrera: Ingeniería en Software


package models

import (
	"errors"
	"fmt"
	"time"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pendiente"
	StatusPaid      OrderStatus = "pagada"
	StatusPrepared  OrderStatus = "preparada"
	StatusShipped   OrderStatus = "enviada"
	StatusDelivered OrderStatus = "entregada"
	StatusCancelled OrderStatus = "cancelada"
)

// Order — todos los campos son privados
type Order struct {
	id        string
	customer  Customer
	items     []CartItem
	total     float64
	status    OrderStatus
	notes     string
	createdAt time.Time
	updatedAt time.Time
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewOrder(id string, customer Customer, cart *Cart) (*Order, error) {
	if id == "" {
		return nil, errors.New("el ID de la orden es obligatorio")
	}
	if err := customer.Validate(); err != nil {
		return nil, fmt.Errorf("datos de cliente inválidos: %w", err)
	}
	if cart.IsEmpty() {
		return nil, errors.New("no se puede crear una orden con el carrito vacío")
	}

	// Copiar ítems del carrito — la orden guarda su propia copia
	items := make([]CartItem, len(cart.items))
	copy(items, cart.items)

	now := time.Now()
	return &Order{
		id:        id,
		customer:  customer,
		items:     items,
		total:     cart.Total(),
		status:    StatusPending,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ============================================================
// GETTERS — solo lectura
// ============================================================

func (o *Order) GetID() string           { return o.id }
func (o *Order) GetCustomer() Customer   { return o.customer }
func (o *Order) GetItems() []CartItem    { return o.items }
func (o *Order) GetTotal() float64       { return o.total }
func (o *Order) GetStatus() OrderStatus  { return o.status }
func (o *Order) GetNotes() string        { return o.notes }
func (o *Order) GetCreatedAt() time.Time { return o.createdAt }
func (o *Order) GetUpdatedAt() time.Time { return o.updatedAt }

// ============================================================
// SETTERS con validación
// ============================================================

// SetNotes permite agregar notas a la orden (instrucciones de entrega, etc.)
func (o *Order) SetNotes(notes string) {
	o.notes = notes
	o.updatedAt = time.Now()
}

// SetStatus es privado — el estado solo cambia a través de los métodos
// de negocio AdvanceStatus y Cancel, no directamente
// (no hay setter público para status: es encapsulación intencional)

// ============================================================
// MÉTODOS DE NEGOCIO — máquina de estados
// ============================================================

// AdvanceStatus avanza al siguiente estado válido
func (o *Order) AdvanceStatus() error {
	switch o.status {
	case StatusPending:
		o.status = StatusPaid
	case StatusPaid:
		o.status = StatusPrepared
	case StatusPrepared:
		o.status = StatusShipped
	case StatusShipped:
		o.status = StatusDelivered
	case StatusDelivered:
		return errors.New("la orden ya fue entregada, no puede avanzar")
	case StatusCancelled:
		return errors.New("la orden cancelada no puede cambiar de estado")
	default:
		return errors.New("estado desconocido")
	}
	o.updatedAt = time.Now()
	return nil
}

// Cancel cancela la orden si aún es posible
func (o *Order) Cancel() error {
	if o.status == StatusShipped || o.status == StatusDelivered {
		return errors.New("no se puede cancelar: la orden ya fue enviada o entregada")
	}
	if o.status == StatusCancelled {
		return errors.New("la orden ya está cancelada")
	}
	o.status    = StatusCancelled
	o.updatedAt = time.Now()
	return nil
}

// IsCancellable informa si la orden puede cancelarse
func (o *Order) IsCancellable() bool {
	return o.status != StatusShipped &&
		o.status != StatusDelivered &&
		o.status != StatusCancelled
}

// IsPending verifica si la orden está pendiente de pago
func (o *Order) IsPending() bool {
	return o.status == StatusPending
}

// Summary retorna un resumen en texto
func (o *Order) Summary() string {
	return fmt.Sprintf("Orden #%s | %s | $%.2f | Estado: %s",
		o.id, o.customer.GetName(), o.total, o.status)
}

// ============================================================
// MarshalJSON para serializar campos privados
// ============================================================

func (o *Order) MarshalJSON() ([]byte, error) {
	itemsJSON := "["
	for i, item := range o.items {
		b, err := item.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if i > 0 {
			itemsJSON += ","
		}
		itemsJSON += string(b)
	}
	itemsJSON += "]"

	customerJSON, err := o.customer.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf(
		`{"id":%q,"customer":%s,"items":%s,"total":%.2f,"status":%q,"notes":%q,"created_at":%q,"updated_at":%q}`,
		o.id, string(customerJSON), itemsJSON, o.total,
		string(o.status), o.notes,
		o.createdAt.Format(time.RFC3339),
		o.updatedAt.Format(time.RFC3339),
	)), nil
}