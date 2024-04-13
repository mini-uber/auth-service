package routes

import (
	"auth-service/pkg/auth/database"
	"auth-service/pkg/auth/jwt"
	"auth-service/pkg/auth/models"
	"auth-service/pkg/auth/utils"
	"encoding/json"
	"fmt"
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
		utils.HandleError(rw, http.StatusBadRequest, "Invalid request payload")
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
		utils.HandleError(rw, http.StatusInternalServerError, "Error querying database")
		return
	}

	if err := user.CheckPassword(loginRequest.Password); err != nil {
		utils.HandleError(rw, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := jwt.GenerateJWT(user.Email)
	if err != nil {
		utils.HandleError(rw, http.StatusInternalServerError, "Error generating token")
		return
	}

	response := AuthResponse{Token: token}
	utils.RespondJSON(rw, http.StatusOK, response)
}

func RegisterUser(rw http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.HandleError(rw, http.StatusBadRequest, "Invalid request payload")
		return
	}
	
	if err := user.HashPassword(user.Password); err != nil {
		utils.HandleError(rw, http.StatusInternalServerError, "Error hashing password")
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
		fmt.Println(err.Error())
		utils.HandleError(rw, http.StatusInternalServerError, "Error creating user")
		return
	}

	token, err := jwt.GenerateJWT(user.Email)
	if err != nil {
		utils.HandleError(rw, http.StatusInternalServerError, "Error generating token")
		return
	}

	response := AuthResponse{Token: token}
	utils.RespondJSON(rw, http.StatusCreated, response)
}
