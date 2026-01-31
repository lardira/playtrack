package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/lardira/playtrack/internal/pkg/envutil"
	"github.com/lardira/playtrack/internal/server"
)

func init() {
	if err := envutil.LoadEnvs(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrChan := make(chan error)

	opts := server.Options{
		Host:        envutil.GetOrDefault("SERVER_HOST", "localhost"),
		Port:        envutil.GetOrDefault("SERVER_PORT", "8080"),
		DatabaseURL: envutil.MustGet("DB_URL"),
		JWTSecret:   envutil.MustGet("JWT_TOKEN_SECRET"),
	}
	server, err := server.New(ctx, opts)
	if err != nil {
		panic(fmt.Errorf("could not set up server: %v", err))
	}
	defer server.Shutdown(ctx)

	go func() {
		defer close(serverErrChan)
		serverErrChan <- server.Run()
	}()

	select {
	case err, ok := <-serverErrChan:
		if ok && err != http.ErrServerClosed {
			log.Printf("error on running: %v", err)
		}

	case <-ctx.Done():
		log.Println("kill signal fired")
	}
}
