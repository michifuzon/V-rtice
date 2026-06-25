package config

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func Cargar() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró .env, usando variables de entorno del sistema")
	}

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		ServerPort: os.Getenv("SERVER_PORT"),
	}

	// Si no hay variables individuales, parsear DATABASE_URL (Render)
	if cfg.DBHost == "" {
		if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
			parseDatabaseURL(dbURL, cfg)
		}
	}

	// Render inyecta PORT; SERVER_PORT es para desarrollo local
	if cfg.ServerPort == "" {
		cfg.ServerPort = os.Getenv("PORT")
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	return cfg
}

func parseDatabaseURL(rawURL string, cfg *Config) {
	if !strings.Contains(rawURL, "://") {
		rawURL = "mysql://" + rawURL
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("Error parseando DATABASE_URL: %v", err)
		return
	}
	cfg.DBUser = u.User.Username()
	cfg.DBPassword, _ = u.User.Password()
	cfg.DBHost = u.Hostname()
	cfg.DBPort = u.Port()
	if cfg.DBPort == "" {
		cfg.DBPort = "3306"
	}
	cfg.DBName = strings.TrimPrefix(u.Path, "/")
}
