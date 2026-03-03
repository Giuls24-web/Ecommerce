// @name: Giuliana Moreta
// @date: 17/02/2026
// @Materia: Programación orientada a objetos
// @Curso: 3er Semestre
// @Carrera: Ingeniería en Software
package main

import (
	"ecommerce/handlers"
	"ecommerce/store"
	"log"
	"net/http"
	"os"
)

func main() {
	s := store.NewStore()
	store.SeedProducts(s)

	productHandler := handlers.NewProductHandler(s)
	cartHandler := handlers.NewCartHandler(s)
	orderHandler := handlers.NewOrderHandler(s)
	inventoryHandler := handlers.NewInventoryHandler(s)

	// Frontend estático
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	// ── CATÁLOGO PÚBLICO ─────────────────────────────────────
	// GET /api/products                → todos los productos
	// GET /api/products?category=rosa  → filtrar por categoría
	// GET /api/products/{id}           → un producto
	// GET /api/products/search?q=texto → búsqueda
	http.HandleFunc("/api/products/search", inventoryHandler.SearchProducts)
	http.HandleFunc("/api/products", productHandler.GetAll)
	http.HandleFunc("/api/products/", productHandler.GetByID)

	// ── CARRITO ──────────────────────────────────────────────
	http.HandleFunc("/api/cart", cartHandler.GetCart)
	http.HandleFunc("/api/cart/add", cartHandler.AddItem)
	http.HandleFunc("/api/cart/remove", cartHandler.RemoveItem)
	http.HandleFunc("/api/cart/clear", cartHandler.ClearCart)

	// ── ÓRDENES ──────────────────────────────────────────────
	// POST /api/orders               → crear orden
	// GET  /api/orders/list          → listar todas
	// GET  /api/orders/{id}          → ver una orden
	// PUT  /api/orders/{id}/status   → avanzar estado
	// PUT  /api/orders/{id}/cancel   → cancelar
	http.HandleFunc("/api/orders", orderHandler.CreateOrder)
	http.HandleFunc("/api/orders/list", orderHandler.ListOrders)
	http.HandleFunc("/api/orders/", orderHandler.HandleByID)

	// ── INVENTARIO (admin) ────────────────────────────────────
	// GET  /api/inventory            → ver todo el inventario
	// POST /api/inventory            → crear producto
	// PUT  /api/inventory/{id}       → editar producto
	// DELETE /api/inventory/{id}     → eliminar producto
	// PUT  /api/inventory/{id}/stock → actualizar solo el stock
	http.HandleFunc("/api/inventory", inventoryHandler.HandleInventory)
	http.HandleFunc("/api/inventory/", inventoryHandler.HandleByID)

	// Puerto dinámico para Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("🌸 FloriLuz iniciado en http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
