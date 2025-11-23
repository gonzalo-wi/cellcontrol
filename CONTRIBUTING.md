# Guía de Contribución

¡Gracias por tu interés en contribuir a CellControl! 🎉

## 📋 Tabla de Contenidos

- [Código de Conducta](#código-de-conducta)
- [¿Cómo Contribuir?](#cómo-contribuir)
- [Configuración del Entorno](#configuración-del-entorno)
- [Proceso de Pull Request](#proceso-de-pull-request)
- [Estándares de Código](#estándares-de-código)
- [Commits](#commits)
- [Reportar Bugs](#reportar-bugs)
- [Sugerir Mejoras](#sugerir-mejoras)

## 🤝 Código de Conducta

Este proyecto se adhiere a un código de conducta. Al participar, se espera que mantengas un ambiente respetuoso y constructivo.

## 🚀 ¿Cómo Contribuir?

1. **Fork** el repositorio
2. **Crea** una rama desde `main`
3. **Implementa** tus cambios
4. **Escribe** tests si es necesario
5. **Asegúrate** de que todos los tests pasen
6. **Commit** tus cambios siguiendo las convenciones
7. **Push** a tu fork
8. **Abre** un Pull Request

## ⚙️ Configuración del Entorno

### Prerrequisitos

- Go 1.24+
- MySQL 8.0+
- Git

### Setup Local

```bash
# Clonar tu fork
git clone https://github.com/TU_USUARIO/cellcontrol.git
cd cellcontrol

# Añadir el repositorio original como remote
git remote add upstream https://github.com/gonzalo-wi/cellcontrol.git

# Instalar dependencias
go mod download

# Copiar configuración de ejemplo
cp .env.example .env

# Editar .env con tus credenciales locales
# Crear la base de datos
mysql -u root -p -e "CREATE DATABASE cellcontrol"

# Ejecutar la aplicación
go run cmd/api/main.go
```

## 🔄 Proceso de Pull Request

1. **Actualiza tu fork** con los últimos cambios:
   ```bash
   git fetch upstream
   git checkout main
   git merge upstream/main
   ```

2. **Crea una rama** con un nombre descriptivo:
   ```bash
   git checkout -b feature/descripcion-corta
   ```

3. **Haz tus cambios** y commitea:
   ```bash
   git add .
   git commit -m "feat: descripción del cambio"
   ```

4. **Push** a tu fork:
   ```bash
   git push origin feature/descripcion-corta
   ```

5. **Abre un Pull Request** en GitHub con:
   - Título descriptivo
   - Descripción detallada de los cambios
   - Referencias a issues relacionados
   - Capturas de pantalla (si aplica)

### Checklist del Pull Request

- [ ] El código compila sin errores
- [ ] Todos los tests pasan (`go test ./...`)
- [ ] He añadido tests para mis cambios
- [ ] He actualizado la documentación si es necesario
- [ ] Mi código sigue los estándares del proyecto
- [ ] Los commits siguen las convenciones

## 📝 Estándares de Código

### Formato

- Usa `gofmt` o `goimports` para formatear el código:
  ```bash
  go fmt ./...
  ```

### Naming Conventions

- **Paquetes**: minúsculas, sin guiones bajos
  ```go
  package userservice
  ```

- **Exportado (público)**: PascalCase
  ```go
  type UserService interface {}
  func NewUserService() {}
  ```

- **No exportado (privado)**: camelCase
  ```go
  type userService struct {}
  func getUserByID() {}
  ```

### Estructura de Archivos

- Un tipo principal por archivo
- Nombrar archivos con snake_case: `user_service.go`
- Agrupar código relacionado

### Comentarios

- Comentar funciones/tipos exportados:
  ```go
  // NewUserService creates a new user service instance
  func NewUserService(repo UserRepository) UserService {
  ```

- Explicar lógica compleja
- Evitar comentarios obvios

### Error Handling

```go
// ✅ Bien
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// ❌ Mal
if err != nil {
    return err  // Sin contexto
}
```

## 💬 Commits

Seguimos [Conventional Commits](https://www.conventionalcommits.org/):

### Formato

```
<tipo>(<scope>): <descripción>

[cuerpo opcional]

[footer opcional]
```

### Tipos

- `feat`: Nueva funcionalidad
- `fix`: Corrección de bug
- `docs`: Cambios en documentación
- `style`: Formato, sin cambios en lógica
- `refactor`: Refactorización de código
- `test`: Añadir o modificar tests
- `chore`: Tareas de mantenimiento

### Ejemplos

```bash
feat(user): agregar endpoint de actualización de usuario
fix(auth): corregir validación de email
docs(readme): actualizar instrucciones de instalación
refactor(service): mejorar estructura del user service
test(handler): añadir tests para CreateUser
```

## 🐛 Reportar Bugs

Usa los [GitHub Issues](https://github.com/gonzalo-wi/cellcontrol/issues) para reportar bugs.

### Template de Bug Report

```markdown
**Descripción del Bug**
Descripción clara y concisa del bug.

**Pasos para Reproducir**
1. Ir a '...'
2. Hacer click en '...'
3. Ver error

**Comportamiento Esperado**
Lo que debería suceder.

**Comportamiento Actual**
Lo que actualmente sucede.

**Screenshots**
Si aplica, añade capturas de pantalla.

**Entorno**
- OS: [e.g. macOS 13.0]
- Go version: [e.g. 1.24.5]
- MySQL version: [e.g. 8.0.33]

**Información Adicional**
Cualquier otro contexto relevante.
```

## 💡 Sugerir Mejoras

Usa los [GitHub Issues](https://github.com/gonzalo-wi/cellcontrol/issues) con la etiqueta `enhancement`.

### Template de Feature Request

```markdown
**¿Es tu feature request relacionado a un problema?**
Descripción clara del problema.

**Describe la solución que te gustaría**
Descripción clara y concisa de lo que quieres que suceda.

**Describe alternativas que has considerado**
Otras soluciones o features que has considerado.

**Contexto adicional**
Cualquier otro contexto, screenshots, etc.
```

## 📚 Recursos Útiles

- [Documentación del Proyecto](./docs/README.md)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Gin Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)

## ❓ Preguntas

Si tienes preguntas, puedes:
- Abrir un [GitHub Issue](https://github.com/gonzalo-wi/cellcontrol/issues)
- Contactar al mantenedor

---

¡Gracias por contribuir! 🙌

