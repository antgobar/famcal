package main

import (
	"log"
	"net/http"
	"time"

	"github.com/antgobar/famcal/config"
	"github.com/antgobar/famcal/handlers"
	"github.com/antgobar/famcal/middleware"
	"github.com/antgobar/famcal/resources"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	cfg := mustLoadConfig()

	mux := http.NewServeMux()
	resources.Load(mux, cfg)
	handlers.Register(mux, cfg)
	stack := middleware.LoadMiddleware()

	server := http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      stack(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("Server starting on", cfg.ServerAddress)
	server.ListenAndServe()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func mustLoadConfig() *config.Config {
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}
	return config

}
