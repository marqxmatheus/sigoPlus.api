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

	servicoRepo := repository.NewServicoRepository(pool)
	servicoService := service.NewServicoService(servicoRepo)
	servicoHandler := handler.NewServicoHandler(servicoService)

	pecaRepo := repository.NewPecaRepository(pool)
	pecaService := service.NewPecaService(pecaRepo)
	pecaHandler := handler.NewPecaHandler(pecaService)

	ordemRepo := repository.NewOrdemServicoRepository(pool)
	ordemService := service.NewOrdemServicoService(ordemRepo, clienteRepo, veiculoRepo, servicoRepo, pecaRepo)
	ordemHandler := handler.NewOrdemServicoHandler(ordemService)

	r := router.New(cfg, router.Handlers{
		Cliente:      clienteHandler,
		Veiculo:      veiculoHandler,
		Servico:      servicoHandler,
		Peca:         pecaHandler,
		OrdemServico: ordemHandler,
	})

	log.Printf("SIGO+ API rodando na porta %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
