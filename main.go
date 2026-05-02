package main

import (
	"fmt"
	"net/http"
	"server/config"
	"server/routes"

	"github.com/go-chi/chi/v5"
)

func main() {
	config.InitDB()
	router := chi.NewRouter()
	routes.SetupRoutes(router)

	http.ListenAndServe(":8080", router)
	fmt.Println("Lancé")
}
