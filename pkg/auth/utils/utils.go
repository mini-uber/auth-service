package utils

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(rw http.ResponseWriter, statusCode int, data interface{}) {
	rw.WriteHeader(statusCode)
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}