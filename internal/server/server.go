package server

import (
	"net/http"

	"github.com/authena-ru/courses-organization/internal/config"
)

type Server struct {
	standardServer *http.Server
}

func New(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		standardServer: &http.Server{
			Addr:         ":" + cfg.HTTP.Port,
			Handler:      handler,
			ReadTimeout:  cfg.HTTP.ReadTimeout,
			WriteTimeout: cfg.HTTP.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	return s.standardServer.ListenAndServe()
}
