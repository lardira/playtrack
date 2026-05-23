package envutil

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvMode string

const (
	EnvModeDevelopment EnvMode = "development"
	EnvModeProduction  EnvMode = "production"
)

func (e EnvMode) Valid() bool {
	switch e {
	case EnvModeDevelopment,
		EnvModeProduction:
		return true
	default:
		return false
	}
}

func (e EnvMode) String() string {
	return string(e)
}

func MustGet(key string) string {
	s, ok := os.LookupEnv(key)
	if !ok || s == "" {
		log.Panicf("must get env %q", key)
	}

	return s
}

func GetOrDefault(key, def string) string {
	if s, ok := os.LookupEnv(key); !ok || s == "" {
		return def
	} else {
		return s
	}
}

func GetEnvMode() EnvMode {
	envMode, ok := os.LookupEnv("ENV_MODE")
	if !ok {
		return EnvModeDevelopment
	}

	if ok := EnvMode(envMode).Valid(); !ok {
		log.Printf("WARN: undefined envmode %q, fallback to %v", envMode, EnvModeDevelopment)
		return EnvModeDevelopment
	}
	return EnvMode(envMode)
}

func LoadEnvs() error {
	envPath := GetOrDefault("ENV_PATH", "./.env")

	if err := godotenv.Load(envPath); err != nil {
		return err
	}
	return nil
}
