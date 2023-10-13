package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	CronEnabled bool   `env:"CRON_ENABLED" env-default:"false"`
	ApodApiKey  string `env:"APOD_API_KEY" env-required:"true"`
	Database    struct {
		User     string `env:"POSTGRES_USER" env-required:"true"`
		Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
		Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
		Port     string `env:"POSTGRES_PORT" env-required:"true"`
		DBname   string `env:"POSTGRES_DB" env-required:"true"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}

		if err := godotenv.Load("../../.env"); err != nil {
			log.Fatalf("error load .env file: %v", err)
		}

		if err := cleanenv.ReadEnv(instance); err != nil {
			log.Fatalf("error read .env file: %v", err)
		}
	})
	return instance
}
