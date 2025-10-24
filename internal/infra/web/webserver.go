package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	chizero "github.com/ironstar-io/chizerolog"
	"github.com/rs/zerolog"
)

type WebServerInterface interface {
	Start()
	Shutdown(ctx context.Context) error
}

type RouteHandler struct {
	Path        string
	Method      string
	HandlerFunc http.HandlerFunc
}

type WebServer struct {
	Server        *http.Server
	Router        chi.Router
	Handlers      []RouteHandler
	WebServerPort int
	Logger        zerolog.Logger
}

func NewWebServer(serverPort int, logger zerolog.Logger, handlers []RouteHandler) *WebServer {
	return &WebServer{
		Server:        nil,
		Router:        chi.NewRouter(),
		Handlers:      handlers,
		WebServerPort: serverPort,
		Logger:        logger,
	}
}

func (s *WebServer) Start() {
	s.Router.Use(chizero.LoggerMiddleware(&s.Logger))
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Recoverer)

	for _, h := range s.Handlers {
		s.Logger.Debug().Msgf("Registering route %s %s", h.Method, h.Path)
		s.Router.MethodFunc(h.Method, h.Path, h.HandlerFunc)
	}

	s.Logger.Info().Msgf("Starting server on port %d", s.WebServerPort)

	s.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.WebServerPort),
		Handler: s.Router,
	}

	go func() {
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Fatal().Err(err).Msg("Failed to start webserver")
		}
	}()
}

func (s *WebServer) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
