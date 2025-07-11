package helpers

import (
	"encoding/json"
	"net/http"
)

// RespondWithJSON writes the given payload as a JSON response with the specified HTTP status code.
//
// It sets the "Content-Type" header to "application/json" and writes the status code before encoding.
//
// Returns any error encountered during JSON encoding.
//
// Example:
//
//	err := RespondWithJSON(w, http.StatusOK, map[string]string{"message": "success"})
//	if err != nil {
//	    // handle error
//	}
func RespondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

// RespondWithError writes a JSON response containing an error message with the specified HTTP status code.
//
// It uses RespondWithJSON internally to send a standardized error response of the form:
// {"error": "message"}
//
// Returns any error encountered during JSON encoding.
//
// Example:
//
//	err := RespondWithError(w, http.StatusBadRequest, "invalid request payload")
//	if err != nil {
//	    // handle error
//	}
func RespondWithError(w http.ResponseWriter, status int, message string) error {
	return RespondWithJSON(w, status, map[string]string{"error": message})
}

// NoContent sends an HTTP 204 No Content response without a body.
//
// Example:
//
//	NoContent(w)  // Sends HTTP status 204 with empty body
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
