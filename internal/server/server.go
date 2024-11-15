package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	router     *mux.Router
	address    string
	httpServer *http.Server
	dbpool     *pgxpool.Pool
}

func NewServer(address string, dbpool *pgxpool.Pool) *Server {
	r := mux.NewRouter()
	return &Server{
		router:  r,
		address: address,
		dbpool:  dbpool,
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
	s.router.Use(s.injectServerContext)
	s.router.HandleFunc("/ws", HandleWebSocket)
}

func (s *Server) injectServerContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "server", s)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
