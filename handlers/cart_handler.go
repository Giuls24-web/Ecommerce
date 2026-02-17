// @name: Giuliana Moreta
// @date: 17/02/2026
// @Materia: Programación orientada a objetos
// @Curso: 3er Semestre
// @Carrera: Ingeniería en Software


package handlers

import (
	"ecommerce/store"
	"net/http"
)

// CartHandler depende del Store
type CartHandler struct {
	store *store.Store
}

// NewCartHandler — constructor del handler
func NewCartHandler(s *store.Store) *CartHandler {
	return &CartHandler{store: s}
}

// GetCart responde a GET /api/cart
// Retorna el carrito actual con todos sus ítems y totales
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	cart := h.store.GetCart()

	// Usamos los GETTERS del cart para leer su estado
	// cart.items sería error — los campos son privados
	// cart.GetItems() es la forma correcta
	_ = cart.GetItems()    // getter: lista de ítems
	_ = cart.GetDiscount() // getter: descuento aplicado

	// El cart se serializa con su propio MarshalJSON()
	respondJSON(w, cart, http.StatusOK)
}

// AddItem responde a POST /api/cart/add
// Body esperado: { "product_id": "lamp-001", "quantity": 2 }
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Struct auxiliar con campos PÚBLICOS para recibir el JSON
	// (necesario porque parseJSON usa encoding/json que requiere campos públicos)
	var body struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := parseJSON(r, &body); err != nil {
		respondError(w, "Cuerpo de solicitud inválido", http.StatusBadRequest)
		return
	}

	// Validaciones en el handler antes de llamar al store
	if body.ProductID == "" {
		respondError(w, "Se requiere product_id", http.StatusBadRequest)
		return
	}
	if body.Quantity <= 0 {
		respondError(w, "La cantidad debe ser mayor a cero", http.StatusBadRequest)
		return
	}

	// El store llama a cart.AddItem() que internamente usa
	// los getters de Product (GetID, GetName, GetPrice, GetStock)
	if err := h.store.AddToCart(body.ProductID, body.Quantity); err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	cart := h.store.GetCart()

	// Ejemplo de uso de getters para loguear info del carrito
	// sin acceder a campos privados directamente
	totalItems := cart.ItemCount()  // getter de comportamiento
	total      := cart.Total()      // getter calculado
	_ = totalItems
	_ = total

	respondJSON(w, cart, http.StatusOK)
}

// RemoveItem responde a POST /api/cart/remove
// Body esperado: { "product_id": "lamp-001" }
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		ProductID string `json:"product_id"`
	}

	if err := parseJSON(r, &body); err != nil {
		respondError(w, "Cuerpo de solicitud inválido", http.StatusBadRequest)
		return
	}

	if body.ProductID == "" {
		respondError(w, "Se requiere product_id", http.StatusBadRequest)
		return
	}

	// RemoveFromCart usa cart.RemoveItem() que internamente
	// busca por productID usando el campo privado
	if err := h.store.RemoveFromCart(body.ProductID); err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	cart := h.store.GetCart()
	respondJSON(w, cart, http.StatusOK)
}

// ClearCart responde a POST /api/cart/clear
// Vacía completamente el carrito
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	h.store.ClearCart() // llama a cart.Clear() internamente

	respondJSON(w, map[string]string{
		"message": "Carrito vaciado exitosamente",
	}, http.StatusOK)
}