package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/lardira/playtrack/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrChan := make(chan error)

	server := server.New(server.Options{
		Host: "localhost",
		Port: 8080, //TODO: -> .env
	})
	defer server.Shutdown()

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
