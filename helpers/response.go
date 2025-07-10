package helpers

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON writes a JSON response with status code
func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

// RespondWithError writes a standardized error message
func RespondWithError(w http.ResponseWriter, status int, message string) error {
	return RespondWithJSON(w, status, map[string]string{"error": message})
}

// NoContent sends a 204 No Content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
