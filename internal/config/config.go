package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	ServerPort    string
	ServerHost    string
	JWTSecret     string
	AllowedOrigin string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBSSLMode:     os.Getenv("DB_SSL_MODE"),
		ServerPort:    os.Getenv("SERVER_PORT"),
		ServerHost:    os.Getenv("SERVER_HOST"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		AllowedOrigin: os.Getenv("ALLOWED_ORIGIN"),
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(cfg *Config) error {
	requiredVars := map[string]string{
		"DB_HOST":        cfg.DBHost,
		"DB_PORT":        cfg.DBPort,
		"DB_USER":        cfg.DBUser,
		"DB_PASSWORD":    cfg.DBPassword,
		"DB_NAME":        cfg.DBName,
		"SERVER_PORT":    cfg.ServerPort,
		"SERVER_HOST":    cfg.ServerHost,
		"JWT_SECRET":     cfg.JWTSecret,
		"ALLOWED_ORIGIN": cfg.AllowedOrigin,
	}

	for name, value := range requiredVars {
		if value == "" {
			return fmt.Errorf("variável de ambiente %s não está definida", name)
		}
	}

	return nil
}
