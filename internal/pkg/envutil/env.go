package envutil

import (
	"log"
	"os"
)

func MustGet(key string) string {
	s := os.Getenv(key)
	if s == "" {
		log.Fatalf("must get env '%s'", key)
	}

	return s
}
