// @name: Giuliana Moreta
// @date: 17/02/2026
// @Materia: Programación orientada a objetos
// @Curso: 3er Semestre
// @Carrera: Ingeniería en Software

// store/store.go
// Store: base de datos en memoria
// Usa getters/setters de los modelos para acceder a los datos
package store

import (
	"ecommerce/models"
	"errors"
	"fmt"
	"sync"
)

type Store struct {
	mu       sync.Mutex
	products map[string]*models.Product
	cart     *models.Cart
	orders   map[string]*models.Order
	orderSeq int
}

func NewStore() *Store {
	return &Store{
		products: make(map[string]*models.Product),
		cart:     models.NewCart(),
		orders:   make(map[string]*models.Order),
		orderSeq: 1,
	}
}

// PRODUCTOS

func (s *Store) AddProduct(p *models.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p.GetID() == "" {
		return errors.New("el producto debe tener un ID")
	}
	s.products[p.GetID()] = p
	return nil
}

func (s *Store) GetProduct(id string) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, exists := s.products[id]
	if !exists {
		return nil, fmt.Errorf("producto '%s' no encontrado", id)
	}
	return p, nil
}

func (s *Store) GetAllProducts() []*models.Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	products := make([]*models.Product, 0, len(s.products))
	for _, p := range s.products {
		products = append(products, p)
	}
	return products
}

func (s *Store) GetProductsByCategory(cat models.Category) []*models.Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []*models.Product
	for _, p := range s.products {
		// Usamos el GETTER en lugar de acceder al campo directamente
		if p.GetCategory() == cat {
			result = append(result, p)
		}
	}
	return result
}

// ============================================================
// CARRITO
// ============================================================

func (s *Store) GetCart() *models.Cart {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cart
}

func (s *Store) AddToCart(productID string, qty int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	product, exists := s.products[productID]
	if !exists {
		return fmt.Errorf("producto '%s' no existe", productID)
	}
	return s.cart.AddItem(product, qty)
}

func (s *Store) RemoveFromCart(productID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cart.RemoveItem(productID)
}

func (s *Store) ClearCart() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cart.Clear()
}

// ÓRDENES

func (s *Store) CreateOrder(customer models.Customer) (*models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cart.IsEmpty() {
		return nil, errors.New("el carrito está vacío")
	}

	orderID := fmt.Sprintf("ORD-%04d", s.orderSeq)
	s.orderSeq++

	order, err := models.NewOrder(orderID, customer, s.cart)
	if err != nil {
		return nil, err
	}

	// Descontar stock usando el getter para obtener el ID y el método DecreaseStock
	for _, item := range order.GetItems() {
		product, exists := s.products[item.GetProductID()]
		if !exists {
			return nil, fmt.Errorf("producto '%s' no encontrado", item.GetProductID())
		}
		// DecreaseStock es un método de negocio encapsulado en Product
		if err := product.DecreaseStock(item.GetQuantity()); err != nil {
			return nil, err
		}
	}

	s.orders[order.GetID()] = order
	s.cart.Clear()
	return order, nil
}

func (s *Store) GetOrder(id string) (*models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, exists := s.orders[id]
	if !exists {
		return nil, fmt.Errorf("orden '%s' no encontrada", id)
	}
	return order, nil
}

func (s *Store) GetAllOrders() []*models.Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	orders := make([]*models.Order, 0, len(s.orders))
	for _, o := range s.orders {
		orders = append(orders, o)
	}
	return orders
}

// ============================================================
// SEED — datos de ejemplo
// ============================================================

func SeedProducts(s *Store) {
	data := []struct {
		id, name, description string
		price                 float64
		stock                 int
		category              models.Category
		imageURL              string
	}{
		{"lamp-001", "Lámpara Rosa Romántica", "Elegante lámpara de mesa con pétalos de rosa en porcelana fría, luz cálida LED.", 49.99, 15, models.CategoryRose, "https://es.bestdealplus.com/product/2280182/Lamparas-de-Mesa-de-flores-rosas-romanticas-para-dormitorio-lampara-de-escritorio-de-vidrio-rosa-soporte-de-camas-Led-moderno-accesorios-de-iluminacion-decoracion-de-boda"},
		{"lamp-002", "Lámpara Girasol Primaveral", "Lámpara de pie inspirada en el girasol, con pétalos de resina dorada.", 89.99, 8, models.CategorySunflower, "https://m.media-amazon.com/images/I/71seWZWMIvL._AC_SL1500_.jpg"},
		{"lamp-003", "Lámpara Loto Zen", "Lámpara de ambiente inspirada en la flor de loto. Emite luz suave y relajante.", 65.00, 12, models.CategoryLotus, "https://mx.pinterest.com/pin/lmpara-de-flor-de-loto--563020390936762429/"},
		{"lamp-004", "Lámpara Margarita Alegre", "Lámpara infantil con forma de margarita multicolor. Segura para niños.", 35.50, 20, models.CategoryDaisy, "https://www.ubuy.ec/es/product/4L3YESXLG-led-waterproof-floating-lotus-light-pond-light-battery-operated-lily-flower-white-light-flower-night-lamp-pack-of-5-2-lily-pad-8?srsltid=AfmBOoo04ldut2JrPMlh9li6tg1nV7bQ0ve5NVCBYLFhSbHvN-olFZNw"},
		{"lamp-005", "Lámpara Rosa Vintage", "Lámpara colgante estilo vintage con motivos de rosas antiguas.", 75.00, 6, models.CategoryRose, "https://m.media-amazon.com/images/I/61AV5CcHo3L._AC_UF894,1000_QL80_.jpg"},
		{"lamp-006", "Lámpara Girasol Mini", "Mini lámpara de escritorio con diseño de girasol. Perfecta para tu espacio de trabajo.", 28.99, 25, models.CategorySunflower, "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTUYSsDLC4OWo4BxenE4-JotSri9Sn48EwcSw&s"},
	}

	for _, d := range data {
		// NewProduct usa validación interna antes de crear
		p, err := models.NewProduct(d.id, d.name, d.description, d.price, d.stock, d.category, d.imageURL)
		if err != nil {
			continue
		}
		s.AddProduct(p)
	}
}
