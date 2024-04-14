package routes

import (
	"auth-service/pkg/auth/handlers"
	"net/http"
)

func RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("/login", handlers.LoginUser)
	r.HandleFunc("/register", handlers.RegisterUser)
}