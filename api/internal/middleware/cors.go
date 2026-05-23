package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

var (
	allowedOrigins = []string{
		"http://localhost:3000",
	}
)

func CORS(handler http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
	}).Handler(handler)
}
