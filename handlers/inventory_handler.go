// handlers/inventory_handler.go — CRUD de productos + búsqueda
package handlers

import (
	"ecommerce/store"
	"encoding/json"
	"net/http"
	"strings"
)

type InventoryHandler struct {
	store *store.Store
}

func NewInventoryHandler(s *store.Store) *InventoryHandler {
	return &InventoryHandler{store: s}
}

// HandleInventory → GET /api/inventory  |  POST /api/inventory
func (h *InventoryHandler) HandleInventory(w http.ResponseWriter, r *http.Request) {
	if corsHeaders(w, r) {
		return
	}
	switch r.Method {
	case http.MethodGet:
		respondJSON(w, h.store.GetAllProducts(), http.StatusOK)
	case http.MethodPost:
		h.createProduct(w, r)
	default:
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// HandleByID → PUT /api/inventory/{id}  |  DELETE /api/inventory/{id}  |  PUT /api/inventory/{id}/stock
func (h *InventoryHandler) HandleByID(w http.ResponseWriter, r *http.Request) {
	if corsHeaders(w, r) {
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/inventory/")

	// PUT /api/inventory/{id}/stock
	if strings.HasSuffix(path, "/stock") {
		id := strings.TrimSuffix(path, "/stock")
		h.updateStock(w, r, id)
		return
	}

	id := path
	if id == "" {
		respondError(w, "ID requerido", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPut:
		h.updateProduct(w, r, id)
	case http.MethodDelete:
		if err := h.store.DeleteProduct(id); err != nil {
			respondError(w, err.Error(), http.StatusNotFound)
			return
		}
		respondJSON(w, map[string]string{"message": "Producto eliminado"}, http.StatusOK)
	default:
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// SearchProducts → GET /api/products/search?q=
func (h *InventoryHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	if corsHeaders(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query().Get("q")
	if q == "" {
		respondJSON(w, h.store.GetAllProducts(), http.StatusOK)
		return
	}
	respondJSON(w, h.store.SearchProducts(q), http.StatusOK)
}

// ── internos ─────────────────────────────────────────────────

func (h *InventoryHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		Category    string  `json:"category"`
		ImageURL    string  `json:"image_url"`
	}
	if err := parseJSON(r, &body); err != nil {
		respondError(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	if body.Name == "" {
		respondError(w, "El nombre es obligatorio", http.StatusBadRequest)
		return
	}
	if body.Price <= 0 {
		respondError(w, "El precio debe ser mayor a cero", http.StatusBadRequest)
		return
	}
	p, err := h.store.CreateProduct(body.Name, body.Description, body.Price,
		body.Stock, strToCategory(body.Category), body.ImageURL)
	if err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, p, http.StatusCreated)
}

func (h *InventoryHandler) updateProduct(w http.ResponseWriter, r *http.Request, id string) {
	var body struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		Category    string  `json:"category"`
		ImageURL    string  `json:"image_url"`
	}
	if err := parseJSON(r, &body); err != nil {
		respondError(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	p, err := h.store.UpdateProduct(id, body.Name, body.Description,
		body.Price, body.Stock, strToCategory(body.Category), body.ImageURL)
	if err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, p, http.StatusOK)
}

func (h *InventoryHandler) updateStock(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPut {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Stock int `json:"stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	p, err := h.store.UpdateStock(id, body.Stock)
	if err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, p, http.StatusOK)
}
