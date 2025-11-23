package db

import (
	"fmt"
	"github.com/gonzalo-wi/cellcontrol/internal/config"
	"github.com/gonzalo-wi/cellcontrol/internal/domain"
	"github.com/gonzalo-wi/cellcontrol/pkg/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		logger.Error("error al conectar la Base de Datos: %v", err)
		panic(fmt.Errorf("no se pudo conectar a la Base de Datos: %w", err))
	}
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		logger.Error("error al migrar la Base de Datos: %v", err)
		panic(fmt.Errorf("no se pudo migrar la Base de Datos: %w", err))
	}
	logger.Info("Base de datos inicializada correctamente")
	return db
}
