package middleware

import (
	"context"
	"net/http"

	"github.com/Etwodev/ramchi/v2/log"
)

// NewLoggingMiddleware creates a Middleware that injects the provided logger
// into the request context. This allows downstream handlers and middleware
// to retrieve the logger via context for structured logging.
//
// This middleware adds the logger under the context key `log.LoggerCtxKey`.
//
// Example usage:
//
//	// Initialize your logger instance (e.g., zerolog.Logger)
//	myLogger := zerolog.New(os.Stdout)
//
//	// Create the logging middleware
//	loggingMiddleware := NewLoggingMiddleware(myLogger)
//
//	// Register middleware with your server
//	s.LoadMiddleware([]Middleware{loggingMiddleware})
//
//	// In an HTTP handler, retrieve the logger:
//	func MyHandler(w http.ResponseWriter, r *http.Request) {
//	    logger := LoggerFromContext(r.Context())
//	    logger.Info().Msg("Handling request")
//	    // ...
//	}
func NewLoggingMiddleware(logger log.Logger) Middleware {
	return NewMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), log.LoggerCtxKey, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}, "ramchi_logger_inject", true, false)
}

// NewCORSMiddleware creates a simple CORS middleware that sets appropriate
// headers to allow cross-origin requests based on the provided allowed origins.
//
// Parameters:
//   - allowedOrigins: a slice of allowed origins (e.g., []string{"https://example.com"}).
//     Use []string{"*"} to allow all origins.
//
// The middleware responds to OPTIONS requests with a 200 status code immediately
// to support CORS preflight requests.
//
// Example usage:
//
//	corsMiddleware := NewCORSMiddleware([]string{"https://example.com", "https://api.example.com"})
//	s.LoadMiddleware([]Middleware{corsMiddleware})
func NewCORSMiddleware(allowedOrigins []string) Middleware {
	return NewMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowed := false

			if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
				allowed = true
			} else {
				for _, o := range allowedOrigins {
					if o == origin {
						allowed = true
						break
					}
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Immediately respond to OPTIONS preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}, "ramchi_cors", true, false)
}
