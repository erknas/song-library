package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerConifg
	PostgresConfig
}

type ServerConifg struct {
	Addr             string        `env:"ADDR"`
	ReadTimeout      time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout     time.Duration `env:"WRITE_TIMEOUT"`
	IdleTimeout      time.Duration `env:"IDLE_TIMEOUT"`
	ThirdPartyAPIURL string        `env:"THIRD_PARTY_API_URL"`
}

type PostgresConfig struct {
	Host          string `env:"POSTGRES_HOST"`
	Port          string `env:"POSTGRES_PORT"`
	User          string `env:"POSTGRES_USER"`
	Password      string `env:"POSTGRES_PASSWORD"`
	DBName        string `env:"POSTGRES_DB"`
	MigrationPath string `env:"MIGRATIONS_PATH"`
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file: %s", err)
	}

	cfg := new(Config)

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("failed to read envs: %s", err)
	}

	return cfg
}
