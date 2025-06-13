package main

import (
	"fmt"
	"log"

	"lockify-back/internal/config"
	"lockify-back/internal/database"
	"lockify-back/internal/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Erro ao inicializar banco de dados: %v", err)
	}

	router := routes.SetupRoutes(db, cfg)

	serverAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Servidor iniciado em http://%s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
