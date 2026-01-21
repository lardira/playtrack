package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/lardira/playtrack/internal/pkg/envutil"
	"github.com/lardira/playtrack/internal/server"
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if err := godotenv.Load(path.Join(dir, ".env")); err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrChan := make(chan error)

	opts := server.Options{
		Host:        "localhost",
		Port:        8080, //TODO: -> .env
		DatabaseURL: envutil.MustGet("DB_URL"),
	}
	server, err := server.New(ctx, opts)
	if err != nil {
		log.Fatalf("error on setting up: %v", err)
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
