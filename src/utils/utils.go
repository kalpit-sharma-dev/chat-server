package utils

import (
	"encoding/json"
	"net/http"
	"regexp"
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

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RemoveAllButNumbersAndPlus(input string) string {
	// Compile the regular expression to match all characters except digits and '+'
	re := regexp.MustCompile(`[^0-9+]`)
	// Replace all matches with an empty string
	return re.ReplaceAllString(input, "")
}
