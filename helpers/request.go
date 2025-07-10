package helpers

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// GetHeader safely retrieves a header key
func GetHeader(r *http.Request, key string) string {
	return strings.TrimSpace(r.Header.Get(key))
}

// GetBearerToken extracts token from Authorization header
func GetBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

// IsJSONRequest checks Content-Type
func IsJSONRequest(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")
}

// GetIP retrieves the real IP address from the request, accounting for proxies.
func GetIP(r *http.Request) string {
	// Try X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	// Try X-Real-IP header
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}

	// Fallback to remote address
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// BindJSON decodes JSON body into the provided struct.
func BindJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

// RouteContext returns the chi.RouteContext from the request context.
func RouteContext(r *http.Request) *chi.Context {
	return chi.RouteContext(r.Context())
}

// URLParam returns a URL parameter from a http.Request object.
func URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
