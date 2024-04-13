package auth

import (
	"auth-service/pkg/auth/routes"
	"net/http"
)

func RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("/login", routes.LoginUser)
	r.HandleFunc("/register", routes.RegisterUser)
}