/*
Package ramchi provides a configurable HTTP server framework with support for
middleware, routers, structured logging, and graceful shutdown.

Example usage:

	package main

	import (
		"encoding/json"
		"net/http"

		"github.com/etwodev/ramchi/v2"
		"github.com/etwodev/ramchi/v2/router"
	)

	func main() {
		s := ramchi.New()

		// Load routers into the server
		s.LoadRouter(Routers())

		// Start the HTTP server (blocks until shutdown)
		s.Start()
	}

	func Routers() []router.Router {
		return []router.Router{
			router.NewRouter("example", Routes(), true, nil),
		}
	}

	func Routes() []router.Route {
		return []router.Route{
			router.NewGetRoute("/demo", true, false, ExampleGetHandler, nil),
		}
	}

	// ExampleGetHandler is a GET handler registered at /example/demo
	func ExampleGetHandler(w http.ResponseWriter, r *http.Request) {
		res, _ := json.Marshal(map[string]string{"success": "ping"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(res)
	}
*/
package ramchi

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	c "github.com/Etwodev/ramchi/v2/config"
	"github.com/Etwodev/ramchi/v2/log"
	"github.com/Etwodev/ramchi/v2/middleware"
	"github.com/Etwodev/ramchi/v2/router"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// Server represents an HTTP server with support for
// configuration, middleware, routers, and structured logging.
type Server struct {
	idle        chan struct{}
	middlewares []middleware.Middleware
	routers     []router.Router
	instance    *http.Server
	logger      log.Logger
}

// New creates a new Server instance with configuration loaded
// and a logger initialized.
//
// It will fatal exit if configuration loading fails.
//
// Example:
//
//	srv := ramchi.New()
func New() *Server {
	err := c.New()
	if err != nil {
		baseLogger := zerolog.New(os.Stdout).With().Timestamp().Str("Group", "ramchi").Logger()
		baseLogger.Fatal().Str("Function", "New").Err(err).Msg("Failed to load config")
	}

	level, err := zerolog.ParseLevel(c.LogLevel())
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	format := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02T15:04:05"}
	baseLogger := zerolog.New(format).With().Timestamp().Str("Group", "ramchi").Logger()

	logger := log.NewZeroLogger(baseLogger)

	return &Server{
		logger: logger,
	}
}

// Logger returns the logger instance used by the server.
//
// Example:
//
//	logger := srv.Logger()
//	logger.Info().Msg("Server logger retrieved")
func (s *Server) Logger() log.Logger {
	return s.logger
}

// LoadRouter appends one or more routers to the server's router list.
//
// Example:
//
//	srv.LoadRouter([]router.Router{myRouter1, myRouter2})
func (s *Server) LoadRouter(routers []router.Router) {
	s.routers = append(s.routers, routers...)
}

// LoadMiddleware appends one or more middleware instances to the server's middleware chain.
//
// Middleware registered here will be applied globally to all routers.
//
// Example:
//
//	srv.LoadMiddleware([]middleware.Middleware{corsMw, loggingMw})
func (s *Server) LoadMiddleware(middlewares []middleware.Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

// Start launches the HTTP server, applying configured middleware and routers,
// and listens for termination signals for graceful shutdown.
//
// It blocks until the server is shut down.
//
// Example:
//
//	srv.Start()
func (s *Server) Start() {
	s.instance = &http.Server{
		Addr:           fmt.Sprintf("%s:%s", c.Address(), c.Port()),
		Handler:        s.handler(),
		ReadTimeout:    time.Duration(c.ReadTimeout()) * time.Second,
		WriteTimeout:   time.Duration(c.WriteTimeout()) * time.Second,
		IdleTimeout:    time.Duration(c.IdleTimeout()) * time.Second,
		MaxHeaderBytes: c.MaxHeaderBytes(),
	}

	s.logger.Debug().
		Str("Port", c.Port()).
		Str("Address", c.Address()).
		Bool("Experimental", c.Experimental()).
		Msg("Server started")

	s.idle = make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		timeout := time.Duration(c.ShutdownTimeout()) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := s.instance.Shutdown(ctx); err != nil {
			s.logger.Warn().Str("Function", "Shutdown").Err(err).Msg("Server shutdown failed!")
		}
		close(s.idle)
	}()

	if c.EnableTLS() {
		s.logger.Info().Msg("Starting HTTPS server")
		if err := s.instance.ListenAndServeTLS(c.TLSCertFile(), c.TLSKeyFile()); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal().Err(err).Msg("HTTPS server failed")
		}
	} else {
		s.logger.Info().Msg("Starting HTTP server")
		if err := s.instance.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal().Err(err).Msg("HTTP server failed")
		}
	}

	<-s.idle

	s.logger.Debug().
		Str("Port", c.Port()).
		Str("Address", c.Address()).
		Bool("Experimental", c.Experimental()).
		Msg("Server stopped")
}

// handler creates and returns the root chi.Mux router for the server.
//
// It initializes the mux with middleware and routers previously loaded.
//
// Example:
//
//	mux := srv.handler()
func (s *Server) handler() *chi.Mux {
	m := chi.NewMux()
	s.initMux(m)
	return m
}

// initMux initializes the provided chi.Mux by registering global middleware,
// routers, and route-level middleware and handlers.
//
// Only middleware and routes enabled and matching the experimental config are registered.
//
// Example:
//
//	mux := chi.NewMux()
//	srv.initMux(mux)
func (s *Server) initMux(m *chi.Mux) {
	if c.EnableRequestLogging() {
		middleware := middleware.NewLoggingMiddleware(s.logger)

		s.logger.Debug().
			Str("Name", middleware.Name()).
			Bool("Experimental", middleware.Experimental()).
			Bool("Status", middleware.Status()).
			Msg("Registering middleware")

		m.Use(middleware.Method())
	}

	if c.EnableCORS() && len(c.AllowedOrigins()) > 0 {
		middleware := middleware.NewCORSMiddleware(c.AllowedOrigins())

		s.logger.Debug().
			Str("Name", middleware.Name()).
			Bool("Experimental", middleware.Experimental()).
			Bool("Status", middleware.Status()).
			Msg("Registering middleware")

		m.Use(middleware.Method())
	}

	for _, middleware := range s.middlewares {
		if middleware.Status() && (middleware.Experimental() == c.Experimental() || !middleware.Experimental()) {
			s.logger.Debug().
				Str("Name", middleware.Name()).
				Bool("Experimental", middleware.Experimental()).
				Bool("Status", middleware.Status()).
				Msg("Registering middleware")

			m.Use(middleware.Method())
		}
	}

	for _, rtr := range s.routers {
		if !rtr.Status() {
			continue
		}

		m.Route("/"+rtr.Prefix(), func(r chi.Router) {
			for _, rmw := range rtr.Middleware() {
				r.Use(rmw)
			}

			for _, rt := range rtr.Routes() {
				if !rt.Status() || (rt.Experimental() != c.Experimental() && rt.Experimental()) {
					continue
				}

				fullPath := "/" + rt.Path()

				s.logger.Debug().
					Bool("Experimental", rt.Experimental()).
					Bool("Status", rt.Status()).
					Str("Method", rt.Method()).
					Str("Path", path.Join("/", rtr.Prefix(), rt.Path())).
					Msg("Registering route")

				finalHandler := http.Handler(rt.Handler())
				for _, mw := range rt.Middleware() {
					finalHandler = mw(finalHandler)
				}

				r.Method(rt.Method(), fullPath, finalHandler)
			}
		})
	}
}
