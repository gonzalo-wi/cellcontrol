# Explicación Detallada: User Service

## Visión General
Este archivo implementa la **capa de servicios (Service Layer)** para la gestión de usuarios. Es el **cerebro de la lógica de negocio** - aquí es donde viven las reglas, validaciones y transformaciones de datos antes de llegar a la base de datos.

---

## ¿Qué es la Capa de Servicios?

La Service Layer es el **intermediario** entre los Handlers (HTTP) y los Repositories (Base de datos).

```
┌─────────────┐
│   HANDLER   │  ← Recibe requests HTTP
│ (Controller)│     Valida entrada básica
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   SERVICE   │  ← ESTÁS AQUÍ - Lógica de negocio
│  (Business) │     Validaciones complejas
└──────┬──────┘     Transformaciones
       │             Orquestación
       ▼
┌─────────────┐
│ REPOSITORY  │  ← Acceso a datos
│    (Data)   │     CRUD operations
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  DATABASE   │  ← MySQL/PostgreSQL
└─────────────┘
```

**¿Por qué esta separación?**
- ✅ El Handler no conoce la BD (solo conoce el Service)
- ✅ El Repository no conoce HTTP (solo conoce datos)
- ✅ El Service orquesta todo sin conocer detalles técnicos

---

## Desglose Línea por Línea

### Package e Imports

```go
package service
```
- Define el paquete `service`
- Todos los servicios de tu app van aquí
- Convención: un servicio por entidad (UserService, CellphoneService, etc.)

```go
import (
	"strings"

	"github.com/gonzalo-wi/cellcontrol/internal/domain"
	"github.com/gonzalo-wi/cellcontrol/internal/repository"
)
```

**Imports explicados:**

1. **`"strings"`**
    - Paquete estándar de Go para manipulación de strings
    - Usado para `TrimSpace()`, `ToLower()`, etc.
    - Limpia y normaliza datos de entrada

2. **`"github.com/gonzalo-wi/cellcontrol/internal/domain"`**
    - Tu paquete de **modelos de dominio**
    - Contiene `domain.User` - la estructura de usuario
    - Define "qué es" un usuario en tu sistema

3. **`"github.com/gonzalo-wi/cellcontrol/internal/repository"`**
    - Tu paquete de **repositorios**
    - Contiene `repository.UserRepository` - la interface de acceso a datos
    - El Service usa esto para guardar/obtener usuarios

**Nota:** La línea en blanco entre `strings` y los demás imports es **convención en Go**:
- Primero: paquetes estándar de Go
- Luego: paquetes externos (terceros)
- Finalmente: tus paquetes internos

---

## Interface UserService

```go
type UserService interface {
	CreateUser(nombre, apellido, email, reparto string) error
	GetAllUsers() ([]domain.User, error)
}
```

### ¿Por qué una Interface?

Una **interface** define el **contrato público** del servicio - qué operaciones están disponibles.

**Ventajas:**
- 🧪 **Testing**: Puedes crear mocks fácilmente
- 🔌 **Desacoplamiento**: Los handlers dependen de la interface, no de la implementación
- 🔄 **Intercambiabilidad**: Puedes cambiar la implementación sin afectar otros componentes
- 📋 **Documentación**: La interface es el "manual de uso" del servicio

### Métodos de la Interface

#### 1. CreateUser

```go
CreateUser(nombre, apellido, email, reparto string) error
```

**Propósito**: Crear un nuevo usuario en el sistema

**Parámetros (4 strings):**
- `nombre`: Nombre del usuario (ej: "Juan")
- `apellido`: Apellido del usuario (ej: "Pérez")
- `email`: Email del usuario (ej: "juan@example.com")
- `reparto`: Departamento o área de trabajo (ej: "Ventas")

**¿Por qué recibe strings y no `domain.User`?**
- 🎯 **Separación de responsabilidades**: El que llama (handler) no necesita construir un `User` completo
- ✅ **Validación centralizada**: El Service se encarga de crear y validar el User
- 🧹 **Limpieza de datos**: El Service limpia y normaliza antes de crear

**Retorno:**
- `error`:
    - `nil` si se creó exitosamente
    - Error si falló (email duplicado, error de BD, etc.)

**Ejemplo de uso:**
```go
err := userService.CreateUser("Juan", "Pérez", "juan@example.com", "Ventas")
if err != nil {
	log.Printf("Error creando usuario: %v", err)
}
```

#### 2. GetAllUsers

```go
GetAllUsers() ([]domain.User, error)
```

**Propósito**: Obtener todos los usuarios del sistema

**Parámetros:**
- Ninguno (obtiene TODOS los usuarios)

**Retorno (2 valores):**
1. `[]domain.User`: Slice con todos los usuarios
2. `error`: Error si falló, o `nil` si todo bien

**Ejemplo de uso:**
```go
users, err := userService.GetAllUsers()
if err != nil {
	log.Printf("Error obteniendo usuarios: %v", err)
	return
}
fmt.Printf("Total de usuarios: %d\n", len(users))
```

---

## Implementación: Estructura userService

```go
type userService struct {
	repo repository.UserRepository
}
```

### Características Importantes

**1. Nombre en minúscula (privado)**
- `userService` vs `UserService`
- Privado → solo accesible dentro del paquete
- Nadie puede instanciarlo directamente desde afuera

**2. Campo `repo`**
- Tipo: `repository.UserRepository` (es una **interface**)
- Almacena la referencia al repositorio
- Permite **inyección de dependencias**

**¿Qué es Dependency Injection (DI)?**
```go
// El Service NO crea su repositorio
// El repositorio se "inyecta" desde afuera
userRepo := repository.NewUserRepository(db)
userService := service.NewUserService(userRepo) // ← Inyección
```

**Ventajas de DI:**
- ✅ **Testing**: Puedes inyectar un mock
- ✅ **Flexibilidad**: Puedes cambiar la implementación
- ✅ **Control**: Quien crea el service controla las dependencias

---

## Constructor: NewUserService

```go
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
```

### Análisis Detallado

**Patrón Factory Method**
- Es el **único** way oficial de crear un `UserService`
- Encapsula la creación de la instancia
- Retorna la interface, no la struct

**Paso a paso:**

1. **Recibe `repo repository.UserRepository`**
    - Ya debe estar inicializado
    - Se pasa desde afuera (inyección de dependencias)

2. **Crea `&userService{repo: repo}`**
    - `&` → toma la dirección (devuelve un puntero)
    - `userService{...}` → crea la struct
    - `repo: repo` → asigna el campo (sintaxis corta cuando nombres coinciden)

3. **Retorna `UserService`**
    - Retorna la **interface**, no el tipo concreto
    - `*userService` implementa `UserService` automáticamente
    - Quien use esto solo ve la interface

**Ejemplo completo de inicialización:**
```go
// 1. Conectar a BD
db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 2. Crear repositorio
userRepo := repository.NewUserRepository(db)

// 3. Crear servicio (inyectando el repo)
userService := service.NewUserService(userRepo)

// 4. Usar el servicio
err := userService.CreateUser("Juan", "Pérez", "juan@mail.com", "IT")
```

---

## Método CreateUser - El Corazón de la Lógica

```go
func (s *userService) CreateUser(nombre, apellido, email, reparto string) error {
	u := &domain.User{
		Nombre:   strings.TrimSpace(nombre),
		Apellido: strings.TrimSpace(apellido),
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Reparto:  strings.TrimSpace(reparto),
	}
	if err := s.repo.CreateUser(u); err != nil {
		return err
	}
	return nil
}
```

### Desglose Línea por Línea

#### Firma del Método

```go
func (s *userService) CreateUser(nombre, apellido, email, reparto string) error
```

- **`(s *userService)`**: Receptor del método
    - `s` es la instancia (convención: primera letra del tipo)
    - `*userService` es el tipo receptor (puntero a struct)
    - Hace que `CreateUser` sea un **método** de `*userService`

- **Parámetros**: 4 strings (datos crudos del usuario)

- **Retorno**: `error` (nil = éxito, error = falló)

#### Construcción del Usuario (Líneas 24-29)

```go
u := &domain.User{
	Nombre:   strings.TrimSpace(nombre),
	Apellido: strings.TrimSpace(apellido),
	Email:    strings.ToLower(strings.TrimSpace(email)),
	Reparto:  strings.TrimSpace(reparto),
}
```

**¿Qué está pasando aquí?**

1. **Crea un puntero a `domain.User`**:
    - `&domain.User{...}` → literal de struct + dirección
    - `u` es de tipo `*domain.User` (puntero)

2. **Aplica transformaciones a cada campo**:

**a) Nombre y Apellido:**
```go
Nombre: strings.TrimSpace(nombre)
```
- `strings.TrimSpace()` elimina espacios al inicio y final
- Entrada: `"  Juan  "` → Salida: `"Juan"`

**b) Email (doble transformación):**
```go
Email: strings.ToLower(strings.TrimSpace(email))
```
Ejecuta de **adentro hacia afuera**:
1. `strings.TrimSpace(email)` → quita espacios
    - `" Juan@EXAMPLE.com "` → `"Juan@EXAMPLE.com"`
2. `strings.ToLower(...)` → convierte a minúsculas
    - `"Juan@EXAMPLE.com"` → `"juan@example.com"`

**¿Por qué normalizar el email?**
- ✅ Emails son case-insensitive: `Juan@Mail.com` = `juan@mail.com`
- ✅ Evita duplicados por mayúsculas/minúsculas
- ✅ Consistencia en la BD

**c) Reparto:**
```go
Reparto: strings.TrimSpace(reparto)
```
- Igual que nombre/apellido, solo quita espacios

**Resultado final:**
```go
// Entrada:
CreateUser("  Juan  ", " Pérez ", " Juan@EXAMPLE.com ", "  IT  ")

// Usuario creado:
&domain.User{
	Nombre:   "Juan",
	Apellido: "Pérez",
	Email:    "juan@example.com",
	Reparto:  "IT",
}
```

#### Guardar en Base de Datos (Líneas 30-33)

```go
if err := s.repo.CreateUser(u); err != nil {
	return err
}
return nil
```

**Análisis detallado:**

**Línea 30:**
```go
if err := s.repo.CreateUser(u); err != nil {
```

Esta línea hace **3 cosas**:
1. `s.repo` → accede al repositorio
2. `.CreateUser(u)` → llama al método del repositorio
3. `err := ...` → asigna el error retornado
4. `if ... != nil` → verifica si hubo error

Es una forma compacta de:
```go
err := s.repo.CreateUser(u)
if err != nil {
	return err
}
```

**Línea 31:**
```go
return err
```
- Si hubo error → lo propaga hacia arriba
- El handler/caller debe decidir qué hacer con el error

**Línea 33:**
```go
return nil
```
- Si llegamos aquí → todo bien
- Retornamos `nil` (sin error)

### Flujo Completo del CreateUser

```
1. Handler llama:
   service.CreateUser("  Juan  ", "Pérez", "JUAN@MAIL.COM", "IT")
              ↓
2. Service limpia y normaliza:
   - Nombre: "Juan" (sin espacios)
   - Email: "juan@mail.com" (minúsculas, sin espacios)
              ↓
3. Service crea &domain.User{...}
              ↓
4. Service llama:
   repo.CreateUser(user)
              ↓
5. Repository ejecuta:
   INSERT INTO users (...) VALUES (...)
              ↓
6. Si éxito → retorna nil
   Si error → retorna error
              ↓
7. Service propaga el resultado al Handler
```

---

## Método GetAllUsers

```go
func (s *userService) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAllUsers()
}
```

### Análisis

**¿Por qué tan simple?**
- Este método es un **pass-through** (delegación directa)
- No hay lógica de negocio que aplicar
- Solo obtiene los datos del repositorio

**Firma:**
```go
func (s *userService) GetAllUsers() ([]domain.User, error)
```
- Receptor: `(s *userService)`
- Sin parámetros
- Retorna: lista de usuarios + error

**Implementación:**
```go
return s.repo.GetAllUsers()
```
- Llama directamente al repositorio
- Retorna exactamente lo que el repositorio retorna

**¿Cuándo añadirías lógica aquí?**

Ejemplos de lógica que podrías agregar:

```go
func (s *userService) GetAllUsers() ([]domain.User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	
	// Lógica de negocio adicional:
	
	// 1. Filtrar usuarios activos
	activeUsers := filterActive(users)
	
	// 2. Ordenar por nombre
	sort.Slice(activeUsers, func(i, j int) bool {
		return activeUsers[i].Nombre < activeUsers[j].Nombre
	})
	
	// 3. Ocultar emails parcialmente (privacidad)
	for i := range activeUsers {
		activeUsers[i].Email = maskEmail(activeUsers[i].Email)
	}
	
	return activeUsers, nil
}
```

Pero si no necesitas transformaciones, el pass-through está perfecto.

---

## Arquitectura Completa: Cómo Encaja Todo

### Flujo de una Request HTTP

```
┌─────────────────────────────────────────────────────┐
│ 1. CLIENTE (Postman/Browser)                       │
│    POST /api/users                                  │
│    Body: {"nombre":"Juan","apellido":"Pérez",...}  │
└───────────────────────┬─────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│ 2. HANDLER (user_handler.go)                       │
│    - Parsea JSON                                    │
│    - Valida campos obligatorios                     │
│    - Extrae: nombre, apellido, email, reparto       │
└───────────────────────┬─────────────────────────────┘
                        │
                        │ service.CreateUser(nombre, apellido, email, reparto)
                        ▼
┌─────────────────────────────────────────────────────┐
│ 3. SERVICE (user_service.go) ◄── ESTÁS AQUÍ       │
│    - Limpia datos (TrimSpace)                      │
│    - Normaliza email (ToLower)                     │
│    - Crea domain.User                              │
│    - Valida reglas de negocio (si las hay)         │
└───────────────────────┬─────────────────────────────┘
                        │
                        │ repo.CreateUser(user)
                        ▼
┌─────────────────────────────────────────────────────┐
│ 4. REPOSITORY (user_repository.go)                 │
│    - Ejecuta db.Create(user)                       │
│    - Traduce Go → SQL                              │
└───────────────────────┬─────────────────────────────┘
                        │
                        │ INSERT INTO users ...
                        ▼
┌─────────────────────────────────────────────────────┐
│ 5. DATABASE (MySQL)                                │
│    - Guarda el registro                            │
│    - Valida constraints (UNIQUE email, NOT NULL)   │
│    - Retorna ID generado                           │
└─────────────────────────────────────────────────────┘
```

### Flujo de Respuesta (camino inverso)

```
Database → Repository → Service → Handler → Cliente
  (ID)       (error)     (error)   (JSON)    (200 OK)
```

---

## Responsabilidades de Cada Capa

| Capa | Responsabilidades | NO debe hacer |
|------|-------------------|---------------|
| **Handler** | - Parsear HTTP<br>- Validar formato<br>- Responder HTTP | - Lógica de negocio<br>- Acceso a BD |
| **Service** | - Lógica de negocio<br>- Validaciones complejas<br>- Transformaciones | - Conocer HTTP<br>- SQL directo |
| **Repository** | - CRUD operations<br>- Queries SQL<br>- Transacciones | - Lógica de negocio<br>- Conocer HTTP |

---

## Ejemplo de Uso Completo

### 1. Inicialización (main.go)

```go
package main

import (
	"log"
	
	"github.com/gonzalo-wi/cellcontrol/internal/config"
	"github.com/gonzalo-wi/cellcontrol/internal/repository"
	"github.com/gonzalo-wi/cellcontrol/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Cargar configuración
	cfg := config.MustLoad()
	
	// Conectar a BD
	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Error conectando a BD:", err)
	}
	
	// Crear repositorio
	userRepo := repository.NewUserRepository(db)
	
	// Crear servicio (inyecta el repositorio)
	userService := service.NewUserService(userRepo)
	
	// Usar el servicio
	err = userService.CreateUser("Juan", "Pérez", "juan@example.com", "IT")
	if err != nil {
		log.Fatal("Error creando usuario:", err)
	}
	log.Println("✅ Usuario creado exitosamente")
	
	// Obtener todos los usuarios
	users, err := userService.GetAllUsers()
	if err != nil {
		log.Fatal("Error obteniendo usuarios:", err)
	}
	
	log.Printf("📋 Total de usuarios: %d\n", len(users))
	for _, user := range users {
		log.Printf("   - %s %s (%s)\n", user.Nombre, user.Apellido, user.Email)
	}
}
```

### 2. En un Handler (user_handler.go)

```go
package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/gonzalo-wi/cellcontrol/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// POST /api/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Nombre   string `json:"nombre" binding:"required"`
		Apellido string `json:"apellido" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Reparto  string `json:"reparto" binding:"required"`
	}
	
	// Parsear JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Llamar al servicio (¡aquí es donde se usa!)
	err := h.service.CreateUser(req.Nombre, req.Apellido, req.Email, req.Reparto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando usuario"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "Usuario creado exitosamente"})
}

// GET /api/users
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo usuarios"})
		return
	}
	
	c.JSON(http.StatusOK, users)
}
```

### 3. Request/Response Example

**Request:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "  Juan  ",
    "apellido": "Pérez",
    "email": "JUAN@EXAMPLE.COM",
    "reparto": "  IT  "
  }'
```

**Lo que pasa internamente:**
1. Handler recibe: `"  Juan  "`, `"JUAN@EXAMPLE.COM"`, etc.
2. Service limpia: `"Juan"`, `"juan@example.com"`, etc.
3. Repository guarda en BD
4. BD almacena datos limpios y normalizados

**Response:**
```json
{
  "message": "Usuario creado exitosamente"
}
```

---

## Ventajas de esta Implementación

### ✅ 1. Lógica de Negocio Centralizada
- Todas las transformaciones están en UN lugar
- Si cambias la regla (ej: emails en mayúsculas), solo cambias el Service

### ✅ 2. Datos Limpios y Consistentes
- Siempre se eliminan espacios
- Emails siempre en minúsculas
- No importa cómo lleguen los datos

### ✅ 3. Separación de Responsabilidades
- Handler: HTTP
- Service: Lógica de negocio
- Repository: Datos

### ✅ 4. Fácil de Testear

```go
// Mock del repositorio para testing
type mockUserRepository struct{}

func (m *mockUserRepository) CreateUser(user *domain.User) error {
	// Verificar que el email esté en minúsculas
	if user.Email != strings.ToLower(user.Email) {
		return errors.New("email debe estar en minúsculas")
	}
	return nil
}

// Test
func TestCreateUser(t *testing.T) {
	mockRepo := &mockUserRepository{}
	userService := service.NewUserService(mockRepo)
	
	// Test que el servicio limpia el email
	err := userService.CreateUser("Juan", "Pérez", "JUAN@MAIL.COM", "IT")
	assert.NoError(t, err) // Debe pasar porque el Service normaliza
}
```

### ✅ 5. Reutilizable
```go
// El mismo servicio puede usarse desde:
// - HTTP handlers
// - gRPC services
// - CLI commands
// - Background jobs
```

---

## Mejoras Posibles

### 🚀 1. Validaciones Más Robustas

```go
func (s *userService) CreateUser(nombre, apellido, email, reparto string) error {
	// Validar que no estén vacíos (después de trim)
	if strings.TrimSpace(nombre) == "" {
		return errors.New("nombre es requerido")
	}
	if strings.TrimSpace(apellido) == "" {
		return errors.New("apellido es requerido")
	}
	
	// Validar formato de email
	if !isValidEmail(email) {
		return errors.New("email inválido")
	}
	
	// Crear usuario...
}
```

### 🚀 2. Verificar Duplicados

```go
func (s *userService) CreateUser(nombre, apellido, email, reparto string) error {
	// Limpiar email
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	
	// Verificar si ya existe
	existing, err := s.repo.GetUserByEmail(cleanEmail)
	if err == nil && existing != nil {
		return errors.New("el email ya está registrado")
	}
	
	// Crear usuario...
}
```

### 🚀 3. Logging

```go
func (s *userService) CreateUser(nombre, apellido, email, reparto string) error {
	logger.Info("Creando usuario: %s %s", nombre, apellido)
	
	u := &domain.User{...}
	
	if err := s.repo.CreateUser(u); err != nil {
		logger.Error("Error creando usuario: %v", err)
		return err
	}
	
	logger.Info("Usuario creado exitosamente con ID: %d", u.ID)
	return nil
}
```

### 🚀 4. Más Métodos en la Interface

```go
type UserService interface {
	CreateUser(nombre, apellido, email, reparto string) error
	GetAllUsers() ([]domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(id uint, nombre, apellido, email, reparto string) error
	DeleteUser(id uint) error
	SearchUsers(query string) ([]domain.User, error)
}
```

---

## Conceptos Clave de Go Usados

| Concepto | Dónde se usa | Explicación |
|----------|--------------|-------------|
| **Interface** | `UserService` | Contrato de métodos |
| **Struct** | `userService` | Implementación concreta |
| **Método con receptor** | `(s *userService)` | Método de una struct |
| **Puntero** | `*userService`, `&domain.User` | Referencia a memoria |
| **Dependency Injection** | Constructor | Pasar dependencias |
| **Error handling** | `if err != nil` | Patrón Go de errores |
| **String manipulation** | `strings.TrimSpace` | Transformaciones |
| **Public/Private** | `UserService` vs `userService` | Visibilidad |

---

## Resumen Ejecutivo

**¿Qué hace este código?**
Implementa la lógica de negocio para gestionar usuarios, limpiando datos y coordinando entre handlers y repositorios.

**¿Por qué está diseñado así?**
- Interface pública → flexibilidad
- Struct privada → encapsulación
- Transformaciones → datos consistentes
- Separación de capas → mantenibilidad

**¿Cuándo se usa?**
Desde los handlers HTTP cuando necesitas crear o listar usuarios.

**¿Qué hace especial?**
- ✅ Limpia espacios en blanco
- ✅ Normaliza emails a minúsculas
- ✅ Separa lógica de negocio de acceso a datos
- ✅ Fácil de testear y mantener

---

## Diagrama de Dependencias

```
┌──────────────────┐
│  user_handler.go │
└────────┬─────────┘
         │ usa
         ▼
┌──────────────────┐
│  UserService     │ ◄── Interface (contrato)
│  (interface)     │
└────────┬─────────┘
         │ implementada por
         ▼
┌──────────────────┐
│  userService     │ ◄── ESTE ARCHIVO
│  (struct)        │
└────────┬─────────┘
         │ usa
         ▼
┌──────────────────┐
│ UserRepository   │ ◄── Interface del repo
│  (interface)     │
└────────┬─────────┘
         │ implementada por
         ▼
┌──────────────────┐
│ userRepository   │ ◄── user_repository.go
│  (struct)        │
└────────┬─────────┘
         │ usa
         ▼
┌──────────────────┐
│  GORM / MySQL    │
└──────────────────┘
```

---

¿Necesitas que profundice en alguna parte específica, implemente alguna mejora, o explique cómo conectar esto con los handlers?
