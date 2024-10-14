package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		DBUrl: getEnv("DB_URL", "postgresql://myuser:mypassword@localhost:5432/dbname?schema=public"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
