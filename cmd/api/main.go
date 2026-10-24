package main

import (
	"log"

	"sigo/internal/config"
	"sigo/internal/database"
	"sigo/internal/handler"
	"sigo/internal/repository"
	"sigo/internal/router"
	"sigo/internal/service"
)

func main() {
	cfg := config.Load()

	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("erro no banco de dados: %v", err)
	}
	defer pool.Close()
	log.Println("conectado ao PostgreSQL")

	clienteRepo := repository.NewClienteRepository(pool)
	clienteService := service.NewClienteService(clienteRepo)
	clienteHandler := handler.NewClienteHandler(clienteService)

	veiculoRepo := repository.NewVeiculoRepository(pool)
	veiculoService := service.NewVeiculoService(veiculoRepo, clienteRepo)
	veiculoHandler := handler.NewVeiculoHandler(veiculoService)

	r := router.New(cfg, router.Handlers{
		Cliente: clienteHandler,
		Veiculo: veiculoHandler,
	})

	log.Printf("SIGO+ API rodando na porta %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
