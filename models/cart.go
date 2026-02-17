// @name: Giuliana Moreta
//@date: 17/02/2026
//@Materia: Programación orientada a objetos
//@Curso: 3er Semestre
//@Carrera: Ingeniería en Software

// models/cart.go
// Clases CartItem y Cart con encapsulación completa
package models

import (
	"errors"
	"fmt"
)

// CLASE CartItem — campos privados

type CartItem struct {
	productID   string
	productName string
	price       float64
	quantity    int
	imageURL    string
}

// Constructor de CartItem
func NewCartItem(productID, productName string, price float64, quantity int, imageURL string) (*CartItem, error) {
	if productID == "" {
		return nil, errors.New("el ID del producto es obligatorio")
	}
	if price <= 0 {
		return nil, errors.New("el precio debe ser mayor a cero")
	}
	if quantity <= 0 {
		return nil, errors.New("la cantidad debe ser mayor a cero")
	}
	return &CartItem{
		productID:   productID,
		productName: productName,
		price:       price,
		quantity:    quantity,
		imageURL:    imageURL,
	}, nil
}

// GETTERS de CartItem
func (ci *CartItem) GetProductID() string   { return ci.productID }
func (ci *CartItem) GetProductName() string { return ci.productName }
func (ci *CartItem) GetPrice() float64      { return ci.price }
func (ci *CartItem) GetQuantity() int       { return ci.quantity }
func (ci *CartItem) GetImageURL() string    { return ci.imageURL }

// SETTER de CartItem — solo quantity tiene setter (lo demás no cambia)
func (ci *CartItem) SetQuantity(qty int) error {
	if qty <= 0 {
		return errors.New("la cantidad debe ser mayor a cero")
	}
	ci.quantity = qty
	return nil
}

// MÉTODO DE NEGOCIO
func (ci *CartItem) Subtotal() float64 {
	return ci.price * float64(ci.quantity)
}

// MarshalJSON para serializar campos privados
func (ci *CartItem) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(
		`{"product_id":%q,"product_name":%q,"price":%.2f,"quantity":%d,"image_url":%q}`,
		ci.productID, ci.productName, ci.price, ci.quantity, ci.imageURL,
	)), nil
}

// ============================================================
// CLASE Cart — campos privados
// ============================================================

type Cart struct {
	items    []CartItem
	discount float64
}

// Constructor de Cart
func NewCart() *Cart {
	return &Cart{
		items:    []CartItem{},
		discount: 0,
	}
}

// GETTERS de Cart
func (c *Cart) GetItems() []CartItem   { return c.items }
func (c *Cart) GetDiscount() float64   { return c.discount }

// SETTER de Cart — el descuento tiene validación
func (c *Cart) SetDiscount(discount float64) error {
	if discount < 0 {
		return errors.New("el descuento no puede ser negativo")
	}
	if discount > c.Subtotal() {
		return errors.New("el descuento no puede ser mayor al subtotal")
	}
	c.discount = discount
	return nil
}

// ============================================================
// MÉTODOS DE NEGOCIO de Cart
// ============================================================

// AddItem agrega un producto al carrito con validación de stock
func (c *Cart) AddItem(product *Product, qty int) error {
	if qty <= 0 {
		return errors.New("la cantidad debe ser mayor a cero")
	}
	if !product.IsAvailableQty(qty) {
		return fmt.Errorf("stock insuficiente para '%s'", product.GetName())
	}

	// Buscar si el producto ya está en el carrito
	for i, item := range c.items {
		if item.productID == product.GetID() {
			newQty := item.quantity + qty
			if !product.IsAvailableQty(newQty) {
				return errors.New("la cantidad supera el stock disponible")
			}
			// Usar el setter con validación
			if err := c.items[i].SetQuantity(newQty); err != nil {
				return err
			}
			return nil
		}
	}

	// Si no existe, crear nuevo CartItem
	item, err := NewCartItem(
		product.GetID(),
		product.GetName(),
		product.GetPrice(),
		qty,
		product.GetImageURL(),
	)
	if err != nil {
		return err
	}
	c.items = append(c.items, *item)
	return nil
}

// RemoveItem elimina un producto del carrito por ID
func (c *Cart) RemoveItem(productID string) error {
	for i, item := range c.items {
		if item.productID == productID {
			c.items = append(c.items[:i], c.items[i+1:]...)
			return nil
		}
	}
	return errors.New("producto no encontrado en el carrito")
}

// Subtotal calcula el total sin descuento
func (c *Cart) Subtotal() float64 {
	total := 0.0
	for _, item := range c.items {
		total += item.Subtotal()
	}
	return total
}

// Total calcula el total aplicando descuento
func (c *Cart) Total() float64 {
	total := c.Subtotal() - c.discount
	if total < 0 {
		return 0
	}
	return total
}

// ItemCount cuenta el total de unidades en el carrito
func (c *Cart) ItemCount() int {
	count := 0
	for _, item := range c.items {
		count += item.quantity
	}
	return count
}

// Clear vacía el carrito
func (c *Cart) Clear() {
	c.items    = []CartItem{}
	c.discount = 0
}

// IsEmpty verifica si el carrito está vacío
func (c *Cart) IsEmpty() bool {
	return len(c.items) == 0
}

// MarshalJSON para serializar campos privados
func (c *Cart) MarshalJSON() ([]byte, error) {
	itemsJSON := "["
	for i, item := range c.items {
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

	return []byte(fmt.Sprintf(
		`{"items":%s,"discount":%.2f,"subtotal":%.2f,"total":%.2f,"item_count":%d}`,
		itemsJSON, c.discount, c.Subtotal(), c.Total(), c.ItemCount(),
	)), nil
}