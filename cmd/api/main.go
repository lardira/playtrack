package main

import (
	"log"

	"github.com/lardira/playtrack/internal/server"
)

func main() {
	opts := server.Options{
		Host: "localhost",
		Port: 8080, //TODO: -> .env
	}
	server := server.New(opts)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
	log.Println("server shutdown")
}
