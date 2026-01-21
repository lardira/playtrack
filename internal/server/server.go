package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/lardira/playtrack/internal/tech"
)

type Options struct {
	Host string
	Port int
}

type Server struct {
	Options

	server *http.Server
	db     *sql.DB
}

func New(opts Options) *Server {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Handler: mux,
	}
	api := humago.New(mux, huma.DefaultConfig("playtrack API", "1.0.0"))

	techHandler := tech.NewHandler()

	techHandler.Register(api)

	return &Server{
		Options: opts,
		server:  &server,
	}
}

func (s *Server) Run() error {
	s.prompt()
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() {
	log.Println("shutting down...")

	if s.db != nil {
		s.db.Close()
	}

	s.server.Shutdown(context.Background())
}

func (s *Server) prompt() {
	log.Println("server is running...")
	log.Printf("api - http://%v\n", s.server.Addr)
	log.Printf("docs - http://%v/docs\n", s.server.Addr)
}
