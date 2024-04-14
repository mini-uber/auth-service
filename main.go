package main

import (
	"auth-service/pkg/auth/routes"
	"auth-service/pkg/database"
	"auth-service/pkg/config"
	"fmt"
	"log"
	"net/http"

)

func main() {
	config := config.LoadConfig()
	database.Connect(config.DBUrl)

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux)
	addr := fmt.Sprintf(":%s", config.Port)
	log.Println("Starting server on ", addr)
	err := http.ListenAndServe(addr, mux)
	log.Fatal(err)
}