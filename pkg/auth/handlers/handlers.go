package handlers

import (
	"auth-service/pkg/database"
	"auth-service/pkg/jwt"
	"auth-service/pkg/auth/models"
	"auth-service/pkg/utils"
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func LoginUser(rw http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(rw, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user := models.User{}
	err := database.DB.QueryRow(
		"SELECT * FROM users WHERE email = $1", 
		loginRequest.Email,
	).Scan(
		&user.ID, 
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.UserType,
	)

	if err != nil {
		http.Error(rw, "Cannot fetch user from database!", http.StatusInternalServerError)
		return
	}

	if err := user.CheckPassword(loginRequest.Password); err != nil {
		http.Error(rw, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := jwt.GenerateJWT(user.Email)
	if err != nil {
		http.Error(rw, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{Token: token}
	utils.RespondJSON(rw, http.StatusOK, response)
}

func RegisterUser(rw http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(rw, "Invalid request payload", http.StatusBadRequest)
		return
	}
	
	if err := user.HashPassword(user.Password); err != nil {
		http.Error(rw, "Error hashing password", http.StatusInternalServerError)
		return
	}

	err := database.DB.QueryRow(
		"INSERT INTO users (first_name, last_name, email, password, user_type) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.UserType,
	).Scan(&user.ID)

	if err != nil {
		http.Error(rw, "Error inserting user into database", http.StatusInternalServerError)
		return
	}

	token, err := jwt.GenerateJWT(user.Email)
	if err != nil {
		http.Error(rw, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{Token: token}
	utils.RespondJSON(rw, http.StatusCreated, response)
}
