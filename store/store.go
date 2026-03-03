// store/store.go — Base de datos en memoria con JSON dinámico
package store

import (
	"ecommerce/models"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Store struct {
	mu       sync.Mutex
	products map[string]*models.Product
	cart     *models.Cart
	orders   map[string]*models.Order
	orderSeq int
	prodSeq  int
}

func NewStore() *Store {
	return &Store{
		products: make(map[string]*models.Product),
		cart:     models.NewCart(),
		orders:   make(map[string]*models.Order),
		orderSeq: 1,
		prodSeq:  7,
	}
}

// ── PRODUCTOS ─────────────────────────────────────────────────────────────────

func (s *Store) AddProduct(p *models.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.products[p.GetID()] = p
	return nil
}

// CreateProduct genera ID automático y crea el producto
func (s *Store) CreateProduct(name, description string, price float64, stock int, category models.Category, imageURL string) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := fmt.Sprintf("lamp-%03d", s.prodSeq)
	s.prodSeq++
	p, err := models.NewProduct(id, name, description, price, stock, category, imageURL)
	if err != nil {
		return nil, err
	}
	s.products[id] = p
	return p, nil
}

// UpdateProduct edita solo los campos que vengan no-vacíos
func (s *Store) UpdateProduct(id, name, description string, price float64, stock int, category models.Category, imageURL string) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.products[id]
	if !ok {
		return nil, fmt.Errorf("producto '%s' no encontrado", id)
	}
	if name != "" {
		if err := p.SetName(name); err != nil {
			return nil, err
		}
	}
	if description != "" {
		p.SetDescription(description)
	}
	if price > 0 {
		if err := p.SetPrice(price); err != nil {
			return nil, err
		}
	}
	if stock >= 0 {
		if err := p.SetStock(stock); err != nil {
			return nil, err
		}
	}
	if category != "" {
		if err := p.SetCategory(category); err != nil {
			return nil, err
		}
	}
	if imageURL != "" {
		p.SetImageURL(imageURL)
	}
	return p, nil
}

// DeleteProduct elimina un producto por ID
func (s *Store) DeleteProduct(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.products[id]; !ok {
		return fmt.Errorf("producto '%s' no encontrado", id)
	}
	delete(s.products, id)
	return nil
}

func (s *Store) GetProduct(id string) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.products[id]
	if !ok {
		return nil, fmt.Errorf("producto '%s' no encontrado", id)
	}
	return p, nil
}

func (s *Store) GetAllProducts() []*models.Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*models.Product, 0, len(s.products))
	for _, p := range s.products {
		out = append(out, p)
	}
	return out
}

func (s *Store) GetProductsByCategory(cat models.Category) []*models.Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []*models.Product
	for _, p := range s.products {
		if p.GetCategory() == cat {
			out = append(out, p)
		}
	}
	return out
}

// SearchProducts busca por nombre, descripción o categoría
func (s *Store) SearchProducts(q string) []*models.Product {
	s.mu.Lock()
	defer s.mu.Unlock()
	ql := strings.ToLower(q)
	var out []*models.Product
	for _, p := range s.products {
		if strings.Contains(strings.ToLower(p.GetName()), ql) ||
			strings.Contains(strings.ToLower(p.GetDescription()), ql) ||
			strings.Contains(strings.ToLower(string(p.GetCategory())), ql) {
			out = append(out, p)
		}
	}
	return out
}

// UpdateStock actualiza solo el stock de un producto
func (s *Store) UpdateStock(id string, qty int) (*models.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.products[id]
	if !ok {
		return nil, fmt.Errorf("producto '%s' no encontrado", id)
	}
	if err := p.SetStock(qty); err != nil {
		return nil, err
	}
	return p, nil
}

// ── CARRITO ───────────────────────────────────────────────────────────────────

func (s *Store) GetCart() *models.Cart {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cart
}

func (s *Store) AddToCart(productID string, qty int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.products[productID]
	if !ok {
		return fmt.Errorf("producto '%s' no existe", productID)
	}
	return s.cart.AddItem(p, qty)
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

// ── ÓRDENES ───────────────────────────────────────────────────────────────────

func (s *Store) CreateOrder(customer models.Customer) (*models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cart.IsEmpty() {
		return nil, errors.New("el carrito está vacío")
	}
	id := fmt.Sprintf("ORD-%04d", s.orderSeq)
	s.orderSeq++
	order, err := models.NewOrder(id, customer, s.cart)
	if err != nil {
		return nil, err
	}
	for _, item := range order.GetItems() {
		p, ok := s.products[item.GetProductID()]
		if !ok {
			return nil, fmt.Errorf("producto '%s' no encontrado", item.GetProductID())
		}
		if err := p.DecreaseStock(item.GetQuantity()); err != nil {
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
	o, ok := s.orders[id]
	if !ok {
		return nil, fmt.Errorf("orden '%s' no encontrada", id)
	}
	return o, nil
}

func (s *Store) GetAllOrders() []*models.Order {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*models.Order, 0, len(s.orders))
	for _, o := range s.orders {
		out = append(out, o)
	}
	return out
}

// AdvanceOrderStatus avanza la máquina de estados de una orden
func (s *Store) AdvanceOrderStatus(id string) (*models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	if !ok {
		return nil, fmt.Errorf("orden '%s' no encontrada", id)
	}
	if err := o.AdvanceStatus(); err != nil {
		return nil, err
	}
	return o, nil
}

// CancelOrder cancela una orden
func (s *Store) CancelOrder(id string) (*models.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	if !ok {
		return nil, fmt.Errorf("orden '%s' no encontrada", id)
	}
	if err := o.Cancel(); err != nil {
		return nil, err
	}
	return o, nil
}

// ── SEED ──────────────────────────────────────────────────────────────────────

func SeedProducts(s *Store) {
	items := []struct {
		id, name, desc, img string
		price               float64
		stock               int
		cat                 models.Category
	}{
		{"lamp-001", "Lámpara Rosa Romántica", "Elegante lámpara con pétalos de rosa en porcelana fría, luz cálida LED.", "https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=400", 49.99, 15, models.CategoryRose},
		{"lamp-002", "Lámpara Girasol Primaveral", "Lámpara de pie con pétalos de resina dorada inspirada en el girasol.", "https://images.unsplash.com/photo-1513475382585-d06e58bcb0e0?w=400", 89.99, 8, models.CategorySunflower},
		{"lamp-003", "Lámpara Loto Zen", "Lámpara de ambiente con flor de loto. Emite luz suave y relajante.", "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400", 65.00, 12, models.CategoryLotus},
		{"lamp-004", "Lámpara Margarita Alegre", "Lámpara infantil multicolor con forma de margarita. Segura para niños.", "https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=400", 35.50, 20, models.CategoryDaisy},
		{"lamp-005", "Lámpara Rosa Vintage", "Lámpara colgante estilo vintage con motivos de rosas antiguas.", "https://images.unsplash.com/photo-1524484485831-a92ffc0de03f?w=400", 75.00, 6, models.CategoryRose},
		{"lamp-006", "Lámpara Girasol Mini", "Mini lámpara de escritorio con diseño de girasol.", "https://images.unsplash.com/photo-1491553895911-0055eca6402d?w=400", 28.99, 25, models.CategorySunflower},
	}
	for _, d := range items {
		p, err := models.NewProduct(d.id, d.name, d.desc, d.price, d.stock, d.cat, d.img)
		if err != nil {
			continue
		}
		s.AddProduct(p)
	}
}
