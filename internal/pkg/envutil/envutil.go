package envutil

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func MustGet(key string) string {
	s, ok := os.LookupEnv(key)
	if !ok || s == "" {
		panic(fmt.Sprintf("must get env '%s'", key))
	}

	return s
}

func GetOrDefault(key, def string) string {
	s, ok := os.LookupEnv(key)
	if !ok || s == "" {
		return def
	}

	return s
}

func LoadEnvs() error {
	envPath := GetOrDefault("ENV_PATH", "./.env")

	if err := godotenv.Load(envPath); err != nil {
		return err
	}
	return nil
}
