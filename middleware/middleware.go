package middleware

import (
	"net/http"
)

type Middleware interface {
	Method()  func(http.Handler) http.Handler
	// Status returns whether the middleware is enabled
	Status()  bool
	// Experimental returns whether the middleware is experimental
	Experimental() bool
}
