package utils

import (
	"encoding/json"
	"net/http"
)

func HandleError(rw http.ResponseWriter, statusCode int, message string) {
	rw.WriteHeader(statusCode)
	rw.Header().Set("Content-Type", "application/json")
	response := map[string]string{"error": message}
	json.NewEncoder(rw).Encode(response)
}

func RespondJSON(rw http.ResponseWriter, statusCode int, data interface{}) {
	rw.WriteHeader(statusCode)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}