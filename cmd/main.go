package main

import (
	"log"
	"net/http"
	"time"

	"github.com/antgobar/famcal/config"
	"github.com/antgobar/famcal/internal/handlers"
	"github.com/antgobar/famcal/internal/middleware"
	"github.com/antgobar/famcal/internal/resources"
	"github.com/joho/godotenv"
)

func main() {
	cfg := config.MustLoadConfig()
	mux := http.NewServeMux()
	resources.Load(mux, cfg)
	handlers.Register(mux, cfg)
	stack := middleware.LoadMiddleware()

	server := http.Server{
		Addr:         "localhost:8090",
		Handler:      stack(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("Server starting on", server.Addr)
	server.ListenAndServe()
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Proceeding without .env file")
	}
}
