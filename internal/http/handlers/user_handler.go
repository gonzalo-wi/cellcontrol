package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gonzalo-wi/cellcontrol/internal/service"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type createUserRequest struct {
	Nombre   string `json:"nombre" binding:"required"`
	Apellido string `json:"apellido" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Reparto  string `json:"reparto" binding:"required"`
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/usuarios", h.CreateUser)
	r.GET("/usuarios", h.ListUsers)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.CreateUser(req.Nombre, req.Apellido, req.Email, req.Reparto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear el usuario"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "usuario creado exitosamente"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.svc.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo obtener usuarios"})
		return
	}
	c.JSON(http.StatusOK, users)
}
