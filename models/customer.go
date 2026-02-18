// @name: Giuliana Moreta
//@date: 17/02/2026
//@Materia: Programación orientada a objetos
//@Curso: 3er Semestre
//@Carrera: Ingeniería en Software

// models/customer.go
// Clase Customer con encapsulación completa
package models

import (
	"errors"
	"strings"
)

// Customer — campos privados
type Customer struct {
	name    string
	email   string
	phone   string
	address string
	city    string
}

func NewCustomer(name, email, phone, address, city string) (*Customer, error) {
	c := &Customer{}

	// Usamos los setters para validar cada campo
	if err := c.SetName(name); err != nil {
		return nil, err
	}
	if err := c.SetEmail(email); err != nil {
		return nil, err
	}
	c.SetPhone(phone)
	if err := c.SetAddress(address); err != nil {
		return nil, err
	}
	if err := c.SetCity(city); err != nil {
		return nil, err
	}

	return c, nil
}

// GETTERS

func (c *Customer) GetName() string    { return c.name }
func (c *Customer) GetEmail() string   { return c.email }
func (c *Customer) GetPhone() string   { return c.phone }
func (c *Customer) GetAddress() string { return c.address }
func (c *Customer) GetCity() string    { return c.city }

// SETTERS con validación

func (c *Customer) SetName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("el nombre del cliente es obligatorio")
	}
	c.name = name
	return nil
}

func (c *Customer) SetEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("el correo electrónico es obligatorio")
	}
	// Validación básica: debe tener @ y un punto después
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.New("el correo electrónico no tiene un formato válido")
	}
	c.email = email
	return nil
}

// SetPhone no es obligatorio, pero limpia espacios
func (c *Customer) SetPhone(phone string) {
	c.phone = strings.TrimSpace(phone)
}

func (c *Customer) SetAddress(address string) error {
	address = strings.TrimSpace(address)
	if address == "" {
		return errors.New("la dirección de entrega es obligatoria")
	}
	c.address = address
	return nil
}

func (c *Customer) SetCity(city string) error {
	city = strings.TrimSpace(city)
	if city == "" {
		return errors.New("la ciudad es obligatoria")
	}
	c.city = city
	return nil
}

// MÉTODOS DE NEGOCIO

// Validate verifica que todos los campos obligatorios estén completos
func (c *Customer) Validate() error {
	if c.name == "" {
		return errors.New("el nombre es obligatorio")
	}
	if c.email == "" {
		return errors.New("el correo es obligatorio")
	}
	if c.address == "" {
		return errors.New("la dirección es obligatoria")
	}
	if c.city == "" {
		return errors.New("la ciudad es obligatoria")
	}
	return nil
}

// FullInfo retorna un resumen del cliente
func (c *Customer) FullInfo() string {
	return c.name + " <" + c.email + "> — " + c.address + ", " + c.city
}

// MarshalJSON para serializar campos privados

func (c *Customer) MarshalJSON() ([]byte, error) {
	return []byte(`{"name":"` + c.name + `","email":"` + c.email +
		`","phone":"` + c.phone + `","address":"` + c.address +
		`","city":"` + c.city + `"}`), nil
}

// UnmarshalJSON permite deserializar JSON al recibir datos del frontend
func (c *Customer) UnmarshalJSON(data []byte) error {
	// Usamos una struct auxiliar temporal con campos públicos
	var aux struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Address string `json:"address"`
		City    string `json:"city"`
	}

	// Importamos encoding/json aquí vía interfaz para evitar import circular
	// Usamos parsing manual simple
	str := string(data)
	aux.Name = extractJSON(str, "name")
	aux.Email = extractJSON(str, "email")
	aux.Phone = extractJSON(str, "phone")
	aux.Address = extractJSON(str, "address")
	aux.City = extractJSON(str, "city")

	c.name = aux.Name
	c.email = aux.Email
	c.phone = aux.Phone
	c.address = aux.Address
	c.city = aux.City
	return nil
}

// extractJSON extrae el valor de un campo JSON simple
func extractJSON(json, key string) string {
	search := `"` + key + `":"`
	start := strings.Index(json, search)
	if start == -1 {
		return ""
	}
	start += len(search)
	end := strings.Index(json[start:], `"`)
	if end == -1 {
		return ""
	}
	return json[start : start+end]
}
