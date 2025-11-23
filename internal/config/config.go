package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	HttpPort    string
	DatabaseDSN string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		Env:         getEnv("APP_ENV", "dev"),
		HttpPort:    getEnv("HTTP_PORT", "8080"),
		DatabaseDSN: getEnv("DATABASE_DSN", "user:password@tcp(localhost:3306)/dbname"),
	}
}

func MustLoad() *Config {
	cfg := Load()
	log.Printf("[config] env=%s port=%s db=%s\n", cfg.Env, cfg.HttpPort, cfg.DatabaseDSN)
	return cfg
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}
