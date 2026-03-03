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
	"strings"
)

type OrderHandler struct {
	store *store.Store
}

func NewOrderHandler(s *store.Store) *OrderHandler {
	return &OrderHandler{store: s}
}

// CreateOrder — POST /api/orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if corsHeaders(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
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
	customer, err := models.NewCustomer(input.Name, input.Email, input.Phone, input.Address, input.City)
	if err != nil {
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
	if corsHeaders(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	respondJSON(w, h.store.GetAllOrders(), http.StatusOK)
}

// HandleByID — router para /api/orders/{id}, /api/orders/{id}/status, /api/orders/{id}/cancel
func (h *OrderHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
	if corsHeaders(w, r) {
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/orders/")

	switch {
	case strings.HasSuffix(path, "/status"):
		id := strings.TrimSuffix(path, "/status")
		h.advanceStatus(w, r, id)
	case strings.HasSuffix(path, "/cancel"):
		id := strings.TrimSuffix(path, "/cancel")
		h.cancelOrder(w, r, id)
	default:
		h.getOrder(w, r, path)
	}
}

func (h *OrderHandler) getOrder(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	order, err := h.store.GetOrder(id)
	if err != nil {
		respondError(w, err.Error(), http.StatusNotFound)
		return
	}
	respondJSON(w, order, http.StatusOK)
}

func (h *OrderHandler) advanceStatus(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPut {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	order, err := h.store.AdvanceOrderStatus(id)
	if err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, order, http.StatusOK)
}

func (h *OrderHandler) cancelOrder(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPut {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	order, err := h.store.CancelOrder(id)
	if err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, order, http.StatusOK)
}
