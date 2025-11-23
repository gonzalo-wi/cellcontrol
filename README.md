# 📱 CellControl

Sistema de gestión y control de celulares corporativos desarrollado en Go con arquitectura limpia (Clean Architecture).

## 🚀 Características

- ✅ API RESTful con Gin Framework
- ✅ Arquitectura en capas (Handler → Service → Repository)
- ✅ ORM con GORM para MySQL
- ✅ Configuración mediante variables de entorno
- ✅ Sistema de logging personalizado
- ✅ Validación automática de datos
- ✅ Separación de responsabilidades (Clean Architecture)

## 🏗️ Arquitectura

```
┌─────────────────────────────────────┐
│      HTTP REQUEST (Cliente)         │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│     HTTP SERVER (Gin)               │
│  - Routing                          │
│  - Middlewares                      │
│  - Health Check                     │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│        HANDLERS                     │
│  - Parseo de JSON                   │
│  - Validación de entrada            │
│  - Respuestas HTTP                  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│        SERVICES                     │
│  - Lógica de negocio                │
│  - Transformaciones                 │
│  - Validaciones complejas           │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│      REPOSITORIES                   │
│  - CRUD operations                  │
│  - Queries con GORM                 │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│      BASE DE DATOS (MySQL)          │
└─────────────────────────────────────┘
```

## 📁 Estructura del Proyecto

```
cellControl/
├── cmd/
│   └── api/
│       └── main.go              # Punto de entrada de la aplicación
├── internal/
│   ├── config/
│   │   └── config.go            # Gestión de configuración
│   ├── db/
│   │   └── db.go                # Conexión a base de datos
│   ├── domain/
│   │   ├── cellphone.go         # Modelo de celular
│   │   ├── history.go           # Modelo de historial
│   │   └── user.go              # Modelo de usuario
│   ├── http/
│   │   ├── server.go            # Configuración del servidor HTTP
│   │   └── handlers/
│   │       ├── cellphone_handler.go
│   │       └── user_handler.go
│   ├── repository/
│   │   ├── cellphone_repository.go
│   │   └── user_repository.go
│   └── service/
│       ├── cellphone_service.go
│       └── user_service.go
├── pkg/
│   └── logger/
│       └── logger.go            # Sistema de logging
├── docs/                        # 📚 Documentación detallada
│   ├── README.md
│   ├── 01_logger.md
│   ├── 02_config.md
│   ├── 03_user_repository.md
│   ├── 04_user_service.md
│   ├── 05_server.md
│   └── 06_user_handler.md
├── .env                         # Variables de entorno (no subir a git)
├── .gitignore
├── go.mod
├── go.sum
└── README.md                    # Este archivo
```

## 🛠️ Tecnologías

- **[Go 1.24+](https://golang.org/)** - Lenguaje de programación
- **[Gin](https://gin-gonic.com/)** - Framework web HTTP
- **[GORM](https://gorm.io/)** - ORM para Go
- **[MySQL](https://www.mysql.com/)** - Base de datos relacional
- **[godotenv](https://github.com/joho/godotenv)** - Gestión de variables de entorno

## ⚙️ Requisitos Previos

- Go 1.24 o superior
- MySQL 8.0 o superior
- Git

## 🚀 Instalación

### 1. Clonar el repositorio

```bash
git clone https://github.com/gonzalo-wi/cellcontrol.git
cd cellcontrol
```

### 2. Instalar dependencias

```bash
go mod download
# o
go mod tidy
```

### 3. Configurar variables de entorno

Crea un archivo `.env` en la raíz del proyecto:

```env
# Entorno
APP_ENV=development

# Servidor HTTP
HTTP_PORT=8080

# Base de Datos
DATABASE_DSN=root:password@tcp(localhost:3306)/cellcontrol?charset=utf8mb4&parseTime=True&loc=Local
```

> **Nota:** Ajusta los valores según tu configuración local de MySQL.

### 4. Crear la base de datos

```sql
CREATE DATABASE cellcontrol CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. Ejecutar migraciones (si las hay)

```bash
# Las migraciones se ejecutarán automáticamente al iniciar la app
# o ejecuta manualmente las migraciones si tienes un script
```

### 6. Ejecutar la aplicación

```bash
go run cmd/api/main.go
```

La aplicación estará disponible en `http://localhost:8080`

## 📡 API Endpoints

### Health Check

```http
GET /health
```

**Respuesta:**
```json
{
  "status": "ok"
}
```

### Usuarios

#### Crear Usuario

```http
POST /api/v1/usuarios
Content-Type: application/json

{
  "nombre": "Juan",
  "apellido": "Pérez",
  "email": "juan@example.com",
  "reparto": "IT"
}
```

**Respuesta (201 Created):**
```json
{
  "message": "usuario creado exitosamente"
}
```

#### Listar Usuarios

```http
GET /api/v1/usuarios
```

**Respuesta (200 OK):**
```json
[
  {
    "id": 1,
    "nombre": "Juan",
    "apellido": "Pérez",
    "email": "juan@example.com",
    "reparto": "IT",
    "created_at": "2025-11-23T14:30:45Z",
    "updated_at": "2025-11-23T14:30:45Z"
  }
]
```

### Celulares

#### Crear Celular

```http
POST /api/v1/celulares
Content-Type: application/json

{
  "marca": "Samsung",
  "modelo": "Galaxy S21",
  "imei": "123456789012345",
  "numero": "+54 9 11 1234-5678"
}
```

#### Listar Celulares

```http
GET /api/v1/celulares
```

## 🧪 Testing

```bash
# Ejecutar todos los tests
go test ./...

# Ejecutar tests con cobertura
go test -cover ./...

# Ejecutar tests de un paquete específico
go test ./internal/service
```

## 📚 Documentación

La documentación detallada de cada componente está disponible en la carpeta [`docs/`](./docs/):

- **[01_logger.md](./docs/01_logger.md)** - Sistema de logging explicado
- **[02_config.md](./docs/02_config.md)** - Configuración y variables de entorno
- **[03_user_repository.md](./docs/03_user_repository.md)** - Patrón Repository
- **[04_user_service.md](./docs/04_user_service.md)** - Lógica de negocio
- **[05_server.md](./docs/05_server.md)** - Servidor HTTP con Gin
- **[06_user_handler.md](./docs/06_user_handler.md)** - Controladores HTTP

Cada documento incluye:
- Explicación conceptual
- Desglose línea por línea del código
- Ejemplos prácticos
- Diagramas de flujo
- Mejoras sugeridas

## 🔧 Desarrollo

### Estructura de Commits

Seguimos [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: agregar endpoint de actualización de usuarios
fix: corregir validación de email
docs: actualizar README con nuevos endpoints
refactor: mejorar estructura del user service
```

### Ejecutar en modo desarrollo

```bash
# Con hot reload (usando air)
air

# O manualmente
go run cmd/api/main.go
```

### Variables de Entorno

| Variable | Descripción | Valor por Defecto |
|----------|-------------|-------------------|
| `APP_ENV` | Entorno de ejecución | `dev` |
| `HTTP_PORT` | Puerto del servidor HTTP | `8080` |
| `DATABASE_DSN` | Cadena de conexión a MySQL | `user:password@tcp(localhost:3306)/dbname` |

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/amazing-feature`)
3. Commit tus cambios (`git commit -m 'feat: add amazing feature'`)
4. Push a la rama (`git push origin feature/amazing-feature`)
5. Abre un Pull Request

## 📝 Buenas Prácticas Implementadas

- ✅ **Clean Architecture** - Separación clara de responsabilidades
- ✅ **Dependency Injection** - Facilita testing y mantenibilidad
- ✅ **Interface Segregation** - Interfaces específicas para cada capa
- ✅ **Error Handling** - Manejo consistente de errores
- ✅ **Validation** - Validación automática con binding tags
- ✅ **Logging** - Sistema de logs estructurado
- ✅ **Configuration** - Gestión de config mediante env vars

## 🐛 Problemas Conocidos

Ninguno por el momento. Si encuentras algún problema, por favor [abre un issue](https://github.com/gonzalo-wi/cellcontrol/issues).

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

## 👤 Autor

**Gonzalo Wiñazki**

- GitHub: [@gonzalo-wi](https://github.com/gonzalo-wi)
- Email: tu-email@example.com

## 🙏 Agradecimientos

- [Gin Framework](https://gin-gonic.com/) - Por el excelente framework web
- [GORM](https://gorm.io/) - Por el poderoso ORM
- Comunidad de Go - Por los recursos y apoyo

---

⭐️ Si este proyecto te fue útil, considera darle una estrella en GitHub!

