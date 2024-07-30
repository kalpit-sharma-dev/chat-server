package utils

import (
	"encoding/json"
	"net/http"
)

// RespondJSON sends a JSON response with the provided status code and data
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// RespondError sends a JSON response with an error message
func RespondError(w http.ResponseWriter, statusCode int, message string) {
	RespondJSON(w, statusCode, map[string]string{"error": message})
}
