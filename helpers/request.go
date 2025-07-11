package helpers

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// GetHeader safely retrieves the value of the specified HTTP header key from the request,
// trimming any leading or trailing whitespace.
//
// Example:
//
//	val := GetHeader(r, "X-Custom-Header")
//	fmt.Println(val) // Output: "header-value"
func GetHeader(r *http.Request, key string) string {
	return strings.TrimSpace(r.Header.Get(key))
}

// GetBearerToken extracts the Bearer token from the Authorization header of the request.
//
// If the Authorization header does not start with "Bearer ", an empty string is returned.
//
// Example:
//
//	token := GetBearerToken(r)
//	fmt.Println(token) // Output: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
func GetBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return ""
}

// IsJSONRequest checks if the Content-Type header of the request indicates JSON content.
//
// Returns true if Content-Type starts with "application/json".
//
// Example:
//
//	if IsJSONRequest(r) {
//	    // handle JSON request
//	}
func IsJSONRequest(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Content-Type"), "application/json")
}

// GetIP attempts to retrieve the real client IP address from the HTTP request,
// accounting for common proxy headers such as "X-Forwarded-For" and "X-Real-IP".
//
// If those headers are not set, it falls back to parsing the remote address.
//
// Example:
//
//	ip := GetIP(r)
//	fmt.Println(ip) // Output: "203.0.113.195"
func GetIP(r *http.Request) string {
	// Try X-Forwarded-For header (may contain multiple IPs, take first)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	// Try X-Real-IP header
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}

	// Fallback to RemoteAddr (host:port)
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// BindJSON decodes the JSON payload from the request body into the destination struct.
//
// It closes the request body after reading and disallows unknown JSON fields.
//
// Example:
//
//	var payload MyStruct
//	err := BindJSON(r, &payload)
//	if err != nil {
//	    // handle JSON decoding error
//	}
func BindJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

// RouteContext retrieves the *chi.Context from the request's context,
// allowing access to route parameters and routing information.
//
// Example:
//
//	ctx := RouteContext(r)
//	fmt.Println(ctx.RoutePattern())
func RouteContext(r *http.Request) *chi.Context {
	return chi.RouteContext(r.Context())
}

// URLParam extracts a URL parameter value from the request using the chi router.
//
// Example:
//
//	userID := URLParam(r, "userID")
//	fmt.Println(userID) // Output: "12345"
func URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
