package config

import (
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	GithubActionSecret string `env:"GITHUB_ACTION_SECRET,required"`
	AppEnv             string `env:"APP_ENV,required"`
	GinMode            string `env:"GIN_MODE"`
	Port               string `env:"PORT" envDefault:"8080"`
	Neo4jURI           string `env:"NEO4J_URI"`
	Neo4jUser          string `env:"NEO4J_USER,required"`
	Neo4jPassword      string `env:"NEO4J_PASSWORD,required"`
	DBName             string `env:"DB_NAME,required"`
}

func NewConfig() *Config {
	if appEnv := os.Getenv("APP_ENV"); appEnv == "" || appEnv == "local" {
		if err := godotenv.Load(); err != nil {
			println("⚠️  No .env file found, relying on system environment variables")
		}
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if os.Getenv("APP_ENV") == "local" || os.Getenv("APP_ENV") == "" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	} else {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse environment variables")
	}

	return &cfg
}
