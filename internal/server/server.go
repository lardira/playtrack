package server

import (
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
	address string
}

func New(opts Options) *Server {
	return &Server{
		Options: opts,
		address: fmt.Sprintf("%s:%d", opts.Host, opts.Port),
	}
}

func (s *Server) Run() error {
	s.prompt()

	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("playtrack API", "1.0.0"))

	techHandler := tech.NewHandler()

	techHandler.Register(api)

	return http.ListenAndServe(s.address, router)
}

func (s *Server) prompt() {
	log.Println("server is running...")
	log.Printf("api - http://%v\n", s.address)
	log.Printf("docs - http://%v/docs\n", s.address)
}
