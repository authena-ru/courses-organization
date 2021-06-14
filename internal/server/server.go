package server

import (
	"github.com/authena-ru/courses-organization/internal/config"
	"net/http"
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
