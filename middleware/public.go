package middleware

import (
	"context"
	"net/http"

	"github.com/Etwodev/ramchi/log"
)

// LoggerInjectionMiddleware returns a Middleware that injects the provided logger
// instance into the request's context. This allows downstream handlers and middleware
// to retrieve the logger directly from the context for structured logging.
// If your preferred logging library is not supported, please raise an issue on this repo.
//
// Usage:
//
//		// Create the logger (e.g., in main.go)
//		myLogger := zerolog.New(format)
//
//		// Create the middleware
//	 func Middlewares() []middleware.Middleware {
//		 return []middleware.Middleware{
//			 middleware.NewLoggingMiddleware(myLogger),
//			 middleware.NewMiddleware(auth.Middleware(), "auth", true, false),
//		 }
//	 }
//
//		// Load the middleware
//		s.LoadMiddleware(Middlewares())
//
//	 // In your handlers, you can retrieve the logger from the context like this:
//
//	 func MyHandler(w http.ResponseWriter, r *http.Request) {
//		 logger := middleware.LoggerFromContext(r.Context())
//		 logger.Info().Msg("Handling request")
//		 // ...
//	 }
func NewLoggingMiddleware(logger log.Logger) Middleware {
	return NewMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), log.LoggerCtxKey, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}, "ramchi_logger_inject", true, false)
}

// NewCORSMiddleware returns a simple CORS middleware.
// allowedOrigins is a list of origins that are allowed. Use ["*"] for allowing all.
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

			// For OPTIONS requests, respond with 200 immediately (CORS preflight)
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}, "ramchi_cors", true, false)
}
