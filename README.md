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

En Go la POO se implementa con `structs` y métodos. La encapsulación se logra escribiendo campos en **minúscula** (privados) y exponiendo **getters** y **setters** públicos (mayúscula). En este proyecto se aplicó POO con **encapsulación**, **constructores con validación**, **composición** y una **máquina de estados** para el flujo de órdenes.

---

### 🧩 Modelo UML técnico (atributos + métodos + relaciones)

**Convenciones UML**
- `-` atributo/método privado  
- `+` método público  
- Multiplicidades: `1`, `0..1`, `1..*`, `0..*`  
- Relación:
  - **◆ composición** (contiene y controla ciclo de vida)
  - **— asociación** (referencia sin propiedad del ciclo de vida)

---

### 📌 Diagrama de clases (UML ASCII con multiplicidades)

```text
┌───────────────────────────────┐
│            Store              │  (1)
├───────────────────────────────┤
│ - mu: sync.Mutex              │
│ - products: map[string]*Prod  │  ◆── (0..*) Product
│ - cart: *Cart                 │  ◆── (1)    Cart
│ - orders: map[string]*Order   │  ◆── (0..*) Order
│ - orderSeq: int               │
│ - prodSeq: int                │
├───────────────────────────────┤
│ + GetAllProducts() []*Product │
│ + GetProduct(id) (*Product,e) │
│ + SearchProducts(q) []*Prod   │
│ + GetProductsByCategory(cat)  │
│ + CreateProduct(...) (*Prod,e)│
│ + UpdateProduct(...) (*Prod,e)│
│ + DeleteProduct(id) error     │
│ + UpdateStock(id,stock)       │
│ + GetCart() *Cart             │
│ + AddToCart(pid,qty) (*Cart,e)│
│ + RemoveFromCart(pid) (*Cart,e)│
│ + ClearCart() *Cart           │
│ + CreateOrder(cust,notes)     │
│ + GetOrder(id) (*Order,e)     │
│ + GetAllOrders() []*Order     │
│ + AdvanceOrderStatus(id)      │
│ + CancelOrder(id)             │
└───────────────────────────────┘
                ◆ (1)
                │
                ▼
┌───────────────────────────────┐
│             Cart              │  (1)
├───────────────────────────────┤
│ - items: []CartItem           │  ◆── (0..*) CartItem
│ - discount: float64           │
├───────────────────────────────┤
│ + AddItem(p *Product,qty) err │
│ + RemoveItem(productID) err   │
│ + Subtotal() float64          │
│ + Total() float64             │
│ + ItemCount() int             │
│ + Clear()                     │
│ + IsEmpty() bool              │
└───────────────────────────────┘
                ◆ (0..*)
                │
                ▼
┌───────────────────────────────┐
│           CartItem            │  (0..*)
├───────────────────────────────┤
│ - productID: string           │  — asociación → Product (1)
│ - productName: string         │
│ - price: float64 (copia)      │
│ - quantity: int               │
│ - imageURL: string            │
├───────────────────────────────┤
│ + SetQuantity(qty) error      │
│ + Subtotal() float64          │
└───────────────────────────────┘
                — (1) referencia por ID
                │
                ▼
┌───────────────────────────────┐
│            Product            │  (0..*)
├───────────────────────────────┤
│ - id: string                  │
│ - name: string                │
│ - description: string         │
│ - price: float64              │
│ - stock: int                  │
│ - category: Category          │
│ - imageURL: string            │
│ - createdAt: time.Time        │
├───────────────────────────────┤
│ + SetName(name) error         │
│ + SetPrice(price) error       │
│ + SetStock(stock) error       │
│ + SetCategory(cat) error      │
│ + DecreaseStock(qty) error    │
│ + IncreaseStock(qty) error    │
│ + MarshalJSON() ([]byte,error)│
└───────────────────────────────┘


┌───────────────────────────────┐
│             Order             │  (0..*)
├───────────────────────────────┤
│ - id: string                  │
│ - customer: Customer          │  ◆── (1) Customer
│ - items: []CartItem           │  ◆── (1..*) CartItem (copia)
│ - total: float64              │
│ - status: OrderStatus         │
│ - notes: string               │
│ - createdAt: time.Time        │
│ - updatedAt: time.Time        │
├───────────────────────────────┤
│ + AdvanceStatus() error       │
│ + Cancel() error              │
│ + IsCancellable() bool        │
│ + Summary() string            │
│ + MarshalJSON() ([]byte,error)│
└───────────────────────────────┘
                ◆ (1)
                │
                ▼
┌───────────────────────────────┐
│           Customer            │  (1)
├───────────────────────────────┤
│ - name: string                │
│ - email: string               │
│ - phone: string               │
│ - address: string             │
│ - city: string                │
├───────────────────────────────┤
│ + Validate() error            │
│ + UnmarshalJSON([]byte) error │
│ + MarshalJSON() ([]byte,error)│
└───────────────────────────────┘

Relaciones (resumen):
- Store (1) ◆── (1) Cart
- Store (1) ◆── (0..*) Product
- Store (1) ◆── (0..*) Order
- Cart  (1) ◆── (0..*) CartItem
- Order (1) ◆── (1) Customer
- Order (1) ◆── (1..*) CartItem
- CartItem (1) —→ (1) Product  (referencia por productID)

```

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

## 🧪 Pruebas realizadas y resultados

Para validar la calidad del proyecto se realizaron **pruebas funcionales** (API y reglas de negocio) y **pruebas exploratorias de usabilidad** (interfaz). El objetivo fue verificar que el flujo completo del e-commerce sea consistente: **catálogo → carrito → checkout → orden → administración**.

---

### 1 Pruebas de usabilidad (interfaz y experiencia)

> Tipo: **prueba exploratoria** (no es un estudio científico formal).  
> Objetivo: comprobar que usuarios no técnicos entienden la interfaz y completan tareas reales sin asistencia.

**Participantes**
- Total: **5 personas** (familiares cercanos)
- Perfil: usuarios sin experiencia técnica (4) + usuario con experiencia básica (1)
- Dispositivos: **3 laptop / 2 celular**
- Duración por participante: **8–12 minutos**

**Escenarios / tareas evaluadas**
- **T1 — Catálogo**
  - Entrar a `/products.html`
  - Filtrar por categoría (ej. “rosa”)
  - Buscar por palabra clave
- **T2 — Carrito**
  - Agregar 2 productos
  - Aumentar / disminuir cantidad
  - Eliminar un producto
- **T3 — Checkout**
  - Completar datos del cliente
  - Confirmar orden
  - Revisar mensaje de confirmación
- **T4 — Admin**
  - Entrar a `/admin.html` e ingresar contraseña
  - Crear un producto
  - Editar un producto
  - Avanzar estado de una orden
  - Cancelar una orden (cuando aplique)

**Métricas y resultados (cuantitativos)**

| Métrica | Resultado |
|--------|----------|
| Usuarios que completaron T1 sin ayuda | 5/5 |
| Usuarios que completaron T2 sin ayuda | 4/5 |
| Usuarios que completaron T3 sin ayuda | 4/5 |
| Usuarios que completaron T4 sin ayuda | 3/5 |
| Tiempo promedio flujo compra (T1–T3) | ~4 min 30 s |
| Errores más frecuentes | confusión en admin (contraseña/acciones), botón de avanzar estado poco obvio, expectativa de mensaje más visible en carrito vacío |

**Hallazgos (cualitativos)**

**Fortalezas detectadas**
- La navegación del catálogo fue considerada “clara y rápida”.
- Las categorías ayudaron a encontrar productos sin leer todo el listado.
- El carrito mostró correctamente subtotal/total y cantidades.

**Oportunidades de mejora detectadas**
- En móvil, algunas acciones del admin requerían más scroll.
- El botón para avanzar estado (▶) no fue intuitivo para 2 participantes.
- 1 participante intentó comprar con carrito vacío y esperaba un mensaje más visible.

**Mejoras aplicadas a partir del feedback**
- Mejora de microcopys (mensajes) en acciones inválidas: carrito vacío, stock insuficiente, validaciones de formulario.
- Mejor feedback visual en admin al crear/editar/eliminar (mensajes de éxito/error).
- Ajustes de distribución responsive en el panel admin para reducir scroll en móvil.

---

### 2 Pruebas funcionales (API + reglas de negocio)

Objetivo: verificar que la lógica de negocio es consistente y que los endpoints mantienen formato uniforme.

#### 2.1 Validaciones principales verificadas

**Productos (`Product`)**
- `NewProduct(...)`:
  - ID obligatorio
  - nombre obligatorio
  - precio > 0
  - stock >= 0
- `DecreaseStock(qty)`:
  - qty > 0
  - no permite vender más del stock disponible
- `SetCategory(cat)`:
  - solo permite categorías definidas (`rosa`, `girasol`, `loto`, `margarita`)

**Cliente (`Customer`)**
- `Validate()`:
  - nombre, email, address y city obligatorios
  - email con formato básico (contiene `@` y `.`)

**Orden (`Order`)**
- `NewOrder(...)`:
  - no permite crear orden si el carrito está vacío
  - copia ítems del carrito para mantener historial
- Máquina de estados:
  - `AdvanceStatus()` solo permite transiciones válidas
  - `Cancel()` no permite cancelar si ya está `enviada` o `entregada`

**Store (memoria + concurrencia)**
- Acceso controlado con `sync.Mutex` para evitar corrupción de datos ante múltiples peticiones.

#### 2.2 Casos de prueba ejecutados (resumen)

| Caso | Entrada | Resultado esperado | Resultado |
|------|--------|-------------------|----------|
| Listar catálogo | GET `/api/products` | JSON con lista de productos | OK |
| Filtrar por categoría | GET `/api/products?category=rosa` | Solo productos categoría rosa | OK |
| Buscar productos | GET `/api/products/search?q=loto` | Coincidencias por nombre/desc/categoría | OK |
| Agregar al carrito | POST `/api/cart/add` `{product_id, quantity}` | Carrito con item agregado/sumado | OK |
| Stock insuficiente | POST add con qty > stock | Error JSON “stock insuficiente” | OK |
| Checkout válido | POST `/api/orders` con cliente válido | Crea orden, limpia carrito | OK |
| Checkout carrito vacío | POST `/api/orders` con carrito vacío | Error JSON “carrito vacío” | OK |
| Avanzar estado | PUT `/api/orders/{id}/status` | Cambia al siguiente estado válido | OK |
| Cancelación válida | PUT `/api/orders/{id}/cancel` antes de enviado | Estado pasa a `cancelada` | OK |
| Cancelación inválida | Cancelar cuando ya está `enviada` | Error JSON “no se puede cancelar” | OK |
| Inventario CRUD | POST/PUT/DELETE `/api/inventory` | Cambios reflejados en catálogo | OK |

#### 2.3 Resultado técnico final
- La API mantuvo el formato estándar:
  - Éxito: `{ "success": true, "data": ... }`
  - Error: `{ "success": false, "error": "..." }`
- Todas las validaciones respondieron con errores descriptivos.
- El Store en memoria no presentó inconsistencias durante pruebas con peticiones repetidas (por uso de `sync.Mutex`).

---

### 3 Conclusión de pruebas

Las pruebas confirmaron que el sistema es funcional y comprensible para usuarios no técnicos. El flujo principal del e-commerce se completó de forma exitosa en la mayoría de casos, y los ajustes aplicados mejoraron la claridad del panel admin y los mensajes de error/estado.

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