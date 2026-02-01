package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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
	Host              string
	Port              string
	DatabaseURL       string
	JWTSecret         string
	CheckPollInterval time.Duration
}

type Server struct {
	Options

	server        *http.Server
	db            *pgxpool.Pool
	healthChecker *tech.HealthChecker
}

func New(ctx context.Context, opts Options) (*Server, error) {
	dbpool, err := db.NewPostgres(ctx, opts.DatabaseURL)
	if err != nil {
		return nil, err
	}

	healthChecker := tech.NewHealthChecker(dbpool, opts.CheckPollInterval, "postgres db")

	mux := http.NewServeMux()
	servMux := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(mux)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", opts.Host, opts.Port),
		Handler: servMux,
	}

	config := huma.DefaultConfig("playtrack API", "1.0.0")
	api := humago.New(mux, config)
	apiV1 := huma.NewGroup(api, "/v1")
	unsecApi := huma.NewGroup(api, "/pub")

	apiV1.UseMiddleware(
		middleware.Authorize(opts.JWTSecret),
	)

	// TODO: use squirell for query building
	gameRepository := game.NewPGRepository(dbpool)
	playerRepository := player.NewPGRepository(dbpool)
	playedGameRepository := player.NewPGPlayedRepository(dbpool)

	techHandler := tech.NewHandler(healthChecker)
	gameHandler := game.NewHandler(gameRepository)
	playerHandler := player.NewHandler(playerRepository, gameRepository, playedGameRepository)
	authHandler := auth.NewHandler(opts.JWTSecret, playerRepository)

	techHandler.Register(apiV1)
	gameHandler.Register(apiV1)
	playerHandler.Register(apiV1)
	authHandler.Register(unsecApi)

	return &Server{
		Options:       opts,
		server:        &server,
		db:            dbpool,
		healthChecker: healthChecker,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	s.prompt()
	go s.healthChecker.Check(ctx)
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
