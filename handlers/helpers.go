// @name: Giuliana Moreta
// @date: 17/02/2026
// @Materia: Programación orientada a objetos
// @Curso: 3er Semestre
// @Carrera: Ingeniería en Software


package handlers

import (
	"ecommerce/models"
	"encoding/json"
	"net/http"
)

// APIResponse — estructura estándar para TODAS las respuestas de la API
// El frontend siempre recibe el mismo formato: success + data o error
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// respondJSON envía una respuesta exitosa en formato JSON
func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
	})
}

// respondError envía una respuesta de error en formato JSON
func respondError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   message,
	})
}

// parseJSON decodifica el cuerpo JSON de una solicitud hacia cualquier struct
// Se usa para leer los datos que manda el frontend (product_id, quantity, etc.)
func parseJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// strToCategory convierte un string de la URL al tipo Category de los modelos
// Ejemplo: "rosa" → models.CategoryRose
// Si la categoría no existe, retorna rosa por defecto
func strToCategory(s string) models.Category {
	switch s {
	case "rosa":
		return models.CategoryRose
	case "girasol":
		return models.CategorySunflower
	case "loto":
		return models.CategoryLotus
	case "margarita":
		return models.CategoryDaisy
	default:
		return models.CategoryRose
	}
}