// @name: Giuliana Moreta
// @date: 17/02/2026
// @Materia: Programación orientada a objetos
// @Curso: 3er Semestre
// @Carrera: Ingeniería en Software


package handlers

import (
	"ecommerce/store"
	"net/http"
	"strings"
)

// ProductHandler depende del Store para acceder a los datos
// Relación de dependencia: handler → store → models
type ProductHandler struct {
	store *store.Store
}

// NewProductHandler — constructor del handler
func NewProductHandler(s *store.Store) *ProductHandler {
	return &ProductHandler{store: s}
}

// GetAll responde a GET /api/products
// Opcionalmente filtra por categoría: GET /api/products?category=rosa
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer parámetro opcional de categoría desde la URL
	category := r.URL.Query().Get("category")

	var result interface{}
	if category != "" {
		// strToCategory convierte "rosa" → models.CategoryRose (usa el tipo encapsulado)
		result = h.store.GetProductsByCategory(strToCategory(category))
	} else {
		result = h.store.GetAllProducts()
	}

	// Los productos se serializan usando su MarshalJSON(),
	// que accede a sus campos privados internamente
	respondJSON(w, result, http.StatusOK)
}

// GetByID responde a GET /api/products/{id}
// Extrae el ID de la URL y busca el producto en el store
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer el ID quitando el prefijo de la ruta
	// Ejemplo: /api/products/lamp-001 → lamp-001
	id := strings.TrimPrefix(r.URL.Path, "/api/products/")
	if id == "" {
		respondError(w, "ID de producto requerido", http.StatusBadRequest)
		return
	}

	product, err := h.store.GetProduct(id)
	if err != nil {
		// El store retorna error descriptivo si no encuentra el producto
		respondError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Accedemos al nombre usando el GETTER, no el campo directo
	// (product.name sería error de compilación porque es privado)
	_ = product.GetName() // ejemplo de uso de getter en el handler

	respondJSON(w, product, http.StatusOK)
}