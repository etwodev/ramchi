package ramchi

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	c "github.com/Etwodev/ramchi/config"
	"github.com/Etwodev/ramchi/log"
	"github.com/Etwodev/ramchi/middleware"
	"github.com/Etwodev/ramchi/router"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Server struct {
	idle        chan struct{}
	middlewares []middleware.Middleware
	routers     []router.Router
	instance    *http.Server
	logger      log.Logger
}

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

func (s *Server) Logger() log.Logger {
	return s.logger
}

func (s *Server) LoadRouter(routers []router.Router) {
	s.routers = append(s.routers, routers...)
}

func (s *Server) LoadMiddleware(middlewares []middleware.Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}

func (s *Server) Start() {
	// Load CORS middleware if enabled in config
	if c.EnableCORS() && len(c.AllowedOrigins()) > 0 {
		corsMw := middleware.NewCORSMiddleware(c.AllowedOrigins())
		s.LoadMiddleware([]middleware.Middleware{corsMw})
	}

	// Load Logging middleware if enabled
	if c.EnableRequestLogging() {
		loggingMw := middleware.NewLoggingMiddleware(s.logger)
		s.LoadMiddleware([]middleware.Middleware{loggingMw})
	}

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

func (s *Server) handler() *chi.Mux {
	m := chi.NewMux()
	s.initMux(m)
	return m
}

func (s *Server) initMux(m *chi.Mux) {
	// Global middleware
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

	// Routers
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
