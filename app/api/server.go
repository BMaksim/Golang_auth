package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	config *Config
	router *mux.Router
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		router: mux.NewRouter(),
	}
}

func (s *Server) StartServer() error {
	return http.ListenAndServe(s.config.Addr, s.router)
}

func (s *Server) setRoutes() {
	s.router.HandleFunc("/", s.handleFirst())
}

func (s *Server) handleFirst() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("first route"))
	}
}
