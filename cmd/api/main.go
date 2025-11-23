package api

import (
	"github.com/gonzalo-wi/cellcontrol/internal/config"
	"github.com/gonzalo-wi/cellcontrol/internal/db"
	httpServer "github.com/gonzalo-wi/cellcontrol/internal/http"
	"github.com/gonzalo-wi/cellcontrol/internal/http/handlers"
	"github.com/gonzalo-wi/cellcontrol/internal/repository"
	"github.com/gonzalo-wi/cellcontrol/internal/service"
	"github.com/gonzalo-wi/cellcontrol/pkg/logger"
)

func main() {
	logger.Init()

	cfg := config.MustLoad()
	dbConn := db.NewDatabase(cfg)

	userRepo := repository.NewUserRepository(dbConn)
	userSvc := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)

	server := httpServer.NewServer(cfg, userHandler)

	if err := server.Run(); err != nil {
		logger.Error("error al iniciar servidor: %v", err)
	}
}
