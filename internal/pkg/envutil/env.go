package envutil

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func MustGet(key string) string {
	s := os.Getenv(key)
	if s == "" {
		log.Fatalf("must get env '%s'", key)
	}

	return s
}

func LoadEnvs() error {
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = "./.env"
	}

	if err := godotenv.Load(envPath); err != nil {
		return err
	}

	return nil
}
