// @name: Giuliana Moreta
//@date: 17/02/2026
//@Materia: Programación orientada a objetos
//@Curso: 3er Semestre
//@Carrera: Ingeniería en Software


package models

import (
	"errors"
	"fmt"
	"time"
)

type Category string

const (
	CategoryRose      Category = "rosa"
	CategorySunflower Category = "girasol"
	CategoryLotus     Category = "loto"
	CategoryDaisy     Category = "margarita"
)

// Product — todos los campos son PRIVADOS (minúscula)
// Nadie puede hacer product.price = -99 desde afuera
// Solo se accede mediante los métodos públicos (getters/setters)
type Product struct {
	id          string
	name        string
	description string
	price       float64
	stock       int
	category    Category
	imageURL    string
	createdAt   time.Time
}

// ============================================================
// CONSTRUCTOR
// ============================================================

func NewProduct(id, name, description string, price float64, stock int, category Category, imageURL string) (*Product, error) {
	if id == "" {
		return nil, errors.New("el ID no puede estar vacío")
	}
	if name == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}
	if price <= 0 {
		return nil, errors.New("el precio debe ser mayor a cero")
	}
	if stock < 0 {
		return nil, errors.New("el stock no puede ser negativo")
	}

	return &Product{
		id:          id,
		name:        name,
		description: description,
		price:       price,
		stock:       stock,
		category:    category,
		imageURL:    imageURL,
		createdAt:   time.Now(),
	}, nil
}

// GETTERS

func (p *Product) GetID() string          { return p.id }
func (p *Product) GetName() string        { return p.name }
func (p *Product) GetDescription() string { return p.description }
func (p *Product) GetPrice() float64      { return p.price }
func (p *Product) GetStock() int          { return p.stock }
func (p *Product) GetCategory() Category  { return p.category }
func (p *Product) GetImageURL() string    { return p.imageURL }
func (p *Product) GetCreatedAt() time.Time { return p.createdAt }

// SETTERS

// SetName valida que el nombre no esté vacío antes de asignarlo
func (p *Product) SetName(name string) error {
	if name == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	p.name = name
	return nil
}

// SetDescription no requiere validación especial
func (p *Product) SetDescription(desc string) {
	p.description = desc
}

// SetPrice valida que el precio sea positivo
func (p *Product) SetPrice(price float64) error {
	if price <= 0 {
		return errors.New("el precio debe ser mayor a cero")
	}
	p.price = price
	return nil
}

// SetCategory valida que la categoría sea una de las permitidas
func (p *Product) SetCategory(cat Category) error {
	switch cat {
	case CategoryRose, CategorySunflower, CategoryLotus, CategoryDaisy:
		p.category = cat
		return nil
	default:
		return errors.New("categoría inválida: " + string(cat))
	}
}

// SetImageURL asigna la URL de la imagen
func (p *Product) SetImageURL(url string) {
	p.imageURL = url
}

// ============================================================
// MÉTODOS DE NEGOCIO (lógica específica del dominio)
// ============================================================

func (p *Product) IsAvailable() bool {
	return p.stock > 0
}

func (p *Product) IsAvailableQty(qty int) bool {
	return p.stock >= qty
}

// DecreaseStock descuenta stock al vender — con validación
func (p *Product) DecreaseStock(qty int) error {
	if qty <= 0 {
		return errors.New("la cantidad debe ser positiva")
	}
	if p.stock < qty {
		return fmt.Errorf("stock insuficiente para '%s': hay %d, se piden %d", p.name, p.stock, qty)
	}
	p.stock -= qty
	return nil
}

// IncreaseStock agrega stock (reabastecimiento o devolución)
func (p *Product) IncreaseStock(qty int) error {
	if qty <= 0 {
		return errors.New("la cantidad debe ser positiva")
	}
	p.stock += qty
	return nil
}

// FormattedPrice retorna el precio con formato de moneda
func (p *Product) FormattedPrice() string {
	return fmt.Sprintf("$%.2f", p.price)
}

// MarshalJSON — necesario para que encoding/json pueda
// serializar los campos privados al responder la API

func (p *Product) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(
		`{"id":%q,"name":%q,"description":%q,"price":%.2f,"stock":%d,"category":%q,"image_url":%q,"created_at":%q}`,
		p.id, p.name, p.description, p.price, p.stock,
		string(p.category), p.imageURL, p.createdAt.Format(time.RFC3339),
	)), nil
}
