package main

import (
	"log"
	"net/http"

	"github.com/antgobar/famcal/handlers"
	"github.com/antgobar/famcal/middleware"
	"github.com/antgobar/famcal/resources"
	"github.com/antgobar/famcal/store"
)

func init() {
	store.MembersStore.PopulateDummy()
}

func main() {

	router := http.NewServeMux()
	resources.Load(router)
	handlers.Register(router)
	stack := middleware.LoadMiddleware()

	const addr = "localhost:8080"

	server := http.Server{
		Addr:    addr,
		Handler: stack(router),
	}

	log.Println("Server starting on", addr)
	server.ListenAndServe()
}
