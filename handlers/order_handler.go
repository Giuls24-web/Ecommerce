// @name: Giuliana Moreta
// @date: 17/02/2026
// @Materia: Programación orientada a objetos
// @Curso: 3er Semestre
// @Carrera: Ingeniería en Software

// handlers/order_handler.go
package handlers

import (
	"ecommerce/models"
	"ecommerce/store"
	"encoding/json"
	"net/http"
)

type OrderHandler struct {
	store *store.Store
}

func NewOrderHandler(s *store.Store) *OrderHandler {
	return &OrderHandler{store: s}
}

// CreateOrder — POST /api/orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Struct auxiliar con campos públicos para recibir el JSON del frontend
	var input struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
		City    string `json:"city"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "Datos del cliente inválidos", http.StatusBadRequest)
		return
	}

	// Crear Customer usando el constructor con validación (setters internos)
	customer, err := models.NewCustomer(input.Name, input.Email, input.Phone, input.Address, input.City)
	if err != nil {
		// El error viene del setter que falló (nombre vacío, email inválido, etc.)
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.store.CreateOrder(*customer)
	if err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, order, http.StatusCreated)
}

// ListOrders — GET /api/orders/list
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	respondJSON(w, h.store.GetAllOrders(), http.StatusOK)
}