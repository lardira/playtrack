package envutil

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func MustGet(key string) string {
	s, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("must get env '%s'", key)
	}

	return s
}

func GetOrDefault(key, def string) string {
	s, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	return s
}

func LoadEnvs() error {
	envPath, ok := os.LookupEnv("ENV_PATH")
	if !ok {
		envPath = "./.env"
	}

	if err := godotenv.Load(envPath); err != nil {
		return err
	}

	return nil
}
