package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Pod-Box/swap2p-backend/api"
	"github.com/Pod-Box/swap2p-backend/repo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"

	"github.com/go-chi/cors"
)

var _ api.ServerInterface = &Server{}

type Server struct {
	httpServer *http.Server
	log        *log.Logger
	db         repo.Repository
}

const (
	applicationJSONContentType = "application/json"
)

type SetupOption func(s *Server) error

func NewServer(cfg *Config, l *log.Logger, opts ...SetupOption) (*Server, error) {
	s := &Server{}
	corsOptions := cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}

	mux := chi.NewRouter()
	mux.Use(middleware.NoCache)
	mux.Use(middleware.SetHeader("Content-Type", applicationJSONContentType))
	mux.Use(cors.Handler(corsOptions))

	mux.Mount("/", api.HandlerWithOptions(s, api.ChiServerOptions{BaseURL: cfg.BasePath}))

	s.httpServer = &http.Server{
		Handler: mux,
		Addr:    cfg.Address,
	}
	s.log = l

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Server) Run() {
	if err := s.httpServer.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info(err)
		} else {
			s.log.Error(err)
		}
	}
}

func SetupWithRepo(r repo.Repository) func(server *Server) error {
	return func(s *Server) error {
		s.db = r
		return nil
	}
}
