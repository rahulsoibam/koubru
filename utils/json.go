package utils

import (
	"encoding/json"
	"net/http"
)

// RespondWithError function to respond with error in JSON format
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON function to return response in JSON format
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	// Set headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithMessage(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"message": message})
}
