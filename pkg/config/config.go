package config

import (
	"errors"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
	"time"
)

type AppConfig struct {
	AppName       string `env:"APP_NAME"`
	AddressServer string `env:"ADDRESS_SERVER"`
	Auth          Auth
	Logger        Logger `envPrefix:"LOGGER_"`
	Telemetry     Telemetry
	MainStorage   struct {
		Postgres PostgresConfig `envPrefix:"POSTGRES_"`
	}
}

type PostgresConfig struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName   string `env:"DB_NAME" envDefault:"postgres"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

type Auth struct {
	JWTSecret        string        `env:"JWT_SECRET"`
	JWTExpiry        time.Duration `env:"JWT_EXPIRY" envDefault:"60m"`
	JWTAlgorithm     string        `env:"JWT_ALGORITHM" envDefault:"HS256"`
	JWTRefreshExpiry time.Duration `env:"JWT_REFRESH_EXPIRY" envDefault:"43200m"`
}

type Logger struct {
	Level      string `env:"LEVEL" envDefault:"info"`
	Output     string `env:"OUTPUT" envDefault:"stdout"`
	FilePath   string `env:"FILE_PATH" envDefault:"./logs/app.log"`
	MaxSizeMB  int    `env:"MAX_SIZE_MB" envDefault:"10"`
	MaxBackups int    `env:"MAX_BACKUPS" envDefault:"5"`
	MaxAgeDays int    `env:"MAX_AGE_DAYS" envDefault:"30"`
	Compress   bool   `env:"COMPRESS" envDefault:"true"`
	AppEnv     string `env:"APP_ENV" envDefault:"development"`
}

type Telemetry struct {
	Host  string `env:"TELEMETRY_HOST" envDefault:"localhost"`
	Port  string `env:"TELEMETRY_PORT" envDefault:"4317"`
	Local bool   `env:"TELEMETRY_LOCAL" envDefault:"true"`
}

func (c *AppConfig) ReadEnvConfig() error {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
		return err
	}
	if err := env.Parse(c); err != nil {
		log.Println("Error parsing config")
		return err
	}
	return nil
}

func (c *AppConfig) Validate() error {
	if c.MainStorage.Postgres.Host == "" {
		return errors.New("no postgres host provided")
	}
	if c.MainStorage.Postgres.Port == "" {
		return errors.New("no postgres port provided")
	}
	if c.MainStorage.Postgres.User == "" {
		return errors.New("no postgres username provided")
	}
	if c.MainStorage.Postgres.Password == "" {
		return errors.New("no postgres password provided")
	}
	if c.MainStorage.Postgres.DBName == "" {
		return errors.New("no postgres database name provided")
	}
	if c.Auth.JWTSecret == "" {
		return errors.New("no jwt secret provided")
	}
	return nil
}
