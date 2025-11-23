# Explicación del Código del Logger
¿Necesitas que profundice en alguna parte específica o que implemente alguna de estas mejoras?

---

5. **Colores**: Para distinguir visualmente en consola
4. **Configuración**: Activar/desactivar niveles según entorno
3. **Logs estructurados**: Formato JSON para parseo automático
2. **Logs a archivo**: Además de consola
1. **Niveles adicionales**: `Debug()`, `Warning()`, `Fatal()`

Si quisieras extender este logger, podrías añadir:

## Posibles Mejoras Futuras

---

```
[ERROR] 2025/11/23 14:30:48 user.go:67 Error al guardar usuario: connection timeout
[ERROR] 2025/11/23 14:30:47 db.go:34 Conexión perdida
INFO: 2025/11/23 14:30:46 auth.go:23 Usuario admin inició sesión desde IP 192.168.1.1
INFO: 2025/11/23 14:30:45 main.go:15 Servidor iniciado en puerto 8080
```
**Salida esperada:**

```
logger.Error("Error al guardar usuario: %v", err)
// Error con formato

logger.Error("Conexión perdida")
// Error simple

logger.Info("Usuario %s inició sesión desde IP %s", username, ip)
// Mensaje con formato

logger.Info("Servidor iniciado en puerto 8080")
// Mensaje simple

import "cellControl/pkg/logger"
```go

## Ejemplos de Uso

---

   - Útil para debugging
   - `log.Lshortfile` muestra dónde se generó el log
5. **Rastreo de código**:

   - `import "tu-proyecto/pkg/logger"`
   - No necesitas configurar nada, solo importar
4. **Inicialización automática**:

   - Fácil de usar: `logger.Info("mensaje")` o `logger.Error("error: %v", err)`
   - Soporta mensajes simples y con formato
3. **Flexibilidad**:

   - Incluyen fecha, hora y ubicación en el código
   - Todos los logs tienen el mismo formato
2. **Consistencia**:

   - Facilita filtrar errores en producción
   - Logs INFO y ERROR van a diferentes streams
1. **Separación de concerns**:

### ✅ Ventajas

## ¿Por Qué Este Diseño?

---

- Tiene el prefijo `[ERROR]` en lugar de `INFO:`
- Sale por `stderr` en lugar de `stdout`
- Usa `errorLogger` en lugar de `infoLogger`
**Funciona igual que `Info()`**, pero:

**Propósito**: Registrar mensajes de error

```
}
	}
		errorLogger.Println(msg)
	} else {
		errorLogger.Printf(msg, args...)
	if len(args) > 0 {
func Error(msg string, args ...any) {
```go
### Función Error()

---

- `args...` al llamar → expande el slice como argumentos individuales
- `args ...any` en la firma → acepta múltiples argumentos
**El operador `...`** (spread):

  - Ejemplo: `Info("Servidor iniciado")`
  - Usa `Println` para simplemente **imprimir** el mensaje
- Si NO hay argumentos:
  
  - Ejemplo: `Info("Usuario %s conectado", username)`
  - Usa `Printf` para **formatear** el mensaje con los argumentos
- Si hay argumentos (`len(args) > 0`):
**¿Por qué la condición `if`?**

- `args ...any`: Argumentos variádicos (0 o más argumentos de cualquier tipo)
- `msg string`: El mensaje a registrar
**Parámetros**:

**Propósito**: Registrar mensajes informativos

```
}
	}
		infoLogger.Println(msg)
	} else {
		infoLogger.Printf(msg, args...)
	if len(args) > 0 {
func Info(msg string, args ...any) {
```go
### Función Info()

---

```
[ERROR] 2025/11/23 14:30:47 db.go:45 Error conectando a la base de datos
INFO: 2025/11/23 14:30:45 main.go:12 Aplicación iniciada
```
**Ejemplo de salida:**

   - El operador `|` combina múltiples flags
   - `log.LstdFlags` → equivale a `log.Ldate|log.Ltime` (combinación estándar)
   - `log.Lshortfile` → añade archivo:línea (logger.go:23)
   - `log.Ltime` → añade la hora (01:23:23)
   - `log.Ldate` → añade la fecha (2009/01/23)
3. **Flags de formato**:

   - `"[ERROR] "` → cada mensaje ERROR comenzará con este texto
   - `"INFO: "` → cada mensaje INFO comenzará con este texto
2. **Prefijo**:

   - `os.Stderr` para ERROR → sale por consola de errores (error estándar)
   - `os.Stdout` para INFO → sale por consola normal (salida estándar)
1. **Destino de salida**:

**`log.New()` recibe 3 parámetros:**

- Perfecta para **inicializar** los loggers
- Se ejecuta **antes** que cualquier otra función
- Se ejecuta **automáticamente** cuando se importa el paquete
- La función `init()` es especial en Go
**¿Por qué `init()`?**

```
}
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
func init() {
```go
### Función init()

---

- Se declaran a nivel de paquete para ser accesibles por todas las funciones
- Son **punteros** a `log.Logger` (el tipo `*` indica puntero)
- **`errorLogger`**: Logger dedicado para mensajes de error
- **`infoLogger`**: Logger dedicado para mensajes informativos
```
)
	errorLogger *log.Logger
	infoLogger  *log.Logger
var (
```go
### Variables Globales

---

- **`os`**: Paquete para interactuar con el sistema operativo (se usa para `os.Stdout` y `os.Stderr`)
- **`log`**: Paquete estándar de Go para logging
```
)
	"os"
	"log"
import (
```go

- Define el paquete `logger` que puede ser importado por otras partes de tu aplicación
```
package logger
```go
### Package y Imports

## Desglose Línea por Línea

---

Este es un **paquete de logging personalizado** para tu aplicación Go. Proporciona funciones simples para registrar mensajes informativos y de error.
## Visión General


