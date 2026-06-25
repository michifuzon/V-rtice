package main

import (
	"fmt"
	"log"

	"github.com/isaacunaa/ticketek-ds2026/backend/internal/config"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/domain"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/router"
	"github.com/isaacunaa/ticketek-ds2026/backend/internal/seed"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Cargar()

	dsn := cfg.DatabaseURL
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error conectando a la base de datos: ", err)
	}
	log.Println("Conexión exitosa a MySQL")

	err = db.AutoMigrate(
		&domain.Usuario{},
		&domain.Evento{},
		&domain.Entrada{},
		&domain.Favorito{},
	)
	if err != nil {
		log.Fatal("Error en AutoMigrate: ", err)
	}
	log.Println("Migraciones aplicadas correctamente")

	seed.Ejecutar(db)

	r := router.Configurar(db)

	log.Printf("Servidor escuchando en puerto %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Error iniciando servidor: ", err)
	}
}
