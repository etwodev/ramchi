package helpers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// URLParam returns the url parameter from a http.Request object.
func URLParam(r *http.Request, key string) string {
	if value := chi.URLParam(r, key); value != "" {
		return value
	}
	return ""
}
