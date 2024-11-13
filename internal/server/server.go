package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router     *mux.Router
	address    string
	httpServer *http.Server
}

func NewServer(address string) *Server {
	r := mux.NewRouter()
	return &Server{
		router:  r,
		address: address,
	}
}

func (s *Server) Run() error {
	s.configureRouter()

	s.httpServer = &http.Server{
		Addr:    s.address,
		Handler: s.router,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/ws", HandleWebSocket)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
