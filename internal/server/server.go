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
	"github.com/lardira/playtrack/internal/domain/auth"
	"github.com/lardira/playtrack/internal/domain/game"
	"github.com/lardira/playtrack/internal/domain/player"
	"github.com/lardira/playtrack/internal/middleware"
	"github.com/lardira/playtrack/internal/tech"
	"github.com/rs/cors"
)

type Options struct {
	Host        string
	Port        string
	DatabaseURL string
	JWTSecret   string
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

	config := huma.DefaultConfig("playtrack API", "1.0.0")
	api := humago.New(mux, config)
	apiV1 := huma.NewGroup(api, "/v1")
	unsecApi := huma.NewGroup(api, "/pub")

	apiV1.UseMiddleware(
		middleware.Authorize(opts.JWTSecret),
	)

	cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)

	// TODO: use squirell for query building
	gameRepository := game.NewPGRepository(dbpool)
	playerRepository := player.NewPGRepository(dbpool)
	playedGameRepository := player.NewPGPlayedRepository(dbpool)

	techHandler := tech.NewHandler(dbpool)
	gameHandler := game.NewHandler(gameRepository)
	playerHandler := player.NewHandler(playerRepository, gameRepository, playedGameRepository)
	authHandler := auth.NewHandler(opts.JWTSecret, playerRepository)

	techHandler.Register(apiV1)
	gameHandler.Register(apiV1)
	playerHandler.Register(apiV1)
	authHandler.Register(unsecApi)

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
