# Explicación Detallada: Config (Configuración)

## Visión General
Este paquete gestiona la **configuración de la aplicación**, cargando valores desde **variables de entorno** y un archivo **`.env`**. Es el centro de configuración que define cómo se comporta tu aplicación en diferentes entornos (desarrollo, producción, etc.).

---

## Desglose Línea por Línea

### Package e Imports
```go
package config
```
- Define el paquete `config` para gestionar la configuración

```go
import (
	"log"
	"os"
	"github.com/joho/godotenv"
)
```
- **`log`**: Para imprimir mensajes de log
- **`os`**: Para acceder a variables de entorno del sistema operativo
- **`godotenv`**: Librería externa que carga variables desde un archivo `.env`

**¿Qué es godotenv?**
- Lee archivos `.env` y carga las variables en `os.Getenv()`
- Muy útil para desarrollo local (evitas hardcodear credenciales)
- En producción, usarías variables de entorno reales del sistema

---

## Estructura Config

```go
type Config struct {
	Env         string
	HttpPort    string
	DatabaseDSN string
}
```

**Define la configuración de tu aplicación:**

1. **`Env`**: Entorno de ejecución
   - Ejemplos: `"dev"`, `"prod"`, `"test"`
   - Útil para cambiar comportamientos según el entorno

2. **`HttpPort`**: Puerto donde corre el servidor HTTP
   - Ejemplo: `"8080"`, `"3000"`
   - Permite cambiar el puerto sin recompilar

3. **`DatabaseDSN`**: Data Source Name (cadena de conexión a la BD)
   - Formato MySQL: `"usuario:contraseña@tcp(host:puerto)/nombre_bd"`
   - Ejemplo: `"root:secret@tcp(localhost:3306)/cellcontrol"`

**DSN significa "Data Source Name"** - contiene toda la info para conectarse a la base de datos.

---

## Función Load()

```go
func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		Env:         getEnv("APP_ENV", "dev"),
		HttpPort:    getEnv("HTTP_PORT", "8080"),
		DatabaseDSN: getEnv("DATABASE_DSN", "user:password@tcp(localhost:3306)/dbname"),
	}
}
```

**Propósito**: Cargar la configuración desde variables de entorno

**Paso a paso:**

1. **`_ = godotenv.Load()`**
   - Intenta cargar el archivo `.env` en la raíz del proyecto
   - El `_` significa que **ignoramos el error** (si no existe `.env`, no pasa nada)
   - ¿Por qué ignorar? En producción quizás no tengas archivo `.env`

2. **`return &Config{...}`**
   - Crea una nueva instancia de `Config`
   - El `&` devuelve un **puntero** a la estructura
   - Rellena cada campo con `getEnv()`

3. **Para cada campo:**
   - `getEnv("NOMBRE_VARIABLE", "valor_por_defecto")`
   - Busca la variable de entorno
   - Si no existe, usa el valor por defecto

**Ejemplo de archivo `.env`:**
```env
APP_ENV=production
HTTP_PORT=8080
DATABASE_DSN=root:mypassword@tcp(localhost:3306)/cellcontrol
```

---

## Función MustLoad()

```go
func MustLoad() *Config {
	cfg := Load()
	log.Printf("[config] env=%s port=%s db=%s\n", cfg.Env, cfg.HttpPort, cfg.DatabaseDSN)
	return cfg
}
```

**Propósito**: Cargar configuración y mostrar un log de confirmación

**¿Por qué "Must"?**
- Convención en Go: funciones con `Must` suelen **paniquear** si algo falla
- Aquí NO panicea, pero el nombre sugiere que **debes** usar esta para inicialización
- Es la versión "verbosa" de `Load()` - muestra lo que cargó

**Salida esperada:**
```
[config] env=production port=8080 db=root:***@tcp(localhost:3306)/cellcontrol
```

**Uso típico:**
```go
// Al iniciar la aplicación
cfg := config.MustLoad()  // Carga Y muestra la config
```

---

## Función getEnv() - La Magia

```go
func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
```

**Propósito**: Obtener una variable de entorno o devolver un valor por defecto

**Parámetros:**
- `key`: Nombre de la variable de entorno (ej: `"HTTP_PORT"`)
- `def`: Valor por defecto si no existe (ej: `"8080"`)

**¿Cómo funciona?**

1. **`os.LookupEnv(key)`** devuelve 2 valores:
   - `val`: El valor de la variable (si existe)
   - `ok`: `true` si existe, `false` si no

2. **`if val, ok := ...; ok`** (idiom común en Go):
   - Declara `val` y `ok` en el `if`
   - Si `ok` es `true` → la variable existe → devuelve `val`
   - Si `ok` es `false` → la variable NO existe → devuelve `def`

**¿Por qué no usar `os.Getenv()`?**
- `os.Getenv("KEY")` devuelve `""` si no existe
- No puedes distinguir entre "no existe" y "existe pero está vacía"
- `os.LookupEnv()` te dice explícitamente si existe

**Ejemplo:**
```go
// Si DATABASE_DSN no está definida
getEnv("DATABASE_DSN", "user:password@tcp(localhost:3306)/dbname")
// → Devuelve "user:password@tcp(localhost:3306)/dbname"

// Si DATABASE_DSN = "root:secret@tcp(db:3306)/prod"
getEnv("DATABASE_DSN", "user:password@tcp(localhost:3306)/dbname")
// → Devuelve "root:secret@tcp(db:3306)/prod"
```

---

## Flujo Completo de Ejecución

### Escenario: Iniciar la Aplicación

1. **Tu código llama:**
```go
cfg := config.MustLoad()
```

2. **`MustLoad()` llama a `Load()`**

3. **`Load()` ejecuta:**
   - Intenta cargar `.env` con `godotenv.Load()`
   - Si `.env` existe → carga sus variables
   - Si NO existe → continúa (sin error)

4. **`Load()` crea `Config` llamando a `getEnv()` 3 veces:**
   - `getEnv("APP_ENV", "dev")` → busca `APP_ENV` o usa `"dev"`
   - `getEnv("HTTP_PORT", "8080")` → busca `HTTP_PORT` o usa `"8080"`
   - `getEnv("DATABASE_DSN", "...")` → busca `DATABASE_DSN` o usa el default

5. **`MustLoad()` imprime la config:**
```
[config] env=dev port=8080 db=user:password@tcp(localhost:3306)/dbname
```

6. **`MustLoad()` devuelve el `*Config`**

7. **Tu código ahora puede usar:**
```go
server.Start(cfg.HttpPort)
db.Connect(cfg.DatabaseDSN)
```

---

## ¿Por Qué Este Diseño?

### ✅ Ventajas

1. **Seguridad**: No hardcodeas credenciales en el código
2. **Flexibilidad**: Cambias configuración sin recompilar
3. **Múltiples entornos**:
   - Desarrollo: `.env` local con datos de prueba
   - Producción: Variables de entorno del servidor
4. **Valores por defecto sensatos**: La app arranca aunque no configures nada
5. **Centralizado**: Toda la config en un solo lugar

### 🎯 Buenas Prácticas Implementadas

- ✅ Usa punteros (`*Config`) para evitar copias innecesarias
- ✅ Funciones simples y reutilizables (`getEnv`)
- ✅ No panicea si falta `.env` (flexible)
- ✅ Log de confirmación (debugging fácil)

---

## Ejemplo de Uso Completo

### Estructura de carpetas:
```
cellControl/
├── .env                    ← Archivo de configuración local
├── cmd/api/main.go
└── internal/config/config.go
```

### Archivo `.env`:
```env
APP_ENV=development
HTTP_PORT=3000
DATABASE_DSN=root:mypassword@tcp(localhost:3306)/cellcontrol
```

### En `main.go`:
```go
package main

import (
	"cellControl/internal/config"
	"cellControl/pkg/logger"
)

func main() {
	// Cargar configuración
	cfg := config.MustLoad()
	
	// Usar la configuración
	logger.Info("Iniciando aplicación en entorno: %s", cfg.Env)
	logger.Info("Servidor escuchando en puerto: %s", cfg.HttpPort)
	
	// Conectar a la base de datos
	// db.Connect(cfg.DatabaseDSN)
	
	// Iniciar servidor
	// server.Start(cfg.HttpPort)
}
```

### Salida:
```
[config] env=development port=3000 db=root:***@tcp(localhost:3306)/cellcontrol
INFO: 2025/11/23 14:30:45 main.go:10 Iniciando aplicación en entorno: development
INFO: 2025/11/23 14:30:45 main.go:11 Servidor escuchando en puerto: 3000
```

---

## Mejoras Posibles

Si quisieras extender este código:

1. **Validación de configuración**:
```go
func (c *Config) Validate() error {
	if c.DatabaseDSN == "" {
		return errors.New("DATABASE_DSN es requerido")
	}
	return nil
}
```

2. **Más campos de configuración**:
```go
type Config struct {
	Env          string
	HttpPort     string
	DatabaseDSN  string
	JWTSecret    string  // Para autenticación
	LogLevel     string  // debug, info, error
	MaxConns     int     // Conexiones máximas a la BD
}
```

3. **Soporte para múltiples archivos `.env`**:
```go
godotenv.Load(".env.local", ".env")  // Carga múltiples archivos
```

4. **Panic si faltan valores críticos**:
```go
func MustLoad() *Config {
	cfg := Load()
	if cfg.DatabaseDSN == "" {
		log.Fatal("DATABASE_DSN es requerido")
	}
	return cfg
}
```

---

## Resumen Rápido

| Función | Propósito | Cuándo Usar |
|---------|-----------|-------------|
| `Load()` | Carga config silenciosamente | Cuando no necesitas logs |
| `MustLoad()` | Carga config + muestra log | Al iniciar la aplicación |
| `getEnv()` | Lee variable o usa default | Uso interno (helper) |

**Patrón de diseño**: **Configuration Management** + **Environment Variables Pattern**

¿Necesitas que profundice en alguna parte o que implemente alguna mejora?

