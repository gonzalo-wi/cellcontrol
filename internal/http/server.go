package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gonzalo-wi/cellcontrol/internal/config"
	"github.com/gonzalo-wi/cellcontrol/internal/http/handlers"
	"github.com/gonzalo-wi/cellcontrol/pkg/logger"
)

type Server struct {
	engine *gin.Engine
	cfg    *config.Config
}

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

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%s", s.cfg.HttpPort)
	logger.Info("escuchando en %s", addr)
	return s.engine.Run(addr)
}
