package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lardira/playtrack/internal/db"
	"github.com/lardira/playtrack/internal/tech"
)

type Options struct {
	Host        string
	Port        string
	DatabaseURL string
}

type Server struct {
	Options

	server *http.Server
	db     *pgxpool.Pool
}

func New(ctx context.Context, opts Options) (*Server, error) {
	dbpool, err := db.NewPostgres(ctx, opts.DatabaseURL)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", opts.Host, opts.Port),
		Handler: mux,
	}
	api := humago.New(mux, huma.DefaultConfig("playtrack API", "1.0.0"))

	techHandler := tech.NewHandler(dbpool)

	techHandler.Register(api)

	return &Server{
		Options: opts,
		server:  &server,
		db:      dbpool,
	}, nil
}

func (s *Server) Run() error {
	s.prompt()
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) {
	log.Println("shutting down...")

	if s.db != nil {
		s.db.Close()
	}

	s.server.Shutdown(ctx)
}

func (s *Server) prompt() {
	log.Println("server is running...")
	log.Printf("api - http://%v\n", s.server.Addr)
	log.Printf("docs - http://%v/docs\n", s.server.Addr)
}
