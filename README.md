# 🌸 FloriLuz — E-Commerce de Lámparas Florales

Sistema de e-commerce desarrollado en **Go** como proyecto académico de Programación Orientada a Objetos. Permite navegar un catálogo de lámparas florales, agregarlas al carrito, realizar pedidos y administrar el inventario desde un panel protegido.

---

## 🛠 Stack Tecnológico

| Capa | Tecnología | Descripción |
|------|-----------|-------------|
| Frontend | HTML + CSS + JavaScript | Interfaz visual. Sin frameworks externos |
| Backend | Go (librería estándar) | Servidor HTTP y lógica de negocio |
| Datos | Memoria RAM | Sin base de datos externa. Los datos viven en memoria mientras el servidor está activo |
| Comunicación | HTTP / JSON | El navegador se comunica con Go mediante peticiones `fetch()` |
| Despliegue | Docker + Render.com | Imagen multi-stage. Puerto dinámico via variable de entorno `PORT` |

---

## 📁 Estructura del Proyecto

```
ecommerce/
├── main.go                    → punto de entrada, arranca el servidor (17 rutas registradas)
├── go.mod                     → módulo Go — 0 dependencias externas
├── Dockerfile                 → imagen multi-stage para despliegue en producción
│
├── models/                    → CLASES del sistema (POO)
│   ├── product.go             → clase Product + tipo Category
│   ├── cart.go                → clases CartItem y Cart
│   ├── customer.go            → clase Customer
│   └── order.go               → clase Order + tipo OrderStatus
│
├── store/
│   └── store.go               → base de datos en memoria (sync.Mutex, CRUD completo)
│
├── handlers/                  → controladores HTTP
│   ├── helpers.go             → respondJSON, respondError, CORS headers
│   ├── product_handler.go     → catálogo público
│   ├── cart_handler.go        → carrito de compras
│   ├── order_handler.go       → órdenes + máquina de estados
│   └── inventory_handler.go   → CRUD de inventario (panel admin)
│
└── frontend/                  → interfaz visual
    ├── index.html             → página de inicio con catálogo destacado
    ├── products.html          → catálogo completo con búsqueda y filtros
    ├── cart.html              → carrito y checkout
    ├── admin.html             → panel de administración (protegido con contraseña)
    └── style.css              → estilos con variables CSS, diseño responsive
```

---

## 🚀 Cómo ejecutar

```bash
# 1. Entrar a la carpeta
cd ecommerce

# 2. Correr el servidor
go run main.go

# 3. Abrir en el navegador
# http://localhost:8080

# Panel de administración:
# http://localhost:8080/admin.html
# Contraseña: floriluz2024

# Para detener: Ctrl + C
```

---

## 🧱 Clases del Sistema (POO)

En Go la POO se implementa con `structs` y métodos. La encapsulación se logra escribiendo campos en **minúscula** (privados) y exponiendo **getters** y **setters** públicos (mayúscula).

```
  ┌─────────────┐              ┌─────────────┐
  │   Product   │◄─────────── │  CartItem   │
  └─────────────┘  composición└─────────────┘
                                      ▲
  ┌─────────────┐                     │ composición
  │  Customer   │◄──────── ┌──────────┴──┐
  └─────────────┘          │    Cart     │
         ▲                 └─────────────┘
         │ composición
  ┌──────┴──────────────┐
  │        Order        │ (máquina de estados)
  └─────────────────────┘
```

---

### 📦 Product — `models/product.go`

Representa una lámpara floral del catálogo.

**Tipo auxiliar:**
```go
type Category string

const (
    CategoryRose      Category = "rosa"
    CategorySunflower Category = "girasol"
    CategoryLotus     Category = "loto"
    CategoryDaisy     Category = "margarita"
)
```

**Atributos (privados):**

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | `string` | Identificador único. Ej: `lamp-001` |
| `name` | `string` | Nombre de la lámpara |
| `description` | `string` | Descripción detallada |
| `price` | `float64` | Precio en dólares. Debe ser > 0 |
| `stock` | `int` | Unidades disponibles. No puede ser negativo |
| `category` | `Category` | Tipo de flor: rosa, girasol, loto, margarita |
| `imageURL` | `string` | URL de la imagen |
| `createdAt` | `time.Time` | Fecha de creación del registro |

**Constructor:**
```go
func NewProduct(id, name, description string, price float64,
    stock int, category Category, imageURL string) (*Product, error)
```
Valida todos los campos antes de crear. Retorna `error` si algo es inválido.

**Getters:**

| Método | Retorna | Descripción |
|--------|---------|-------------|
| `GetID()` | `string` | ID del producto |
| `GetName()` | `string` | Nombre |
| `GetDescription()` | `string` | Descripción |
| `GetPrice()` | `float64` | Precio |
| `GetStock()` | `int` | Stock disponible |
| `GetCategory()` | `Category` | Categoría (tipo de flor) |
| `GetImageURL()` | `string` | URL de imagen |
| `GetCreatedAt()` | `time.Time` | Fecha de creación |

**Setters (con validación):**

| Método | Retorna | Validación |
|--------|---------|-----------|
| `SetName(name string)` | `error` | No puede estar vacío |
| `SetDescription(desc string)` | `void` | Sin validación especial |
| `SetPrice(price float64)` | `error` | Debe ser mayor a cero |
| `SetStock(stock int)` | `error` | No puede ser negativo |
| `SetCategory(cat Category)` | `error` | Debe ser una de las 4 categorías válidas |
| `SetImageURL(url string)` | `void` | Sin validación especial |

**Métodos de negocio:**

| Método | Retorna | Descripción |
|--------|---------|-------------|
| `IsAvailable()` | `bool` | `true` si stock > 0 |
| `IsAvailableQty(qty int)` | `bool` | `true` si stock >= qty |
| `DecreaseStock(qty int)` | `error` | Descuenta stock al vender. Error si no alcanza |
| `IncreaseStock(qty int)` | `error` | Agrega stock (devoluciones / reabastecimiento) |
| `FormattedPrice()` | `string` | Precio formateado: `"$49.99"` |
| `MarshalJSON()` | `[]byte, error` | Serializa campos privados a JSON para la API |

---

### 👤 Customer — `models/customer.go`

Representa al cliente que realiza la compra.

**Atributos (privados):**

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `name` | `string` | Nombre completo. Obligatorio |
| `email` | `string` | Correo electrónico. Obligatorio. Debe tener `@` y `.` |
| `phone` | `string` | Teléfono. Opcional |
| `address` | `string` | Dirección de entrega. Obligatorio |
| `city` | `string` | Ciudad. Obligatorio |

**Constructor:**
```go
func NewCustomer(name, email, phone, address, city string) (*Customer, error)
```
Llama internamente a cada setter. La validación queda centralizada en los setters.

**Getters:** `GetName()`, `GetEmail()`, `GetPhone()`, `GetAddress()`, `GetCity()`

**Setters (con validación):**

| Método | Validación |
|--------|-----------|
| `SetName(name string) error` | No puede estar vacío. Elimina espacios con `TrimSpace` |
| `SetEmail(email string) error` | Obligatorio. Debe contener `@` y `.` |
| `SetPhone(phone string)` | Opcional. Solo limpia espacios |
| `SetAddress(address string) error` | No puede estar vacío |
| `SetCity(city string) error` | No puede estar vacío |

**Métodos de negocio:**

| Método | Descripción |
|--------|-------------|
| `Validate() error` | Verifica que todos los campos obligatorios estén completos |
| `FullInfo() string` | Resumen: `"nombre <email> — dirección, ciudad"` |
| `MarshalJSON()` | Serializa campos privados a JSON |
| `UnmarshalJSON(data []byte)` | Deserializa el JSON que llega del formulario del frontend |

---

### 🛒 CartItem y Cart — `models/cart.go`

Dos clases que trabajan juntas para el carrito de compras.

#### CartItem

Representa un producto individual dentro del carrito. Guarda una **copia** del precio al momento de agregar (si el precio del producto cambia después, el carrito mantiene el precio original).

**Atributos (privados):**

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `productID` | `string` | ID del producto referenciado |
| `productName` | `string` | Nombre (copia al momento de agregar) |
| `price` | `float64` | Precio unitario (copia, no cambia) |
| `quantity` | `int` | Cantidad de unidades |
| `imageURL` | `string` | URL de la imagen |

**Métodos:**

| Método | Retorna | Descripción |
|--------|---------|-------------|
| `NewCartItem(...)` | `*CartItem, error` | Constructor con validación |
| `GetProductID()` | `string` | Getter del ID |
| `GetProductName()` | `string` | Getter del nombre |
| `GetPrice()` | `float64` | Getter del precio |
| `GetQuantity()` | `int` | Getter de la cantidad |
| `GetImageURL()` | `string` | Getter de la imagen |
| `SetQuantity(qty int)` | `error` | Setter: valida que qty > 0 |
| `Subtotal()` | `float64` | `precio × cantidad` |
| `MarshalJSON()` | `[]byte, error` | Serializa a JSON |

#### Cart

Contenedor de múltiples CartItems.

**Atributos (privados):**

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `items` | `[]CartItem` | Lista de ítems en el carrito |
| `discount` | `float64` | Descuento aplicado. Por defecto 0 |

**Métodos:**

| Método | Retorna | Descripción |
|--------|---------|-------------|
| `NewCart()` | `*Cart` | Constructor: crea carrito vacío |
| `GetItems()` | `[]CartItem` | Getter de ítems |
| `GetDiscount()` | `float64` | Getter del descuento |
| `SetDiscount(d float64)` | `error` | Valida que no supere el subtotal |
| `AddItem(p *Product, qty int)` | `error` | Agrega producto. Si ya existe, suma cantidad |
| `RemoveItem(productID string)` | `error` | Elimina un ítem por ID |
| `Subtotal()` | `float64` | Suma todos los subtotales sin descuento |
| `Total()` | `float64` | Subtotal menos descuento |
| `ItemCount()` | `int` | Total de unidades (no de ítems únicos) |
| `Clear()` | `void` | Vacía el carrito |
| `IsEmpty()` | `bool` | `true` si no hay ítems |
| `MarshalJSON()` | `[]byte, error` | Serializa todo el carrito a JSON |

---

### 📋 Order — `models/order.go`

Orden de compra confirmada. Implementa una **máquina de estados**.

**Tipo auxiliar:**
```go
type OrderStatus string

const (
    StatusPending   OrderStatus = "pendiente"
    StatusPaid      OrderStatus = "pagada"
    StatusPrepared  OrderStatus = "preparada"
    StatusShipped   OrderStatus = "enviada"
    StatusDelivered OrderStatus = "entregada"
    StatusCancelled OrderStatus = "cancelada"
)
```

**Máquina de estados:**
```
pendiente → pagada → preparada → enviada → entregada
    │           │         │                  (fin)
    └───────────┴─────────┴──── cancelada
               (desde cualquier estado antes de enviada)
```

**Atributos (privados):**

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | `string` | ID único. Ej: `ORD-0001` |
| `customer` | `Customer` | Copia completa del cliente al momento de la compra |
| `items` | `[]CartItem` | Copia de los ítems del carrito |
| `total` | `float64` | Total al momento de crear la orden |
| `status` | `OrderStatus` | Estado actual en la máquina de estados |
| `notes` | `string` | Notas opcionales de entrega |
| `createdAt` | `time.Time` | Fecha de creación |
| `updatedAt` | `time.Time` | Fecha de última modificación |

**Constructor:**
```go
func NewOrder(id string, customer Customer, cart *Cart) (*Order, error)
```
Valida cliente y carrito. Copia los ítems del carrito (la orden es independiente del carrito original).

**Getters:** `GetID()`, `GetCustomer()`, `GetItems()`, `GetTotal()`, `GetStatus()`, `GetNotes()`, `GetCreatedAt()`, `GetUpdatedAt()`

**Setters:**

| Método | Descripción |
|--------|-------------|
| `SetNotes(notes string)` | Asigna notas y actualiza `updatedAt` |

> ⚠️ **`status` NO tiene setter público.** El estado solo puede cambiar a través de `AdvanceStatus()` y `Cancel()`. Esto protege la máquina de estados: nadie puede poner una orden en un estado arbitrario.

**Métodos de negocio:**

| Método | Retorna | Descripción |
|--------|---------|-------------|
| `AdvanceStatus()` | `error` | Avanza al siguiente estado válido. Error si ya llegó al final |
| `Cancel()` | `error` | Cancela la orden. Error si ya fue enviada o entregada |
| `IsCancellable()` | `bool` | `true` si todavía se puede cancelar |
| `IsPending()` | `bool` | `true` si está en estado pendiente |
| `Summary()` | `string` | Resumen: `"Orden #ORD-0001 | Cliente | $49.99 | pendiente"` |
| `MarshalJSON()` | `[]byte, error` | Serializa todos los campos privados incluyendo customer e items |

---

## 🗄 Store — `store/store.go`

Base de datos en memoria. Todos los handlers comparten la misma instancia del Store.

**Atributos:**

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `mu` | `sync.Mutex` | Evita corrupción de datos si dos peticiones llegan al mismo tiempo |
| `products` | `map[string]*Product` | Catálogo de productos indexado por ID |
| `cart` | `*Cart` | El carrito activo |
| `orders` | `map[string]*Order` | Historial de órdenes indexado por ID |
| `orderSeq` | `int` | Contador para IDs de órdenes: ORD-0001, ORD-0002... |
| `prodSeq` | `int` | Contador para IDs de productos: lamp-007, lamp-008... |

**Métodos de productos:** `AddProduct`, `CreateProduct`, `UpdateProduct`, `DeleteProduct`, `GetProduct`, `GetAllProducts`, `GetProductsByCategory`, `SearchProducts`, `UpdateStock`

**Métodos del carrito:** `GetCart`, `AddToCart`, `RemoveFromCart`, `ClearCart`

**Métodos de órdenes:** `CreateOrder`, `GetOrder`, `GetAllOrders`, `AdvanceOrderStatus`, `CancelOrder`

**Flujo de `CreateOrder` (el más importante):**
1. Verifica que el carrito no esté vacío
2. Genera el ID automáticamente (`ORD-0001`, `ORD-0002`...)
3. Llama a `models.NewOrder()` que valida cliente y copia ítems
4. Por cada ítem, llama a `product.DecreaseStock()` usando `GetProductID()` y `GetQuantity()`
5. Si algún `DecreaseStock` falla, retorna error y no guarda nada
6. Guarda la orden en el mapa
7. Llama a `cart.Clear()`
8. Retorna la orden creada

---

## 🌐 API REST — 17 Endpoints

### Catálogo público

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/products` | Todos los productos. Acepta `?category=rosa` |
| GET | `/api/products/{id}` | Un producto por ID |
| GET | `/api/products/search?q=` | Búsqueda por nombre, descripción o categoría |

### Carrito

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/cart` | Estado actual del carrito |
| POST | `/api/cart/add` | Body: `{"product_id":"lamp-001","quantity":2}` |
| POST | `/api/cart/remove` | Body: `{"product_id":"lamp-001"}` |
| POST | `/api/cart/clear` | Vacía el carrito |

### Órdenes

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/api/orders` | Crea una orden con datos del cliente |
| GET | `/api/orders/list` | Lista todas las órdenes |
| GET | `/api/orders/{id}` | Consulta una orden específica |
| PUT | `/api/orders/{id}/status` | Avanza al siguiente estado |
| PUT | `/api/orders/{id}/cancel` | Cancela la orden |

### Inventario (admin)

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/api/inventory` | Lista todo el inventario con stocks |
| POST | `/api/inventory` | Crea un producto nuevo (ID autogenerado) |
| PUT | `/api/inventory/{id}` | Edita un producto existente |
| DELETE | `/api/inventory/{id}` | Elimina un producto |
| PUT | `/api/inventory/{id}/stock` | Actualiza solo el stock |

**Formato de respuesta (siempre el mismo):**
```json
// Éxito:
{ "success": true, "data": { ... } }

// Error:
{ "success": false, "error": "mensaje descriptivo" }
```

---

## 💻 Frontend

| Archivo | Descripción |
|---------|-------------|
| `index.html` | Página principal con hero y catálogo destacado |
| `products.html` | Catálogo completo con búsqueda en tiempo real y tabs de categorías |
| `cart.html` | Carrito con formulario de checkout. Muestra 4 estados: cargando, vacío, con ítems, orden confirmada |
| `admin.html` | Panel de administración protegido con contraseña. Dashboard, CRUD de inventario y gestión de órdenes |
| `style.css` | Estilos con variables CSS, navbar sticky con efecto glass, responsive completo |

### Panel de Administración

Accesible en `/admin.html`. Requiere contraseña (`floriluz2024`). Incluye:

- **Dashboard** — estadísticas en tiempo real: total de productos, órdenes, productos agotados y stock bajo
- **Inventario** — tabla completa con badges de stock. Permite crear, editar, actualizar stock y eliminar productos
- **Órdenes** — tabla con todas las órdenes. Botón para avanzar estado (▶) y cancelar (✖)

La autenticación usa `sessionStorage`: al cerrar la pestaña o el navegador, se pide la contraseña nuevamente.

---

## ✅ Conceptos de POO Aplicados

| Concepto | Cómo se aplicó |
|----------|---------------|
| **Encapsulación** | Todos los campos de las clases son privados (minúscula). Solo se accede mediante getters y setters |
| **Constructores** | Cada clase tiene `New<Clase>()` que valida antes de crear el objeto |
| **Getters** | Métodos `Get<Campo>()` para lectura. Ej: `GetPrice()`, `GetStock()` |
| **Setters** | Métodos `Set<Campo>()` con validación antes de modificar. Ej: `SetPrice()` valida > 0 |
| **Métodos de negocio** | Lógica encapsulada en la clase. Ej: `DecreaseStock()`, `AdvanceStatus()` |
| **Composición** | `Order` contiene `Customer` y `[]CartItem`. `Cart` contiene `[]CartItem` |
| **Manejo de errores** | Cada método que puede fallar retorna `error`. Los errores se propagan hasta el handler |
| **Interfaces implícitas** | `MarshalJSON()` implementa la interfaz `json.Marshaler` de Go en cada clase |
| **Máquina de estados** | `Order` controla transiciones válidas. `status` no tiene setter público para proteger la lógica |

---

## 🐳 Despliegue con Docker

```bash
# Construir la imagen
docker build -t floriluz .

# Correr el contenedor
docker run -p 8080:8080 floriluz
```

El `Dockerfile` usa multi-stage build: compila el binario en `golang:1.21-alpine` y lo copia a una imagen `alpine` limpia. La imagen final no contiene Go instalado, solo el binario.

Para desplegar en **Render.com**: conectar el repositorio de GitHub, seleccionar Docker como runtime y hacer deploy. El sistema lee el puerto de la variable de entorno `PORT` automáticamente.

---

## 👩‍💻 Autora

Proyecto desarrollado como parte de la materia de **Programación Orientada a Objetos**  
Giuliana Moreta — Tercer Semestre — Ingeniería en Software — UIDE — 2024