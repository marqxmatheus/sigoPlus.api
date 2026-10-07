package main

import (
	"log"

	"sigo/internal/config"
	"sigo/internal/router"
)

func main() {
	cfg := config.Load()

	r := router.New(cfg)

	log.Printf("SIGO+ API rodando na porta %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
