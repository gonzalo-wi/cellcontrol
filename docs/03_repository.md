# Explicación Detallada: User Repository

## Visión General
Este archivo implementa el **patrón Repository** para gestionar la persistencia de datos de usuarios en la base de datos. Actúa como una capa de abstracción entre tu lógica de negocio y la base de datos, usando **GORM** como ORM (Object-Relational Mapping).

---

## ¿Qué es el Patrón Repository?

El patrón Repository es una **capa de abstracción** que:
- ✅ Separa la lógica de acceso a datos de la lógica de negocio
- ✅ Facilita cambiar de base de datos sin afectar el resto del código
- ✅ Permite hacer testing fácilmente (con mocks)
- ✅ Centraliza todas las operaciones de base de datos

**Flujo de datos:**
```
Handler → Service → Repository → Base de Datos
                                    ↓
                                  GORM
                                    ↓
                                  MySQL
```

---

## Desglose Línea por Línea

### Package e Imports

```go
package repository
```
- Define el paquete `repository`
- Todos los repositorios de tu app irán aquí

```go
import (
	"github.com/gonzalo-wi/cellcontrol/internal/domain"
	"gorm.io/gorm"
)
```

**Imports explicados:**

1. **`github.com/gonzalo-wi/cellcontrol/internal/domain`**
    - Importa tus **modelos de dominio** (estructuras de datos)
    - Contiene `domain.User` - la definición de cómo es un usuario
    - Es TU código, no una librería externa

2. **`gorm.io/gorm`**
    - **GORM**: ORM (Object-Relational Mapper) para Go
    - Traduce código Go → SQL automáticamente
    - Maneja conexiones, queries, migraciones, etc.
    - Muy popular en el ecosistema Go

---

## Interfaz UserRepository

```go
type UserRepository interface {
	CreateUser(user []domain.User) error
	GetAllUsers() ([]domain.User, error)
}
```

### ¿Qué es una Interface en Go?

Una **interface** define un **contrato** - qué métodos debe tener algo, sin definir cómo los implementa.

**¿Por qué usar una interface aquí?**

1. **Abstracción**: El código que usa esto no sabe si usa MySQL, PostgreSQL, o memoria
2. **Testing**: Puedes crear un mock que implemente esta interface
3. **Flexibilidad**: Puedes cambiar la implementación sin cambiar quien la usa
4. **Dependency Injection**: Facilita inyectar dependencias

### Métodos de la Interface

#### 1. CreateUser
```go
CreateUser(user []domain.User) error
```

**Propósito**: Crear uno o más usuarios en la base de datos

**Parámetros:**
- `user []domain.User`: Un **slice** (lista) de usuarios a crear
- ¿Por qué slice? Permite crear múltiples usuarios en una sola operación (más eficiente)

**Retorno:**
- `error`:
    - `nil` si todo fue bien
    - Un error si algo falló (conexión perdida, constraint violado, etc.)

**Ejemplo de uso:**
```go
users := []domain.User{
	{Name: "Juan", Email: "juan@example.com"},
	{Name: "María", Email: "maria@example.com"},
}
err := repo.CreateUser(users)
if err != nil {
	log.Fatal("Error creando usuarios:", err)
}
```

#### 2. GetAllUsers
```go
GetAllUsers() ([]domain.User, error)
```

**Propósito**: Obtener TODOS los usuarios de la base de datos

**Parámetros:**
- Ninguno (no necesita parámetros porque trae TODO)

**Retorno (2 valores):**
1. `[]domain.User`: Slice con todos los usuarios encontrados
2. `error`: Error si algo falló, o `nil` si todo bien

**Ejemplo de uso:**
```go
users, err := repo.GetAllUsers()
if err != nil {
	log.Fatal("Error obteniendo usuarios:", err)
}
for _, user := range users {
	fmt.Printf("Usuario: %s (%s)\n", user.Name, user.Email)
}
```

---

## Implementación Concreta

### Estructura userRepository (privada)

```go
type userRepository struct {
	db *gorm.DB
}
```

**Características importantes:**

1. **Nombre en minúscula**: `userRepository` vs `UserRepository`
    - En Go: minúscula = **privado** (no exportado)
    - Solo accesible dentro del paquete `repository`
    - Nadie de afuera puede crear esto directamente

2. **Campo `db`**:
    - `*gorm.DB`: Puntero a la conexión de base de datos de GORM
    - Este es el "motor" que ejecutará las queries SQL
    - Se pasa al crear el repositorio

**¿Por qué tener dos tipos (interface + struct)?**
```
UserRepository (interface pública)  ← Lo que exportas
        ↑
        |
        | implementa
        |
userRepository (struct privada)    ← Implementación interna
```

Ventajas:
- Los usuarios del paquete trabajan con la **interface**
- Tú puedes cambiar la **implementación** sin afectar a nadie
- Puedes tener **múltiples implementaciones** de la misma interface

---

### Constructor: NewUserRepository

```go
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
```

**Propósito**: Crear una nueva instancia del repositorio (patrón Factory)

**¿Cómo funciona?**

1. **Recibe `db *gorm.DB`**:
    - La conexión a la base de datos
    - Ya debe estar inicializada
    - Se pasa desde `main.go` o donde inicialices la app

2. **Crea `&userRepository{db: db}`**:
    - `&` → toma la dirección (crea un puntero)
    - `userRepository{...}` → crea la struct
    - `db: db` → asigna el campo db con el parámetro recibido

3. **Retorna `UserRepository` (la interface)**:
    - Aunque crea un `*userRepository` (puntero a struct)
    - Lo retorna como `UserRepository` (interface)
    - Esto funciona porque `*userRepository` implementa `UserRepository`

**Ejemplo de uso:**
```go
// En tu main.go o inicialización
db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
userRepo := repository.NewUserRepository(db)

// userRepo es de tipo UserRepository (interface)
// pero internamente es un *userRepository (struct)
```

**Ventaja de este patrón:**
- Encapsulación: No expones la implementación
- Control: Solo tú decides cómo crear repositorios
- Consistencia: Siempre se crean de la misma forma

---

## Implementación de Métodos

### 1. CreateUser - Crear Usuarios

```go
func (r *userRepository) CreateUser(user []domain.User) error {
	return r.db.Create(user).Error
}
```

**Desglose completo:**

#### Firma del método
```go
func (r *userRepository) CreateUser(user []domain.User) error
```

- `(r *userRepository)`: **Receptor** del método
    - `r` es el nombre de la variable (convención: primera letra del tipo)
    - `*userRepository` es el tipo (puntero a struct)
    - Significa: "este método pertenece a `*userRepository`"
    - Dentro del método, accedes a la struct con `r`

- `CreateUser`: Nombre del método (debe coincidir con la interface)

- `(user []domain.User)`: Parámetro - slice de usuarios

- `error`: Tipo de retorno

#### Implementación
```go
return r.db.Create(user).Error
```

**Paso a paso:**

1. **`r.db`**: Accede al campo `db` de la struct (la conexión GORM)

2. **`.Create(user)`**: Método de GORM
    - Ejecuta un `INSERT INTO users ...` en SQL
    - Automáticamente convierte `[]domain.User` a filas SQL
    - Devuelve `*gorm.DB` (para encadenar métodos)

3. **`.Error`**: Campo de `*gorm.DB`
    - Contiene el último error ocurrido
    - `nil` si todo bien
    - Error específico si algo falló

4. **`return ...`**: Devuelve el error directamente

**SQL Generado (aproximadamente):**
```sql
INSERT INTO users (name, email, phone, created_at, updated_at) 
VALUES 
  ('Juan', 'juan@example.com', '123456', NOW(), NOW()),
  ('María', 'maria@example.com', '789012', NOW(), NOW());
```

**¿Qué puede salir mal?**
- Conexión a BD perdida → error
- Email duplicado (si hay constraint UNIQUE) → error
- Campos requeridos faltantes → error
- Timeout de query → error

---

### 2. GetAllUsers - Obtener Todos los Usuarios

```go
func (r *userRepository) GetAllUsers() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	return users, err
}
```

**Desglose completo:**

#### Firma del método
```go
func (r *userRepository) GetAllUsers() ([]domain.User, error)
```

- `(r *userRepository)`: Receptor - método de `*userRepository`
- `GetAllUsers()`: Sin parámetros (trae TODO)
- `([]domain.User, error)`: Retorna 2 valores (patrón común en Go)

#### Implementación línea por línea

**Línea 1:**
```go
var users []domain.User
```
- Declara una variable `users` de tipo slice de `domain.User`
- Valor inicial: `nil` (slice vacío)
- GORM llenará este slice con los resultados

**Línea 2:**
```go
err := r.db.Find(&users).Error
```

Desglose:
1. **`r.db`**: Conexión GORM
2. **`.Find(&users)`**: Método de GORM
    - `Find`: Equivalente a `SELECT * FROM users`
    - `&users`: Pasa la **dirección** del slice
    - ¿Por qué `&`? GORM necesita modificar el slice original
    - GORM llena `users` con los resultados automáticamente
3. **`.Error`**: Obtiene el error (si hubo)
4. **`err :=`**: Asigna el error a la variable `err`

**Línea 3:**
```go
return users, err
```
- Retorna ambos valores
- `users`: Puede ser vacío (`[]`) o contener usuarios
- `err`: Puede ser `nil` (éxito) o un error

**SQL Generado:**
```sql
SELECT * FROM users;
```

**Comportamiento según casos:**

| Caso | `users` | `err` |
|------|---------|-------|
| 3 usuarios en BD | `[{...}, {...}, {...}]` | `nil` |
| 0 usuarios en BD | `[]` (slice vacío) | `nil` |
| Error de conexión | `[]` o `nil` | Error |
| Tabla no existe | `[]` o `nil` | Error |

---

## Flujo Completo de Uso

### Ejemplo Real: Crear y Listar Usuarios

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/gonzalo-wi/cellcontrol/internal/domain"
	"github.com/gonzalo-wi/cellcontrol/internal/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. Conectar a la base de datos
	dsn := "root:password@tcp(localhost:3306)/cellcontrol?charset=utf8mb4&parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error conectando a BD:", err)
	}
	
	// 2. Crear el repositorio
	userRepo := repository.NewUserRepository(db)
	
	// 3. Crear usuarios
	newUsers := []domain.User{
		{Name: "Juan Pérez", Email: "juan@example.com", Phone: "123456789"},
		{Name: "María García", Email: "maria@example.com", Phone: "987654321"},
	}
	
	err = userRepo.CreateUser(newUsers)
	if err != nil {
		log.Fatal("Error creando usuarios:", err)
	}
	fmt.Println("✅ Usuarios creados exitosamente")
	
	// 4. Obtener todos los usuarios
	users, err := userRepo.GetAllUsers()
	if err != nil {
		log.Fatal("Error obteniendo usuarios:", err)
	}
	
	// 5. Mostrar usuarios
	fmt.Printf("\n📋 Total de usuarios: %d\n", len(users))
	for i, user := range users {
		fmt.Printf("%d. %s (%s) - Tel: %s\n", 
			i+1, user.Name, user.Email, user.Phone)
	}
}
```

**Salida esperada:**
```
✅ Usuarios creados exitosamente

📋 Total de usuarios: 2
1. Juan Pérez (juan@example.com) - Tel: 123456789
2. María García (maria@example.com) - Tel: 987654321
```

---

## Arquitectura: Cómo Encaja en tu Proyecto

```
┌─────────────────────────────────────────────────┐
│                   HANDLER                       │
│  (user_handler.go)                             │
│  - Recibe HTTP requests                        │
│  - Valida input                                │
│  - Llama al Service                            │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│                   SERVICE                       │
│  (user_service.go)                             │
│  - Lógica de negocio                           │
│  - Validaciones complejas                      │
│  - Llama al Repository                         │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│                 REPOSITORY   ◄── ESTÁS AQUÍ    │
│  (user_repository.go)                          │
│  - Acceso a datos                              │
│  - Queries SQL (via GORM)                      │
│  - CRUD operations                             │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│                BASE DE DATOS                    │
│  (MySQL)                                       │
│  - Tabla: users                                │
│  - Almacenamiento persistente                  │
└─────────────────────────────────────────────────┘
```

---

## Ventajas de esta Implementación

### ✅ Separación de Responsabilidades
- El repository SOLO se encarga de datos
- No tiene lógica de negocio
- No sabe nada de HTTP

### ✅ Testeable
```go
// Mock para testing
type mockUserRepository struct{}

func (m *mockUserRepository) CreateUser(user []domain.User) error {
	return nil // Simula éxito
}

func (m *mockUserRepository) GetAllUsers() ([]domain.User, error) {
	return []domain.User{{Name: "Test"}}, nil
}
```

### ✅ Mantenible
- Si cambias de MySQL a PostgreSQL: solo cambias el repository
- Si cambias de GORM a SQL puro: solo cambias el repository
- El resto del código NO se entera

### ✅ Reutilizable
```go
// Puedes usar el mismo repositorio en múltiples servicios
userService := service.NewUserService(userRepo)
authService := service.NewAuthService(userRepo)
```

---

## Limitaciones y Mejoras Posibles

### ⚠️ Limitaciones Actuales

1. **GetAllUsers sin paginación**
    - Si tienes 1 millón de usuarios, los trae TODOS
    - Puede colapsar la memoria

2. **No hay filtros**
    - No puedes buscar por email, nombre, etc.

3. **No hay actualización ni eliminación**
    - Solo CREATE y READ, faltan UPDATE y DELETE

4. **CreateUser acepta slice pero usualmente creas de 1 en 1**
    - Más común: `CreateUser(user domain.User)` (singular)

### 🚀 Mejoras Sugeridas

#### 1. Añadir más operaciones CRUD

```go
type UserRepository interface {
	CreateUser(user *domain.User) error           // Singular
	CreateUsers(users []domain.User) error        // Plural
	GetAllUsers(page, pageSize int) ([]domain.User, error)  // Con paginación
	GetUserByID(id uint) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id uint) error
}
```

#### 2. Paginación en GetAllUsers

```go
func (r *userRepository) GetAllUsers(page, pageSize int) ([]domain.User, error) {
	var users []domain.User
	offset := (page - 1) * pageSize
	err := r.db.Limit(pageSize).Offset(offset).Find(&users).Error
	return users, err
}
```

#### 3. Buscar por ID

```go
func (r *userRepository) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
```

#### 4. Actualizar usuario

```go
func (r *userRepository) UpdateUser(user *domain.User) error {
	return r.db.Save(user).Error
}
```

#### 5. Eliminar usuario

```go
func (r *userRepository) DeleteUser(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}
```

---

## Conceptos Clave de GORM Usados

| Método GORM | SQL Equivalente | Uso en Código |
|-------------|-----------------|---------------|
| `Create(user)` | `INSERT INTO users VALUES (...)` | Crear registros |
| `Find(&users)` | `SELECT * FROM users` | Obtener todos |
| `First(&user, id)` | `SELECT * FROM users WHERE id = ? LIMIT 1` | Obtener uno |
| `Save(user)` | `UPDATE users SET ... WHERE id = ?` | Actualizar |
| `Delete(&user, id)` | `DELETE FROM users WHERE id = ?` | Eliminar |

---

## Resumen Ejecutivo

**¿Qué hace este código?**
Define una capa de acceso a datos para usuarios usando el patrón Repository.

**¿Por qué está diseñado así?**
- Interface pública → flexibilidad y testing
- Struct privada → encapsulación
- GORM → abstracción de SQL

**¿Cuándo se usa?**
Desde el Service layer cuando necesitas persistir o recuperar usuarios.

**¿Qué falta?**
- Update y Delete
- Búsquedas por filtros
- Paginación
- Manejo de transacciones

---

¿Necesitas que profundice en alguna parte específica o que implemente alguna de las mejoras sugeridas?
