# Explicación Detallada: HTTP Server

## Visión General
Este archivo implementa el **servidor HTTP** de tu aplicación usando el framework **Gin**. Es el punto de entrada de todas las peticiones HTTP y se encarga de configurar rutas, middlewares y arrancar el servidor web.

---

## ¿Qué es Gin?

**Gin** es un framework web HTTP para Go, similar a Express.js en Node.js o Flask en Python.

**Características:**
- ✅ Router rápido y eficiente
- ✅ Middleware incorporado
- ✅ Validación de JSON automática
- ✅ Manejo de errores
- ✅ Muy popular en el ecosistema Go

---

## Desglose Línea por Línea

### Package e Imports

```go
package http
```
- Define el paquete `http` (nota: no confundir con el paquete estándar `net/http`)
- Contiene todo lo relacionado con el servidor HTTP

```go
import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gonzalo-wi/cellcontrol/internal/config"
	"github.com/gonzalo-wi/cellcontrol/internal/http/handlers"
	"github.com/gonzalo-wi/cellcontrol/pkg/logger"
)
```

**Imports explicados:**

1. **`"fmt"`** - Paquete estándar de Go
    - Usado para formatear strings
    - En este caso: `fmt.Sprintf()` para construir la dirección del servidor

2. **`"github.com/gin-gonic/gin"`** - Framework web
    - El motor HTTP principal
    - Proporciona routing, middleware, JSON handling, etc.

3. **`"github.com/gonzalo-wi/cellcontrol/internal/config"`** - Tu paquete de configuración
    - Contiene `config.Config` con la configuración de la app
    - Necesario para obtener el puerto HTTP

4. **`"github.com/gonzalo-wi/cellcontrol/internal/http/handlers"`** - Tus handlers
    - Contiene los handlers (controladores) HTTP
    - En este caso: `handlers.UserHandler`

5. **`"github.com/gonzalo-wi/cellcontrol/pkg/logger"`** - Tu logger personalizado
    - Para registrar mensajes informativos
    - Usado para mostrar en qué puerto escucha el servidor

---

## Estructura Server

```go
type Server struct {
	engine *gin.Engine
	cfg    *config.Config
}
```

### Análisis de Campos

**1. `engine *gin.Engine`**
- Tipo: Puntero al motor de Gin
- **¿Qué es `gin.Engine`?**
    - El "corazón" del servidor HTTP
    - Maneja routing (rutas)
    - Ejecuta middlewares
    - Procesa requests y genera responses
- **¿Por qué es un puntero?**
    - Es un objeto grande
    - Se pasa por referencia para eficiencia
    - Gin siempre trabaja con punteros

**2. `cfg *config.Config`**
- Tipo: Puntero a la configuración
- Contiene: puerto, DSN de BD, entorno, etc.
- **¿Para qué se guarda?**
    - Para acceder al puerto cuando se arranca el servidor
    - Podría usarse para configuración adicional (timeouts, SSL, etc.)

**¿Por qué encapsular en una struct?**
```go
// Sin struct (mal):
func RunServer(engine *gin.Engine, cfg *config.Config) { ... }

// Con struct (bien):
server := NewServer(cfg, userHandler)
server.Run()
```

Ventajas:
- ✅ Encapsulación: agrupa datos relacionados
- ✅ Métodos: puedes añadir `Shutdown()`, `Restart()`, etc.
- ✅ Estado: mantiene configuración y engine juntos
- ✅ Testing: más fácil de mockear

---

## Constructor: NewServer

```go
func NewServer(cfg *config.Config, userHandler *handlers.UserHandler) *Server {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	userHandler.RegisterRoutes(api)

	return &Server{
		engine: r,
		cfg:    cfg,
	}
}
```

### Análisis Completo

#### Firma del Constructor

```go
func NewServer(cfg *config.Config, userHandler *handlers.UserHandler) *Server
```

**Parámetros (Dependency Injection):**

1. `cfg *config.Config`
    - Configuración de la aplicación
    - Contiene el puerto HTTP, etc.

2. `userHandler *handlers.UserHandler`
    - Handler de usuarios ya inicializado
    - Inyectado desde afuera (no se crea aquí)

**Retorno:**
- `*Server`: Puntero a la nueva instancia del servidor

**Patrón:** Constructor con inyección de dependencias

---

#### Paso 1: Crear el Engine de Gin

```go
r := gin.Default()
```

**`gin.Default()` crea un nuevo engine con middleware predeterminado:**

1. **Logger middleware**: Registra todas las requests
   ```
   [GIN] 2025/11/23 - 14:30:45 | 200 | 2.456789ms | 127.0.0.1 | POST /api/v1/usuarios
   ```

2. **Recovery middleware**: Recupera de panics
    - Si hay un panic → lo captura
    - Retorna HTTP 500 en lugar de crashear
    - Muy importante en producción

**Alternativa:**
```go
r := gin.New()  // Sin middleware
r.Use(gin.Logger())
r.Use(gin.Recovery())
// Mismo resultado que gin.Default()
```

**Variable `r`:**
- Convención común para el router/engine
- Abreviatura de "router"

---

#### Paso 2: Health Check Endpoint

```go
r.GET("/health", func(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
})
```

**Desglose completo:**

**`r.GET(...)`**
- Registra una ruta para peticiones HTTP GET
- Otros métodos: `POST`, `PUT`, `DELETE`, `PATCH`

**`"/health"`**
- Ruta del endpoint
- URL completa: `http://localhost:8080/health`

**`func(c *gin.Context) { ... }`**
- **Handler inline** (función anónima)
- `c *gin.Context`: Contexto de la petición
    - Contiene: request, response, headers, params, etc.
    - Es el objeto principal en Gin

**`c.JSON(200, gin.H{"status": "ok"})`**
- `c.JSON()`: Envía respuesta JSON
- `200`: Código de estado HTTP (OK)
- `gin.H{...}`: Atajo para `map[string]any`
    - `gin.H{"status": "ok"}` = `{"status": "ok"}`

**¿Para qué sirve /health?**

Este endpoint es un **health check** (verificación de salud):

```bash
# Prueba si el servidor está vivo
curl http://localhost:8080/health
# Respuesta: {"status":"ok"}
```

**Casos de uso:**
- ✅ **Monitoreo**: Herramientas como Kubernetes, AWS ELB lo usan
- ✅ **Load balancers**: Verifican si el servidor responde
- ✅ **Testing**: Verificar que el servidor arrancó
- ✅ **CI/CD**: Verificar deployment exitoso

**Mejora común:**
```go
r.GET("/health", func(c *gin.Context) {
	// Verificar conexión a BD
	if err := db.Ping(); err != nil {
		c.JSON(503, gin.H{"status": "unhealthy", "error": "database down"})
		return
	}
	c.JSON(200, gin.H{
		"status": "ok",
		"version": "1.0.0",
		"uptime": time.Since(startTime).String(),
	})
})
```

---

#### Paso 3: Crear Grupo de Rutas API

```go
api := r.Group("/api/v1")
```

**¿Qué es un Route Group?**

Un **grupo de rutas** aplica un **prefijo común** a todas las rutas dentro de él.

**Sin grupos:**
```go
r.POST("/api/v1/usuarios", handler.CreateUser)
r.GET("/api/v1/usuarios", handler.ListUsers)
r.POST("/api/v1/celulares", handler.CreateCellphone)
// Repetitivo: /api/v1 en cada ruta
```

**Con grupos:**
```go
api := r.Group("/api/v1")
api.POST("/usuarios", handler.CreateUser)        // → /api/v1/usuarios
api.GET("/usuarios", handler.ListUsers)          // → /api/v1/usuarios
api.POST("/celulares", handler.CreateCellphone)  // → /api/v1/celulares
```

**Ventajas:**
- ✅ DRY (Don't Repeat Yourself)
- ✅ Versionado de API (v1, v2, etc.)
- ✅ Middleware específico por grupo
- ✅ Cambio de prefijo en un solo lugar

**Estructura de URLs resultante:**
```
http://localhost:8080/health           ← Fuera del grupo
http://localhost:8080/api/v1/usuarios  ← Dentro del grupo
http://localhost:8080/api/v1/celulares ← Dentro del grupo
```

**Middleware en grupos:**
```go
api := r.Group("/api/v1")
api.Use(AuthMiddleware())  // Solo aplica a rutas de /api/v1

admin := api.Group("/admin")
admin.Use(AdminMiddleware()) // Solo aplica a /api/v1/admin/*
```

---

#### Paso 4: Registrar Rutas del Handler

```go
userHandler.RegisterRoutes(api)
```

**¿Qué hace esto?**

Llama al método `RegisterRoutes` del `UserHandler`, pasándole el grupo `api`.

**En `user_handler.go`:**
```go
func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/usuarios", h.CreateUser)
	r.GET("/usuarios", h.ListUsers)
}
```

**Resultado:**
- `POST /api/v1/usuarios` → `h.CreateUser`
- `GET /api/v1/usuarios` → `h.ListUsers`

**¿Por qué este diseño?**

| Enfoque | Código |
|---------|--------|
| **Malo (acoplado)** | `r.POST("/usuarios", handler.CreateUser)` directamente aquí |
| **Bueno (desacoplado)** | `handler.RegisterRoutes(api)` |

**Ventajas del enfoque actual:**
- ✅ **Separación de responsabilidades**: Cada handler define sus propias rutas
- ✅ **Escalabilidad**: Fácil añadir más handlers
- ✅ **Mantenibilidad**: Rutas junto al código que las maneja
- ✅ **Testing**: Puedes testear el handler independientemente

**Escalabilidad:**
```go
// Fácil añadir más recursos
userHandler.RegisterRoutes(api)
cellphoneHandler.RegisterRoutes(api)
historyHandler.RegisterRoutes(api)
authHandler.RegisterRoutes(api)
```

---

#### Paso 5: Retornar el Server

```go
return &Server{
	engine: r,
	cfg:    cfg,
}
```

**Crea y retorna:**
- Nuevo `Server` struct
- Campo `engine`: el motor de Gin configurado
- Campo `cfg`: la configuración recibida
- `&` → retorna un puntero

**Estado del servidor en este punto:**
```
Server {
  engine: *gin.Engine {
    routes: [
      GET /health
      POST /api/v1/usuarios
      GET /api/v1/usuarios
    ]
    middlewares: [Logger, Recovery]
  }
  cfg: *config.Config {
    HttpPort: "8080"
    DatabaseDSN: "..."
  }
}
```

**Importante:** El servidor aún **NO está corriendo**, solo está **configurado**.

---

## Método Run: Arrancar el Servidor

```go
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%s", s.cfg.HttpPort)
	logger.Info("escuchando en %s", addr)
	return s.engine.Run(addr)
}
```

### Análisis Línea por Línea

#### Línea 1: Firma del Método

```go
func (s *Server) Run() error
```

- **Receptor**: `(s *Server)` - método de Server
- **Nombre**: `Run` - arranca el servidor
- **Retorno**: `error` - porque puede fallar (puerto ocupado, permisos, etc.)

#### Línea 2: Construir Dirección

```go
addr := fmt.Sprintf(":%s", s.cfg.HttpPort)
```

**¿Qué hace?**
- `fmt.Sprintf()`: Formatea un string
- `":%s"`: Template con placeholder
- `s.cfg.HttpPort`: El puerto de la configuración (ej: "8080")

**Ejemplos:**
```go
// Si HttpPort = "8080"
addr = ":8080"

// Si HttpPort = "3000"
addr = ":3000"
```

**Formato de dirección en Go:**
- `:8080` → Escucha en todas las interfaces en el puerto 8080
- `localhost:8080` → Solo escucha en localhost
- `192.168.1.10:8080` → Solo en esa IP específica

**Equivalentes:**
- `:8080` = `0.0.0.0:8080` = escucha en todas las interfaces de red

#### Línea 3: Log Informativo

```go
logger.Info("escuchando en %s", addr)
```

**Registra un mensaje informativo:**
```
INFO: 2025/11/23 14:30:45 server.go:36 escuchando en :8080
```

**¿Por qué es importante?**
- ✅ Confirma que el servidor arrancó
- ✅ Muestra en qué puerto (útil si es configurable)
- ✅ Ayuda en debugging
- ✅ Se ve en logs de producción

#### Línea 4: Arrancar el Engine

```go
return s.engine.Run(addr)
```

**¿Qué hace `s.engine.Run(addr)`?**

1. **Inicia un servidor HTTP** en la dirección especificada
2. **Bloquea** el hilo de ejecución (no retorna hasta que el servidor se detenga)
3. **Escucha peticiones** entrantes
4. **Procesa requests** según las rutas configuradas
5. **Retorna error** si falla al arrancar

**Este método es bloqueante:**
```go
server.Run()  // ← El programa se queda aquí
fmt.Println("Esto nunca se ejecuta (a menos que el servidor se detenga)")
```

**Internamente hace:**
```go
// Simplificación de lo que hace Gin
http.ListenAndServe(addr, s.engine)
```

**Errores comunes:**
```go
// Error: Puerto ocupado
panic: listen tcp :8080: bind: address already in use

// Error: Puerto inválido (< 1024 sin permisos)
panic: listen tcp :80: bind: permission denied
```

---

## Flujo Completo de Inicialización

### Desde main.go (Ejemplo)

```go
package main

import (
	"log"
	
	"github.com/gonzalo-wi/cellcontrol/internal/config"
	"github.com/gonzalo-wi/cellcontrol/internal/http"
	"github.com/gonzalo-wi/cellcontrol/internal/http/handlers"
	"github.com/gonzalo-wi/cellcontrol/internal/repository"
	"github.com/gonzalo-wi/cellcontrol/internal/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. Cargar configuración
	cfg := config.MustLoad()
	
	// 2. Conectar a base de datos
	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Error conectando a BD:", err)
	}
	
	// 3. Crear repositorio
	userRepo := repository.NewUserRepository(db)
	
	// 4. Crear servicio (inyecta repo)
	userService := service.NewUserService(userRepo)
	
	// 5. Crear handler (inyecta servicio)
	userHandler := handlers.NewUserHandler(userService)
	
	// 6. Crear servidor (inyecta config y handlers)
	server := http.NewServer(cfg, userHandler)
	
	// 7. Arrancar servidor (BLOQUEA AQUÍ)
	if err := server.Run(); err != nil {
		log.Fatal("Error arrancando servidor:", err)
	}
}
```

**Salida esperada:**
```
[config] env=dev port=8080 db=root:***@tcp(localhost:3306)/cellcontrol
INFO: 2025/11/23 14:30:45 server.go:36 escuchando en :8080
[GIN-debug] Listening and serving HTTP on :8080
```

---

## Diagrama de Arquitectura

```
┌─────────────────────────────────────────────────┐
│              MAIN.GO                            │
│  - Inicializa todo                             │
│  - Inyecta dependencias                        │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│         HTTP SERVER (server.go)  ◄── ESTÁS AQUÍ│
│  - Configura Gin                               │
│  - Define rutas                                │
│  - Arranca servidor HTTP                       │
└──────┬─────────────────────────────────┬────────┘
       │                                 │
       │ /health                         │ /api/v1/*
       ▼                                 ▼
┌─────────────┐              ┌──────────────────────┐
│   Inline    │              │   USER HANDLER       │
│   Handler   │              │  - CreateUser        │
│             │              │  - ListUsers         │
└─────────────┘              └──────┬───────────────┘
                                    │
                                    ▼
                        ┌──────────────────────┐
                        │    USER SERVICE      │
                        │  - Lógica negocio    │
                        └──────┬───────────────┘
                               │
                               ▼
                   ┌──────────────────────┐
                   │   USER REPOSITORY    │
                   │  - Acceso a datos    │
                   └──────┬───────────────┘
                          │
                          ▼
                   ┌──────────────────────┐
                   │      MYSQL           │
                   └──────────────────────┘
```

---

## Petición HTTP Completa

### Request
```bash
POST http://localhost:8080/api/v1/usuarios
Content-Type: application/json

{
  "nombre": "Juan",
  "apellido": "Pérez",
  "email": "juan@example.com",
  "reparto": "IT"
}
```

### Flujo Interno

```
1. Cliente → :8080
           ↓
2. Gin Engine (server.go)
   - Gin Logger middleware (registra)
   - Gin Recovery middleware (captura panics)
           ↓
3. Router match: POST /api/v1/usuarios
           ↓
4. UserHandler.CreateUser (user_handler.go)
   - Valida JSON
   - Extrae campos
           ↓
5. UserService.CreateUser (user_service.go)
   - Limpia datos
   - Normaliza email
           ↓
6. UserRepository.CreateUser (user_repository.go)
   - INSERT INTO users
           ↓
7. MySQL
   - Guarda registro
   - Retorna ID
           ↓
8. Respuesta sube por las capas
           ↓
9. Cliente ← {"message": "usuario creado exitosamente"}
```

### Log en Consola

```
[GIN] 2025/11/23 - 14:30:45 | 201 |  12.456789ms |  127.0.0.1 | POST     /api/v1/usuarios
INFO: 2025/11/23 14:30:45 server.go:36 escuchando en :8080
```

---

## Conceptos Clave de Gin

| Método/Tipo | Propósito | Ejemplo |
|-------------|-----------|---------|
| `gin.Default()` | Crea engine con middleware | `r := gin.Default()` |
| `gin.New()` | Crea engine sin middleware | `r := gin.New()` |
| `r.GET()` | Registra ruta GET | `r.GET("/path", handler)` |
| `r.POST()` | Registra ruta POST | `r.POST("/path", handler)` |
| `r.Group()` | Crea grupo de rutas | `api := r.Group("/api")` |
| `gin.Context` | Contexto de request | `func(c *gin.Context)` |
| `c.JSON()` | Responder con JSON | `c.JSON(200, data)` |
| `gin.H` | Map para JSON rápido | `gin.H{"key": "value"}` |
| `r.Run()` | Arranca el servidor | `r.Run(":8080")` |

---

## Mejoras Posibles

### 🚀 1. Graceful Shutdown

```go
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%s", s.cfg.HttpPort)
	logger.Info("escuchando en %s", addr)
	
	srv := &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}
	
	// Goroutine para arrancar servidor
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error en servidor: %v", err)
		}
	}()
	
	// Esperar señal de terminación
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	logger.Info("Apagando servidor...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return srv.Shutdown(ctx)
}
```

### 🚀 2. CORS Middleware

```go
func NewServer(cfg *config.Config, userHandler *handlers.UserHandler) *Server {
	r := gin.Default()
	
	// CORS para permitir peticiones desde frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	
	// ... resto del código
}
```

### 🚀 3. Rate Limiting

```go
func NewServer(cfg *config.Config, userHandler *handlers.UserHandler) *Server {
	r := gin.Default()
	
	// Limitar a 100 requests por minuto
	r.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))
	
	// ... resto del código
}
```

### 🚀 4. Múltiples Handlers

```go
func NewServer(
	cfg *config.Config,
	userHandler *handlers.UserHandler,
	cellphoneHandler *handlers.CellphoneHandler,
	historyHandler *handlers.HistoryHandler,
) *Server {
	r := gin.Default()
	
	r.GET("/health", healthCheck)
	
	api := r.Group("/api/v1")
	userHandler.RegisterRoutes(api)
	cellphoneHandler.RegisterRoutes(api)
	historyHandler.RegisterRoutes(api)
	
	return &Server{engine: r, cfg: cfg}
}
```

### 🚀 5. Modo de Entorno

```go
func NewServer(cfg *config.Config, userHandler *handlers.UserHandler) *Server {
	// Configurar modo según entorno
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	
	r := gin.Default()
	// ... resto
}
```

---

## Testing del Server

```go
package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/gonzalo-wi/cellcontrol/internal/config"
	httpserver "github.com/gonzalo-wi/cellcontrol/internal/http"
	"github.com/gonzalo-wi/cellcontrol/internal/http/handlers"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	// Setup
	cfg := &config.Config{HttpPort: "8080"}
	mockUserHandler := &handlers.UserHandler{}
	server := httpserver.NewServer(cfg, mockUserHandler)
	
	// Crear request de prueba
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	
	// Ejecutar
	server.engine.ServeHTTP(w, req)
	
	// Verificar
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
}
```

---

## Resumen Ejecutivo

**¿Qué hace este código?**
Configura y arranca un servidor HTTP usando Gin con rutas versionadas y health check.

**¿Por qué este diseño?**
- ✅ Encapsulación en struct Server
- ✅ Dependency injection de handlers
- ✅ Versionado de API con grupos
- ✅ Separación entre configuración y ejecución

**¿Cuándo se usa?**
En `main.go` para inicializar el servidor web de la aplicación.

**Endpoints disponibles:**
- `GET /health` → Health check
- `POST /api/v1/usuarios` → Crear usuario
- `GET /api/v1/usuarios` → Listar usuarios

---

¿Necesitas que profundice en alguna parte específica o que implemente alguna mejora?
