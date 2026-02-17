// main.go - Punto de entrada del sistema e-commerce de Lámparas Florales
// Este archivo inicializa el servidor HTTP y registra todas las rutas
package main

import (
	"log"
	"net/http"

	"ecommerce/handlers"
	"ecommerce/store"
)

func main() {
	// Inicializamos la tienda (base de datos en memoria con productos de ejemplo)
	s := store.NewStore()
	store.SeedProducts(s)

	// Creamos los handlers (controladores) pasándoles la tienda
	productHandler := handlers.NewProductHandler(s)
	cartHandler := handlers.NewCartHandler(s)
	orderHandler := handlers.NewOrderHandler(s)

	// Servir archivos estáticos del frontend (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	// Rutas de la API REST
	// Productos
	http.HandleFunc("/api/products", productHandler.GetAll)
	http.HandleFunc("/api/products/", productHandler.GetByID)

	// Carrito
	http.HandleFunc("/api/cart", cartHandler.GetCart)
	http.HandleFunc("/api/cart/add", cartHandler.AddItem)
	http.HandleFunc("/api/cart/remove", cartHandler.RemoveItem)
	http.HandleFunc("/api/cart/clear", cartHandler.ClearCart)

	// Órdenes
	http.HandleFunc("/api/orders", orderHandler.CreateOrder)
	http.HandleFunc("/api/orders/list", orderHandler.ListOrders)

	log.Println("Servidor de Lámparas Florales iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
