package main

import (
	"log"

	"sigo/internal/config"
	"sigo/internal/database"
	"sigo/internal/router"
)

func main() {
	cfg := config.Load()

	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("erro no banco de dados: %v", err)
	}
	defer pool.Close()
	log.Println("conectado ao PostgreSQL")

	r := router.New(cfg)

	log.Printf("SIGO+ API rodando na porta %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
